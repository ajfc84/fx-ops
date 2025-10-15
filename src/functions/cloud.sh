#!/bin/sh


start_imds()
{
    config_dir="${1}"

    python3 -m http.server --directory "${config_dir}"
}

cloud_init_validation()
{
    cloud-init schema --config-file user-data
}