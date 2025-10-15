#!/bin/sh


k8s_template() {
  template="$1"
  image_version="$2"

  if [ -z "$template" ] || [ -z "$image_version" ]; then
    echo "Usage: k8s_template <template> <image_version>" >&2
    return 1
  fi

  base_dir="${BASE_DIR:-.}"
  file="${base_dir}/k8s.${template}.yaml"

  if [ ! -f "$file" ]; then
    echo "Error: Template file not found: $file" >&2
    return 1
  fi

  sed "s|\${IMAGE_VERSION}|${image_version}|g" "$file" | jq --raw-input --slurp .
}
