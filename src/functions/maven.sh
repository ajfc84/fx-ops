#!/bin/sh


maven_deploy()
(
    CWD=$(pwd)
    cd "${SUB_PROJECT_DIR}"

    mvn --settings "${CI_PROJECT_DIR}/ci_settings.xml" -DskipTests deploy

    cd "${CWD}"
)

maven_install()
{
    image_version="$1"
    is_release="$2"

    CWD=$(pwd)
    cd "${SUB_PROJECT_DIR}"

    if [ "${is_release}" == "true" ];
    then
        mvn versions:set -DnewVersion="${image_version}"
        mvn versions:commit
    fi
    mvn -DskipTests clean install

    cd "${CWD}"
}
