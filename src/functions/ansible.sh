#!/bin/sh


ansible()
{
    playbook="${1}"
    inventory="${2}"
    run_env="${3}"
    shift 3

    if [ "${run_env}" = "shell" ];
    then
      ansible-playbook "${BASE_DIR}/ansible/playbooks/${playbook}-playbook.yaml" \
          -i "${BASE_DIR}/ansible/inventories/${inventory}-inventory.yaml" \
          -e "ansible_user=${EC2_USER}"
    elif [ "${run_env}" = "docker" ];
    then
      echo "ansible-playbook" \
           "${BASE_DIR}/playbooks/${playbook}-playbook.yaml" \
           "-i ${BASE_DIR}/inventories/${inventory}-inventory.yaml" \
           $@
    else
      echo "ansible: invalid run_env" && exit 1
    fi
}

ansible_collection_init()
{
  namespace="${1}"
  name="${2}"

  ansible-galaxy collection init "${namespace}.${name}"
}

ansible_role_init()
{
  name="${1}"

  ansible-galaxy role init ${name}
}