CI_PROJECT_NAME="fx-ops"
echo "CI_PROJECT_NAME: ${CI_PROJECT_NAME}"
REGISTRY_IMAGE="ajfc84/gitlab-default"
echo "REGISTRY_IMAGE: ${REGISTRY_IMAGE}"
CI_REGISTRY="registry-1.docker.io"
echo "CI_REGISTRY: ${CI_REGISTRY}"
CI_REGISTRY_USER="ajfc84"
echo "CI_REGISTRY_USER: ${CI_REGISTRY_USER}"
# if [ -z "$CI" ];
# then
#   echo -n "
#   Select Version:
#       1) ${LATEST_VERSION}
#       *) ${IMAGE_VERSION}(default)
#   or press Enter(default).
#   "
#   read selection
#   case $selection in
#       1) IMAGE_VERSION=${LATEST_VERSION};;
#       *) IMAGE_VERSION=${IMAGE_VERSION};;
#   esac
#   echo "IMAGE_VERSION: ${IMAGE_VERSION}"
# fi