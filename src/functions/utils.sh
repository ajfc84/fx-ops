#!/bin/sh


get_release()
{
    pre_release="local"
    if [ "${CI_ENVIRONMENT_NAME}" = "develop" ];
    then
        pre_release="alpha"
    elif [ "${CI_ENVIRONMENT_NAME}" = "qa" ];
    then
        pre_release="beta"
    elif [ "${CI_ENVIRONMENT_NAME}" = "uat" ];
    then
        pre_release="rc"
    elif [ "${CI_ENVIRONMENT_NAME}" = "nonprod" ];
    then
        pre_release="rc"
    elif [ "${CI_ENVIRONMENT_NAME}" = "ops" ];
    then
        pre_release="alpha"
    fi
    echo "${pre_release}"
}

decompress()
{
    unzip -o ${CI_PROJECT_DIR}/fx-assembly/target/\*.zip -d "${BASE_DIR}"
}

dotenv()
{
    tmp_dir="/tmp/data"
    
    mkdir -p $tmp_dir
    echo "IMAGE_VERSION=${IMAGE_VERSION}" > "${tmp_dir}/.env"

    if [ -n "$CI" ];
    then
        cat "${CI_PROJECT_DIR}/.env.${CI_ENVIRONMENT_NAME}" | grep -v -e "^CI_" >> "${tmp_dir}/.env"
    else
        cat "${CI_PROJECT_DIR}/.env.${CI_ENVIRONMENT_NAME}" >> "${tmp_dir}/.env"
    fi
    
    cat "${tmp_dir}/.env"
}

exportall()
{
    tmp_dir="/tmp/data"

    set -a
    . "${tmp_dir}/.env"
    set +a
}

read_passwd()
{
  stty -echo
  printf "${1}: "
  read ${1}
  printf "\n"
  stty echo
}

mk_build()
{
    build_dir="${SUB_PROJECT_DIR}/$1"
    shift

    rm -rf "$build_dir"
    mkdir -p "$build_dir"

    for sub in "$@"; do
        src_dir="${CI_PROJECT_DIR}/$sub"

        if [ ! -d "$src_dir" ]; then
            echo "Error: directory '$src_dir' not found" >&2
            return 1
        fi

        cp -R "$src_dir"/. "$build_dir/"
    done
}
