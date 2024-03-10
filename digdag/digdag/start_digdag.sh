#!/bin/bash
set -euo pipefail

CONFIG_PATH=/home/digdag/server.properties

function usage() {
    cat << EOF
Usage: ./start_digdag.sh COMMAND

Digdag start server.

Commands:
    web       API server
    scheduler API server, workflow executor, schedule executor
    agent     API server, agent
    all       API server, workflow executor, schedule executor, agent
EOF
    exit 1
}

function build_server_conf() {
    cat << EOF > "${CONFIG_PATH}"
server.bind = 0.0.0.0
server.port = 65432
server.admin.bind = 0.0.0.0
server.admin.port = 65433
# server.access-log.path = /var/log/digdag/server
# server.access-log.pattern = json
server.http.io-threads = 2
server.http.worker-threads = 16
server.http.enable-http2 = true

database.type = postgresql
database.user = ${DB_USERNAME}
database.password = ${DB_PASSWORD}
database.host = ${DB_HOST}
database.port = ${DB_PORT}
database.database = ${DB_DATABASE}
database.maximumPoolSize = 32

log-server.type = s3
log-server.s3.bucket = ${LOG_S3_BUCKET}
log-server.s3.path = ${LOG_S3_PREFIX}/log
log-server.s3.direct_download = false
log-server.s3.path-style-access = true

archive.type = s3
archive.s3.bucket = ${LOG_S3_BUCKET}
archive.s3.path = ${LOG_S3_PREFIX}/archive
archive.s3.path-style-access = true

digdag.secret-encryption-key = ${DIGDAG_SECRET_KEY}
EOF
    if [[ ${LOG_S3_ENDPOINT:-x} != "x" ]]; then
        echo "log-server.s3.endpoint = ${LOG_S3_ENDPOINT}" >> "${CONFIG_PATH}"
        echo "archive.s3.endpoint = ${LOG_S3_ENDPOINT}" >> "${CONFIG_PATH}"
    fi
}

case "${1}" in
    web)
        build_server_conf
        exec java -jar /usr/local/bin/digdag server \
            --disable-executor-loop \
            --disable-local-agent \
            --disable-scheduler \
            -c "${CONFIG_PATH}"
        ;;
    scheduler)
        build_server_conf
        exec java -jar /usr/local/bin/digdag server \
            --disable-local-agent \
            -c "${CONFIG_PATH}"
        ;;
    agent)
        build_server_conf
        exec java -jar /usr/local/bin/digdag server \
            --disable-executor-loop \
            --disable-scheduler \
            -c "${CONFIG_PATH}"
        ;;
    all)
        build_server_conf
        exec java -jar /usr/local/bin/digdag server \
            -c "${CONFIG_PATH}"
        ;;
    test-conf)
        CONFIG_PATH=/dev/stdout build_server_conf
        exit 0
        ;;
    *)
        usage
        ;;
esac
