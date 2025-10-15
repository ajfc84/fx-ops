#!/bin/sh


kvm()
{
    # Ubuntu 22.04
    apt update
    result=$(egrep -c '(vmx|svm)' /proc/cpuinfo)
    if [ ${result} -ge 0 ];
    then
        echo "Virtualization enabled"
    else
        echo "Enable virtualization in BIOS"
    fi
    apt install -y cpu-checker
    result=$(kvm-ok)
    case "${result}" in
        *"can be used")
            echo "KVM enabled"
            ;;
        *)
            echo "Enable KVM"
            ;;
    esac
    echo """
    Installing...
    1. qemu-kvm -> An opensource emulator and virtualization package that provides hardware emulation.
    2. virt-manager -> A Qt-based graphical interface for managing virtual machines via the libvirt daemon.
    3. libvirt-daemon-system -> A package that provides configuration files required to run the libvirt daemon.
    4. virtinst -> A  set of command-line utilities for provisioning and modifying virtual machines.
    5. libvirt-clients -> A set of client-side libraries and APIs for managing and controlling virtual machines & hypervisors from the command line.
    6. bridge-utils -> A set of tools for creating and managing bridge devices.
    """
    apt install -y qemu-kvm virt-manager libvirt-daemon-system virtinst libvirt-clients bridge-utils
    systemctl enable --now libvirtd
    systemctl start libvirtd
    systemctl status libvirtd
    usermod -aG kvm $USER
    usermod -aG libvirt $USER
    echo "logout/login"
    wget https://cloud-images.ubuntu.com/jammy/current/jammy-server-cloudimg-amd64-disk-kvm.img
    qemu-img info https://cloud-images.ubuntu.com/jammy/current/jammy-server-cloudimg-amd64-disk-kvm.img
}

vm_install()
{
    virt-install --name control001 \
        --ram 1024 \
        --vcpus=1 \
        --import \
        --disk path=~/ubuntu-server.img,format=qcow2 \
        --os-variant ubuntu22.04 \
        --network bridge=virbr0,model=virtio \
        --graphics none \
        --noautoconsole \
        --cloud-init ds=nocloud;s=http://192.168.18.3/
}

vm_list()
{
    virsh list --all
}

vm_os_list()
{
    virt-install --osinfo list
}

vm_xml()
{
    name="${1}"

    virsh dumpxml ${name}
}

vm_addr()
{
    name="${1}"

    virsh domifaddr ${name}
}

vm_edit_live()
{
    name="${1}"

    virsh edit ${name}
}

vm_delete()
{
    name="${1}"

    virsh undefine ${name} --remove-all-storage
}