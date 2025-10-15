#!/bin/sh


npm_build() {
    src_dir="${CI_PROJECT_DIR}/${1}/"
    build_dir="${SUB_PROJECT_DIR}/build/${2}/"

    if [ ! -d "$src_dir" ]; then
        echo "Error: directory '$src_dir' not found" >&2
        return 1
    fi

    if [ ! -d "$build_dir" ]; then
        echo "Error: directory '$build_dir' not found" >&2
        return 1
    fi

    echo "Running npm install and Vite build in $src_dir"
    (cd "$src_dir" && npm install && npx vite build --outDir "$build_dir") || {
        echo "Error: failed to build React project with Vite"
        return 1
    }
}
