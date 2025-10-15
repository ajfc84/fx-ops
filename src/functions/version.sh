#!/bin/sh


version()
{
    CWD=$(pwd)
    cd "${CI_PROJECT_DIR}"
    OLD_IFS=$IFS
    IFS=.

    # TODO "git tag 0.0.1" for new projects
    if [ "${CI_ENVIRONMENT_NAME}" != "local" ];
    then
        git_fetch_tags
    fi

    last_version=$(git tag | grep -E "^[0-9]+\.[0-9]+\.[0-9]+$" | awk -F'.' '{printf "%03d.%03d.%03d\n", $1, $2, $3}' | sort | tail -n 1)
    if [ -z "$last_version" ];
    then
        last_version="000.000.000"
    fi

    if [ "${CI_ENVIRONMENT_NAME}" = "main" ];
    then
        for v in $last_version;
        do
            set -- "$@" $(echo "${v}" | sed -e "s/[0]\{,2\}//")
        done
        major=$((${1})) 
        minor=$((${2}))
        patch=$((${3}))
        shift $#

        if [ "$PATCH" = "true" ];
        then
            patch=$(expr $patch + 1)
            new_version="${major}.${minor}.${patch}"
        else
            minor=$(expr $minor + 1)
            new_version="${major}.${minor}.0"
        fi
    else
        pre_release=$(get_release)
        dev_version=$(git tag | grep -E "^[0-9]+\.[0-9]+\.[0-9]+-${pre_release}[0-9]+$" | sed -e "s/-${pre_release}/./;" | awk -F'.' '{printf "%03d.%03d.%03d.%03d\n", $1, $2, $3, $4}' | sort | tail -n 1)
        if [ -z "$dev_version" ];
        then
            dev_version="${last_version}.000"
        fi
        for v in $last_version;
        do
            set -- "$@" $(echo "${v}" | sed -e "s/[0]\{,2\}//")
        done
        last_version_major=$((${1})) 
        last_version_minor=$((${2}))
        last_version_patch=$((${3}))
        shift $#
        for v in $dev_version;
        do
            set -- "$@" $(echo "${v}" | sed -e "s/[0]\{,2\}//")
        done
        dev_version_major=$((${1})) 
        dev_version_minor=$((${2}))
        dev_version_patch=$((${3}))
        dev_version_meta=$((${4}))
        shift $#

        if [ "${last_version_major}" =  "${dev_version_major}" ] \
        && [ "${last_version_minor}" -ge  "${dev_version_minor}" ];
        then
                minor=$(expr $last_version_minor + 1)
                new_version="${dev_version_major}.${minor}.0-${pre_release}1"
        else
            if [ "$PATCH" = "true" ];
            then
                patch=$(expr $dev_version_patch + 1)
                new_version="${dev_version_major}.${dev_version_minor}.${patch}-${pre_release}1"
            else
                meta=$(expr ${dev_version_meta} + 1)
                new_version="${dev_version_major}.${dev_version_minor}.${dev_version_patch}-${pre_release}${meta}"
            fi
        fi
    fi

    IFS=$OLD_IFS
    cd "${CWD}"
    echo "${new_version}"
}
