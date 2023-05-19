#!/bin/bash
set -euo pipefail
cd "$(dirname "${0}")/../"

function usage() {
    cat << EOF
Usage: ./scripts/push_workflows.sh

Digdag push workflows
EOF
    exit 1
}

revision=$(date +%Y-%m-%dT%H:%M:%S)

for project_dir in projects/*; do
    echo "Push project '${project_dir}'"
    project=${project_dir##*/}
    digdag push "${project}" \
        --project "${project_dir}" \
        --revision "${revision}" \
        --copy-outgoing-symlinks
done
