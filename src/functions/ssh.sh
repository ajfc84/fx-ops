#!/bin/sh


ssh_keygen_rsa()
{
    ssh-keygen -t rsa
}

ssh_copy_cert()
{
  username="${1}"
  remote_host="${2}"

  ssh-copy-id "${username}@${remote_host}"
}