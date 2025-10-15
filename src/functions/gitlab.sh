#!/bin/sh


gitlab_authenticate()
{
    HOST=${1#https://}; HOST=${HOST#http://}
    CONFIG_FILE="$HOME/.config/glab-cli/config.yml"

    if [ -z "$HOST" ]; 
    then
        echo "Error: Missing hostname argument. Usage: gitlab_authenticate <hostname>" >&2
        return 1
    fi

    if [ ! -f "$CONFIG_FILE" ]; 
    then
        echo "No glab config found. Running login for $HOST..." >&2
        glab auth login --hostname "$HOST"
    fi

    GITLAB_TOKEN=$(sed -E 's/([!]*null)+//g' ${CONFIG_FILE} | yq -r '.hosts."'"${HOST}"'".token' | sed 's/["'"'"']//g')

    if [ -z "$GITLAB_TOKEN" ]; 
    then
        echo "Token not found or invalid in config file" >&2
        read_passwd "GITLAB_TOKEN"
        if [ -z "$GITLAB_TOKEN" ]; 
        then
            echo " No token provided. Aborting." >&2
            return 1
        fi
    fi
    glab auth login --hostname "$HOST" --token "$GITLAB_TOKEN"

#    export GITLAB_TOKEN
    echo "Authenticated to $HOST" >&2
}

# Main API call function:
# Usage: 
# Example: gitlab_api GET /user
#          gitlab_api POST /projects/:id/issues "title=Bug&description=Desc"
gitlab_api()
{
    METHOD=$(echo "$1" | tr '[:lower:]' '[:upper:]')
    RESOURCE_PATH=$2
    DATA=$3
    if [ -z "$METHOD" ] || [ -z "$RESOURCE_PATH" ]; 
    then
        echo "Error: Missing argument. 
          Usage: 
            gitlab_api METHOD RESOURCE_PATH [DATA]
          Example:
            gitlab_api GET /user
            gitlab_api POST /projects/:id/issues title=Bug&description=Desc
        " >&2
        return 1
    elif [ -z "$CI_SERVER_URL" ]; 
    then
        echo "Error: CI_SERVER_URL is not set." >&2
        return 1
    fi
    HOST=${CI_SERVER_URL#https://}; HOST=${HOST#http://}

    gitlab_authenticate ${HOST} || return 1

    if [ "$METHOD" = "GET" ]; 
    then
        glab api --paginate --hostname "$HOST" "$RESOURCE_PATH"
    else
        glab api --paginate --method "$METHOD" --hostname "$HOST" "$RESOURCE_PATH" --form "$DATA"
    fi
}

gitlab_vars() {
    PROJECT_ID=$1
    ENVIRONMENT_NAME=$2

    if [ -z "$PROJECT_ID" ] || [ -z "$ENVIRONMENT_NAME" ]; then
        echo "Usage: gitlab_vars <project_id> <environment_name>" >&2
        return 1
    fi

    RESPONSE=$(gitlab_api GET "/projects/${PROJECT_ID}/variables") || return 1

    VAR_COUNT=$(echo "$RESPONSE" | jq 'length')

    if [ "$VAR_COUNT" -eq 0 ]; then
        echo "No variables found in project '$PROJECT_ID'" >&2
        return 1
    fi

    touch /tmp/gitlab_vars.env
    echo "$RESPONSE" | jq -c '.[]' | while read -r VAR; do
        ENV_SCOPE=$(echo "$VAR" | jq -r '.environment_scope')
        VAR_TYPE=$(echo "$VAR" | jq -r '.variable_type')

        if [ "$ENV_SCOPE" != "$ENVIRONMENT_NAME" ] \
        && [ "$ENV_SCOPE" != "*" ] \
        || [ "$VAR_TYPE" != "env_var" ]; 
        then
            continue
        fi

        NAME=$(echo "$VAR" | jq -r '.key')
        VALUE=$(echo "$VAR" | jq -r '.value')
        MASKED=$(echo "$VAR" | jq -r '.masked')

        if [ "$MASKED" = "true" ]; then
            echo "${NAME}=*****"
        else
            echo "${NAME}=${VALUE}"
        fi

        ESCAPED_VALUE=$(printf '%s' "$VALUE" | sed -e 's/\\/\\\\/g' -e 's/"/\\"/g')
        echo "${NAME}=\"${ESCAPED_VALUE}\"" >> /tmp/gitlab_vars.env
    done
    set -a
    . /tmp/gitlab_vars.env
    set +a
    rm -fr /tmp/gitlab_vars.env
}

release()
{
    release-cli create \
    --name "v${1}" \
    --description "$(cat ${CI_PROJECT_DIR}/NOTES)" \
    --tag-name "${1}" \
    --assets-link "{\"url\":\"${2}\",\"name\":\"Documentation\",\"link_type\":\"other\"}" \
    --assets-link "{\"url\":\"${2}/openid\",\"name\":\"API documentation\",\"link_type\":\"other\"}" \
    --assets-link "{\"url\":\"${2}/admin\",\"name\":\"Administration UI\",\"link_type\":\"other\"}" \
    --assets-link "{\"url\":\"https://${CI_REGISTRY_IMAGE}:${1}\",\"name\":\"Docker Image\",\"link_type\":\"other\"}"
}

gitlab_release() {
  tag="$1"
  base_url="$2"
  HOST=${3#https://}; HOST=${HOST#http://}
  changelog=$(/ops/functions/notes.py "${CI_PROJECT_DIR:-.}/CHANGELOG" "${IMAGE_VERSION}")

  if [ -z "$tag" ] || [ -z "$base_url" ];
  then
    echo "Usage: gitlab_release <tag> <base_url>" >&2
    return 1
  fi

  if [ -z "$GITLAB_TOKEN" ]; then
    echo "Warning: GITLAB_TOKEN is not set. The command may fail without authentication." >&2
  fi

  assets_json=$(cat <<EOF
[
  {"name": "Documentation", "url": "${base_url}", "link_type": "other"},
  {"name": "API documentation", "url": "${base_url}/openid", "link_type": "other"},
  {"name": "Administration UI", "url": "${base_url}/admin", "link_type": "other"}
]
EOF
)

  original_remote=$(git remote get-url origin)
  release_remote=$(git remote get-url "$HOST")
  git remote set-url origin "$release_remote"

  echo "Creating release v$tag..."

  glab release create "$tag" \
    --name "v$tag" \
    --notes "$changelog" \
    --assets-links "$assets_json"

  git remote set-url origin "$original_remote"

  if [ $? -eq 0 ]; then
    echo "Release v$tag successfully created."
  else
    echo "Failed to create a release." >&2
    return 1
  fi
}

gitlab_read() {
  environment="$1"
  image_version="$2"
  project_id="$3"
  yaml_path="$4"
  ref="${5:-main}"         # branch ou commit (default: main)

  if [ -z "$environment" ] || [ -z "$image_version" ] || [ -z "$project_id" ] || [ -z "$yaml_path" ]; then
    echo "Usage: infraops <environment> <image_version> <project_id> <file_path> [ref]" >&2
    return 1
  fi

  if [ -z "$GITLAB_TOKEN" ]; then
    echo "Warning: GITLAB_TOKEN is not set. The command may fail without authentication." >&2
  fi

  if [ -z "$CI_SERVER_URL" ] || [ -z "$CD_PROJECT_ID" ]; then
    echo "Error: CI_SERVER_URL or CD_PROJECT_ID not set." >&2
    return 1
  fi

  # URL-encode file path
  file_enc=$(printf '%s' "$yaml_path" | jq -s -R -r @uri)

  response=$(curl -sS \
    --header "Authorization: Bearer ${GITLAB_TOKEN}" \
    "${CI_SERVER_URL}/api/v4/projects/${CD_PROJECT_ID}/repository/files/${file_enc}/raw?ref=${ref}")

  echo "$response" | sed -E \
    -e "s|(image: .*:)[0-9]+\.[0-9]+\.[0-9]+(-[a-z]+[0-9]+)?|\1${image_version}|" \
    -e "/- name: *(IMAGE_VERSION|DD_VERSION)/{n;s|value: *[0-9]+\.[0-9]+\.[0-9]+(-[a-z]+[0-9]+)?|value: ${image_version}|}" \
    | jq -R -s .
}

gitlab_commit()
{
  environment="$1"
  image_version="$2"
  project_id="$3"
  file_path="$4"

  if [ -z "$environment" ] || [ -z "$image_version" ] || [ -z "$project_id" ] || [ -z "$file_path" ]; then
    echo "Usage: infraops <environment> <image_version> <project_id> <file_path>" >&2
    return 1
  fi

  if [ -z "$GITLAB_TOKEN" ]; then
    echo "Warning: GITLAB_TOKEN is not set. The command may fail without authentication." >&2
  fi

  yaml_path="${file_path}/${file_path}.yaml"
  commit_msg="v${image_version}"

payload=$(cat <<-JSON
{
    "id": "${project_id}",
    "branch": "main",
    "commit_message": "${commit_msg}",
    "actions": [{"action": "update","file_path": "${yaml_path}","content": $(gitlab_read "$environment" "$image_version" "$project_id" "$yaml_path")}]
}
JSON
)

  echo "Committing update to '${yaml_path}' in project ${project_id}..."

  echo "${payload}"
  curl \
  --request POST \
  --data "$payload" \
  --header "AUTHORIZATION: Bearer ${GITLAB_TOKEN}" \
  --header "Content-Type: application/json" \
  "${CI_SERVER_URL}/api/v4/projects/${project_id}/repository/commits"

  echo ""
  if [ $? -eq 0 ]; then
    echo "Commit sent successfully."
  else
    echo "Failed to send commit." >&2
    return 1
  fi
}
