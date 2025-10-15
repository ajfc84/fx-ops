#!/bin/sh


cron_logs()
{
    grep CRON /var/log/syslog
}

tcp_list_interfaces_with_status()
{
    tcpdump -D
}

tcp_capture_packets()
{
    tcpdump --interface $1
}

logical_volume_list()
{
    lvs
}

volume_goup_list()
{
    vgs
}

logical_volume_extend()
{
    size="${1}"
    vol_grp_path="${2}"

    lvextend -L ${size} ${vol_grp_path}
}