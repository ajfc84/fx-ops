#!/bin/sh


release_version() {
    changelog_file="${CI_PROJECT_DIR}/CHANGELOG"

    versions=$(grep -oP '(?<=### \*\*Version )[^*]+' "$changelog_file")
    echo "$versions" | sed -n '2p'
}
