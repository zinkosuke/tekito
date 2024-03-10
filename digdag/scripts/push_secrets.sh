#!/bin/bash
set -euo pipefail
cd "$(dirname "${0}")/../"

usage() {
    cat << EOF
Usage: ./push_secrets.sh

Digdag push secrets
.secrets配下に {project_name}.json を準備しておく
EOF
    exit 1
}

cd /tmp/.secrets
for secret_file in *.json; do
    echo "Push secret '${secret_file}'"
    project=${secret_file%%.*}
    digdag secrets \
        --project "${project}" \
        --set "@${secret_file}"
done
