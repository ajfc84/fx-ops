#!/bin/sh


openssl_private_key()
{
  private_key_filename=${1}

  openssl genrsa -out "${private_key_filename}.pem"
}

openssl_public_key()
{
  private_key_filename=${1}
  public_key_filename=${2}

  openssl rsa -in "${private_key_filename}.pem" -out "${public_key_filename}.pem" -pubout
}

openssl_encrypt()
{
  public_key_file="${1}"
  file_to_encrypt="${2}"

  openssl pkeyutl -encrypt -pubin -inkey "${public_key_file}" -in "${file_to_encrypt}" -out "${file_to_encrypt}.enc"
}

openssl_decrypt()
{
  private_key_file="${1}"
  file_to_decrypt="${2}"

  openssl pkeyutl -decrypt -inkey "${private_key_file}" -in "${file_to_decrypt}" -out "${file_to_encrypt}.dec"
}
