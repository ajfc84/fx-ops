#!/bin/sh


# argocd() {
#     sed -e "s/\${IMAGE_VERSION}/${2}/" "${BASE_DIR}/k8s.${1}.yaml" | jq --raw-input --slurp
# }

# infraops()
# {
# PAYLOAD=$(cat <<-JSON
# {
#     "id": "${3}",
#     "branch": "main",
#     "commit_message": "v${IMAGE_VERSION}",
#     "actions": [{"action": "update","file_path": "${4}/${4}.yaml","content": $(argocd $1 $2)}]
# }
# JSON
# )
#     echo "${PAYLOAD}"
#     curl \
#     --request POST \
#     --data "$PAYLOAD" \
#     --header "AUTHORIZATION: Bearer ${ARGOCD_TOKEN}" \
#     --header "Content-Type: application/json" \
#     "${CI_SERVER_URL}/api/v4/projects/${3}/repository/commits"
# }
