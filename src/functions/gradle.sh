#!/bin/sh


gradle()
{
    CWD=$(pwd)
    cd "${CI_PROJECT_DIR}"

    ./gradlew $@

    cd "${CWD}"
}
