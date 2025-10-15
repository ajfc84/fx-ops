#!/bin/sh


if [ "$(basename $(readlink /proc/$$/exe))" != "dash" ];
then 
    echo "Cant run this script $0 with $(basename $(readlink /proc/$$/exe))";
    exit; 
fi

cat <<EOF
 ______   __     ______   ______     __         __     __   __     ______    
/\  == \ /\ \   /\  == \ /\  ___\   /\ \       /\ \   /\ "-.\ \   /\  ___\   
\ \  _-/ \ \ \  \ \  _-/ \ \  __\   \ \ \____  \ \ \  \ \ \-.  \  \ \  __\   
 \ \_\    \ \_\  \ \_\    \ \_____\  \ \_____\  \ \_\  \ \_\\"\_\  \ \_____\ 
  \/_/     \/_/   \/_/     \/_____/   \/_____/   \/_/   \/_/ \/_/   \/_____/ 
                                                                             
EOF

set -e

if [ -f /.dockerenv ];
then
    echo "Running in Container mode"
else
    echo "Pipeline must run in a container"
    exit 0
fi

command -v docker >/dev/null 2>&1 || {
    echo "ERROR: Docker is not installed. Run './ops.sh -s' to install dependencies."
    exit 1
}
if [ -n "$DOCKER_HOST" ]; then
    echo "Checking Docker connection via DOCKER_HOST=$DOCKER_HOST ..."
    if ! docker version >/dev/null 2>&1; then
        echo "ERROR: Could not connect to Docker daemon at $DOCKER_HOST"
        exit 1
    fi
elif [ ! -S /var/run/docker.sock ];
then
    echo "ERROR: Docker socket not found at /var/run/docker.sock"
    exit 1
fi
if [ ! -f "$HOME/.ssh/id_rsa" ]; then
    echo "ERROR: Required private key not found at $HOME/.ssh/id_rsa"
    exit 1
else
    chmod 600 "$HOME/.ssh/id_rsa" 2>/dev/null || {
        echo "Warning: Could not chmod $HOME/.ssh/id_rsa (may already be correct)"
    }
fi

main_file=""
if [ -f "${CI_PROJECT_DIR}/main.yaml" ]; then
  main_file="${CI_PROJECT_DIR}/main.yaml"
elif [ -f "${CI_PROJECT_DIR}/main.yml" ]; then
  main_file="${CI_PROJECT_DIR}/main.yml"
else
  echo "Erro: main.yaml ou main.yml não encontrado em ${CI_PROJECT_DIR}" >&2
  exit 1
fi

functions="/ops/functions"
for f in "${functions}"/*.sh;
do
    . "${f}"
done

. ${CI_PROJECT_DIR}/env.sh

last_stage=$(echo "${1}" | xargs)
if [ "${last_stage}" = "version" ];
then
    echo "${IMAGE_VERSION}"
elif [ "${last_stage}" = "sops" ];
then
    sops_dec_enc "${main_file}"
elif [ "${last_stage}" = "build" ] || [ "${last_stage}" = "install" ] || [ "${last_stage}" = "deploy" ];
then
    # secrets & alias
    if [ "${CI_ENVIRONMENT_NAME}" != "local" ] \
    && [ -z "${CI}" ];
    then
        #gitlab_vars "${CI_PROJECT_ID}" "${CI_ENVIRONMENT_NAME}"
        sops_read "${main_file}"
        if [ -z "$REGISTRY_TOKEN" ];
        then
            echo "ERROR: REGISTRY_TOKEN not set from GitLab vars." >&2
            exit 1
        fi
        export CI_REGISTRY_PASSWORD="$REGISTRY_TOKEN"
    fi

    . ${CI_PROJECT_DIR}/secret.sh

    CHANGELOG=$(/ops/functions/notes.py "${CI_PROJECT_DIR:-.}/CHANGELOG" "${IMAGE_VERSION}")
    if [ -z "$CHANGELOG" ];
    then
        echo "ERROR: CHANGELOG not set for version: ${IMAGE_VERSION}" >&2
        exit 1
    fi

    args=$(echo "${2}" | xargs)
    if [ -n "${args}" ];
    then
        is_multi_project=false
        projects=$(echo "${args}" | awk '{split($0, a, " "); print a[1]}')
    else
        is_multi_project=true
        projects=$(yq -r '.projects[]' "${main_file}")
    fi

    if [ "${is_multi_project}" = "true" ] && [ "${last_stage}" != "build" ];
    then
        stage_order="build ${last_stage}"
    else
        stage_order="${last_stage}"
    fi

    for stage in ${stage_order};
    do
        for project in ${projects};
        do
            SUB_PROJECT_DIR="${CI_PROJECT_DIR}/${project}"
            if [ -f "${SUB_PROJECT_DIR}"/"${stage}".sh ];
            then
                echo "${stage}ing ${project} project in ${SUB_PROJECT_DIR}"
                #(
                if [ -f "${SUB_PROJECT_DIR}"/env.sh ];
                then
                    . "${SUB_PROJECT_DIR}/env.sh"
                fi
                if [ -f "${SUB_PROJECT_DIR}"/secret.sh ];
                then
                    . "${SUB_PROJECT_DIR}/secret.sh" "${stage}"
                fi
                . "${SUB_PROJECT_DIR}/${stage}.sh"
                #)
                if [ "${stage}" = "build" ];
                then
                    if [ "${is_multi_project}" = "true" ];
                    then
                        DOCKERFILE=$(find "${SUB_PROJECT_DIR}" -type f -name "Dockerfile.${CI_ENVIRONMENT_NAME}" | head -n 1)
                        if [ -f "${DOCKERFILE}" ];
                        then
                            docker_push
                        else
                            echo "WARN: DOCKERFILE not found." >&2
                        fi
                    fi
                fi
            fi
        done
        if [ "${stage}" = "build" ];
        then
            if [ "${is_multi_project}" = "false" ];
            then
                git_tag "${IMAGE_VERSION}"
            else
                if [ "${CI_ENVIRONMENT_NAME}" != "local" ];
                then
                    # gitlab_release "${IMAGE_VERSION}" "https://${DOMAIN}" "${CI_SERVER_URL}"
                    git_tag "${IMAGE_VERSION}"
                fi
            fi
        fi
    done
else
    echo "\e[31mUnexpected option: ${1}"
    echo "\n"
    exit 1
fi
