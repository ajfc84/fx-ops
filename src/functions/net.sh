#!/bin/sh


net_install()
{
    curl -s https://install.zerotier.com | sudo bash
}

net_join()
{
    network_id="${1}"

    zerotier-cli join "${network_id}"
}

ip_bridge_create()
{
    name="${1}"

    ip link add ${name} type bridge
}

ip_bridge_show()
{
    ip link show type bridge
}

ip_bridge_add_interface()
{
    interface="${1}"
    bridge="${2}"

    ip link set ${interface} master ${bridge}
}