#!/bin/sh


. "${CI_PROJECT_DIR}/buildSh/artifacts.sh"


echo "INFO: Uploading distribution archive..."
upload_artifacts "${CI_PROJECT_NAME}-v${LATEST_VERSION}.zip"
