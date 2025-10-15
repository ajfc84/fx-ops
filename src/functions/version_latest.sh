#!/bin/sh


latest_version()
{
    CWD=$(pwd)
    cd "${CI_PROJECT_DIR}"
    OLD_IFS=$IFS
    IFS=.

    if [ "${CI_ENVIRONMENT_NAME}" != "local" ];
    then
        git_fetch_tags
    fi

    # last_version=$(git tag | egrep -e "^[0-9]+.[0-9]+.[0-9]+$" | sed -e "s/\b\([0-9]\)\b/00\1/g; s/\b\([0-9]\{2\}\)\b/0\1/g;" | sort | tail -n 1)
    last_version=$(git tag | grep -E "^[0-9]+\.[0-9]+\.[0-9]+$" | awk -F'.' '{printf "%03d.%03d.%03d\n", $1, $2, $3}' | sort | tail -n 1)
    if [ -z "$last_version" ];
    then
        last_version="000.000.000"
    fi

    if [ "${CI_ENVIRONMENT_NAME}" = "main" ];
    then
        # set -- $last_version
        # major=$(expr "$1" + 0)
        # minor=$(expr "$2" + 0)
        # patch=$(expr "$3" + 0)
        # shift $#
        for v in $last_version;
        do
            set -- "$@" $(echo "${v}" | sed -e "s/[0]\{,2\}//")
        done
        major=$((${1})) 
        minor=$((${2}))
        patch=$((${3}))
        shift $#

        echo "${major}.${minor}.${patch}"
    else
        pre_release=$(get_release)
        # dev_version=$(git tag | egrep -e "^[0-9]+.[0-9]+.[0-9]+-${pre_release}[0-9]+$" | sed -e "s/-${pre_release}/./; s/\b\([0-9]\)\b/00\1/g; s/\b\([0-9]\{2\}\)\b/0\1/g;" | sort | tail -n 1)
        dev_version=$(git tag | grep -E "^[0-9]+\.[0-9]+\.[0-9]+-${pre_release}[0-9]+$" | sed -e "s/-${pre_release}/./;" | awk -F'.' '{printf "%03d.%03d.%03d.%03d\n", $1, $2, $3, $4}' | sort | tail -n 1)
        if [ -z "$dev_version" ];
        then
            dev_version="${last_version}.000"
        fi

        # set -- $dev_version
        # major=$(expr "$1" + 0)
        # minor=$(expr "$2" + 0)
        # patch=$(expr "$3" + 0)
        # meta=$(expr "$4" + 0)
        # shift $#
        for v in $dev_version;
        do
            set -- "$@" $(echo "${v}" | sed -e "s/[0]\{,2\}//")
        done
        major=$((${1})) 
        minor=$((${2}))
        patch=$((${3}))
        meta=$((${4}))
        shift $#

        echo "${major}.${minor}.${patch}-${pre_release}${meta}"
    fi

    IFS=$OLD_IFS
    cd "${CWD}"
}
