#!/bin/sh


git_tag()
{
    CWD=$(pwd)
    cd "${CI_PROJECT_DIR}"
    image_version="${1}"
    changelog=$(/ops/functions/notes.py "${CI_PROJECT_DIR:-.}/CHANGELOG" "${image_version}")

    git tag "${image_version}" -m "${changelog}"
    git push origin "${image_version}" || true # ignore error if remote does not exist
    echo "Tagged version: ${image_version}"

    cd "${CWD}"
}

git_fetch_tags()
{
    CWD=$(pwd)
    cd "${CI_PROJECT_DIR}"

    git fetch --tags || true # ignore error if remote does not exist

    cd "${CWD}"
}

git_init_bare()
{
    git init --bare
}

git_remote_add()
{
    name="${1}"
    user="${2}"
    url="${3}"
    group="${4}"
    project="${5}"

    git remote add "${name}" "${user}@${url}:${group}/${project}.git"
}

git_set_default_branch()
{
  name="${1}"

  git symbolic-ref HEAD "refs/heads/${name}"
}

git_clone()
{
  hostname="${1}"
  path="${2}/"
  project="${3}"
  port=":${4}"
  git clone "git@${hostname}${port}:${path}${project}.git"
}

git_current_branch()
{
    CWD=$(pwd)
    cd "${CI_PROJECT_DIR}"

    git branch --show-current

    cd "${CWD}"
}

git_archive
{
     git archive --remote=root@78.47.67.129:/mycicdmain/fx.git develop fx-ops | tar -x
}