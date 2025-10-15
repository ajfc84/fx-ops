#!/bin/sh


aws_configure()
{
    if [ -n "$CI" ];
    then
        echo "Unsupported operation"
        exit 1
    else
        aws configure
    fi
}
