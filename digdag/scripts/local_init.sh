#!/bin/bash
set -x

echo '{"test": "aaaaaaaaaaaaaa"}' > .secrets/mydag.json

# shellcheck disable=SC2016
docker-compose exec scheduler bash -c 'aws --endpoint-url "${LOG_S3_ENDPOINT}" s3api create-bucket --bucket "${LOG_S3_BUCKET}"'

docker-compose exec scheduler /tmp/scripts/push_workflows.sh
docker-compose exec scheduler /tmp/scripts/push_secrets.sh
