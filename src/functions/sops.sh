#!/bin/sh


sops_read()
{
    FILE="$1"

    [ -f "$FILE" ] || { echo "Usage: sops_read <file>" >&2; return 1; }
    grep -q '^sops:' "$FILE" || { echo "Error: '$FILE' is not encrypted with SOPS" >&2; return 1; }

    eval "$(
        sops -d "$FILE" \
        | yq -r --arg env "${CI_ENVIRONMENT_NAME}" '
            .secrets as $s
            | ($s.global // {}) * ($s[$env] // {})
            | to_entries[]
            | select(.value != null)
            | "export " + .key + "=" + ((.value | tostring) | @sh)
        '
    )" || { echo "Error: exporting secrets failed" >&2; return 1; }

    sops -d "$FILE" \
    | yq -r --arg env "${CI_ENVIRONMENT_NAME}" '
        .secrets as $s
        | ($s.global // {}) * ($s[$env] // {})
        | to_entries[]
        | select(.value != null)
        | .key + "=*****"
    '
}

sops_dec_enc()
{
    main_file="$1"

    if grep -q '^sops:' "$main_file"; then
        echo "Decrypting '$main_file'..."
        sops -d -i "$main_file" || {
            echo "Error: Failed to decrypt '$main_file'" >&2
            return 1
        }
    else
        echo "Encrypting '$main_file'..."
        sops -e -i "$main_file" || {
            echo "Error: Failed to encrypt '$main_file'" >&2
            return 1
        }
    fi
}