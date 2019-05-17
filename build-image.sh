#!/bin/bash
set -o pipefail
set -o errexit

readonly IMAGE_TAG=${IMAGE_TAG-local}

main() {
  local readonly full_name=backup-exporter:$IMAGE_TAG

  docker build -t $full_name .
}

main "$@"

