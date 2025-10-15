#!/bin/sh


docker_build() {
    ci_registry_image="${1}"
    image_version="${2}"
    srcs_dir="${SUB_PROJECT_DIR}/${3}/"
    shift 3

    author_name="$(git config --get user.name)"
    author_email="$(git config --get user.email)"
    authors="${author_name} <${author_email}>"
    revision="$(git rev-parse --short HEAD)"
    changelog="$(printf %s "$CHANGELOG" | tr '\n' ' ' | sed 's/[[:space:]]\+/ /g' | cut -c1-1000)"

    docker buildx build --no-cache \
        $@ \
        --label "org.opencontainers.image.title=${CI_PROJECT_NAME}" \
        --label "org.opencontainers.image.description=${PROJECT_DESCRIPTION:-n/a}" \
        --label "org.opencontainers.image.url=${DOMAIN}" \
        --label "org.opencontainers.image.source=${CI_SERVER_URL}" \
        --label "org.opencontainers.image.version=${image_version}" \
        --label "org.opencontainers.image.revision=${revision}" \
        --label "org.opencontainers.image.created=$(date -u '+%Y-%m-%dT%H:%M:%SZ')" \
        --label "org.opencontainers.image.authors=${authors}" \
        --label "org.opencontainers.image.documentation=${DOMAIN}" \
        --label "org.opencontainers.image.licenses=Apache-2.0" \
        --label "org.opencontainers.image.vendor=Fx" \
        --label "org.opencontainers.image.notes=${changelog}" \
        -t "${ci_registry_image}:${image_version}" \
        -f "${srcs_dir}/Dockerfile.${CI_ENVIRONMENT_NAME}" \
        "${srcs_dir}"

    echo
    echo "Labels for ${ci_registry_image}:${image_version}:"
    docker inspect "${ci_registry_image}:${image_version}" --format '{{ json .Config.Labels }}' | jq .
}

docker_login()
{
    if [ "${CI_ENVIRONMENT_NAME}" = "local" ];
    then
        echo "Local environment does not need registry login"
    elif [ "${CI_REGISTRY}" = "AWS" ];
    then
        echo "AWS registry login: "
        aws ecr get-login-password --region sa-east-1 | docker login --username "${CI_REGISTRY_USER}" --password-stdin "${CI_REGISTRY}"
    elif [ "${CI_REGISTRY}" = "registry-1.docker.io" ];
    then
        echo "DockerHub registry login: "
        echo "${DOCKER_HUB_PASSWORD}" | docker login --username "${CI_REGISTRY_USER}" --password-stdin
    else
        echo "${CI_REGISTRY} registry login: "
        echo "${CI_REGISTRY_PASSWORD}" | docker login --username "${CI_REGISTRY_USER}" --password-stdin "${CI_REGISTRY}"
    fi
}

docker_push()
{
    if [ "${CI_ENVIRONMENT_NAME}" = "local" ];
    then
        echo "Local environment does not need registry push"
    else
        docker_login
        if [ "${CI_REGISTRY}" = "registry-1.docker.io" ];
        then
            docker push "${REGISTRY_IMAGE}:${IMAGE_VERSION}"
        else
            docker tag "${REGISTRY_IMAGE}:${IMAGE_VERSION}" "${CI_REGISTRY}/${REGISTRY_IMAGE}:${IMAGE_VERSION}"
            docker push "${CI_REGISTRY}/${REGISTRY_IMAGE}:${IMAGE_VERSION}"
        fi
    fi
}

docker_pull()
{
    ci_registry_image="${1}"
    image_version="${2}"

    if docker image inspect "${ci_registry_image}:${image_version}" >/dev/null 2>&1; then
        echo "INFO: Image ${ci_registry_image}:${image_version} already exists locally"
    else
        docker pull "${CI_REGISTRY}/${ci_registry_image}:${image_version}" || {
            echo "Error: could not pull the image."
            return 1
        }
    fi
}

docker_run()
{
    project_name="${1}"
    ci_registry_image="${2}"
    image_version="${3}"
    shift 3

    if [ "$(docker ps -a | grep -c ${project_name})" -gt 0 ];
    then
        docker container rm -f "${project_name}"
    fi

    docker_login
    image=$(docker_pull "${ci_registry_image}" "${image_version}")

    if [ -z "$CI" ];
    then    
        docker run -d \
        --network host \
        --restart always \
        $@ \
        --name "${project_name}" \
        "${image}"
    else
        docker run -d \
        --restart always \
        $@ \
        --name "${project_name}" \
        "${image}"
    fi
}

docker_exec()
{
    project_name="${1}"
    ci_registry_image="${2}"
    image_version="${3}"
    cmd="$(echo ${4} | awk '{split($0, array); print array[1]}')"
    args="$(echo ${4} | awk '{split($0, array); for (i=2; i <= length(array); i++) print array[i]}')"
    shift 4

    if [ "$(docker ps -a | grep -c ${project_name})" -gt 0 ];
    then
        docker container rm -f "${project_name}"
    fi

    docker_login
    image=$(docker_pull "${ci_registry_image}" "${image_version}")

    if [ -z "$CI" ];
    then
        docker run -it \
        --network host \
        $@ \
        --name "${project_name}" \
        "${image}" \
        ${cmd} \
        ${args}
    else
        docker run -it \
        $@ \
        --name "${project_name}" \
        "${image}" \
        ${cmd} \
        ${args}
    fi
}

docker_commit_container_to_image()
{
    docker commit ${container_id} ${image_name}
}

docker_run_shell()
{
    docker run -it --entrypoint sh ${image_name}
}