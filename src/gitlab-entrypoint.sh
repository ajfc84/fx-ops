#!/bin/sh

if [ -n "$CI" ];
then
    # this is how GitLab expects the entrypoint to end, if provided
    # will execute scripts from stdin
    exec /bin/bash
else
    exec pipeline.sh "$@"
fi
