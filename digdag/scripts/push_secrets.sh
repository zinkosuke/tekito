#!/bin/bash
set -euo pipefail
cd "$(dirname "${0}")/../"

cd .secrets

for secret_file in *.json; do
    echo "Push secret '${secret_file}'"
    project=${secret_file%%.*}
    digdag secrets \
        --project "${project}" \
        --set "@${secret_file}"
done
