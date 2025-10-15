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

    CWD=$(pwd)
    cd "${SUB_PROJECT_DIR}"

    mvn versions:set -DnewVersion="${image_version}"
    mvn versions:commit
    mvn -DskipTests clean install

    cd "${CWD}"
}
