#!/bin/bash
set -x

echo '{"test": "aaaaaaaaaaaaaa"}' > .secrets/mydag.json

docker-compose exec \
    -e AWS_ACCESS_KEY_ID=root \
    -e AWS_SECRET_ACCESS_KEY=password \
    web aws --endpoint-url http://minio:9000 \
    s3api create-bucket --bucket digdag-log

docker-compose exec web ./scripts/push_workflows.sh
docker-compose exec web ./scripts/push_secrets.sh
