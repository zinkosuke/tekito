#!/bin/bash

usage() {
    cat << EOF
Usage: ./scripts/start_digdag.sh COMMAND

Digdag start server.

Commands:
    web       API server
    scheduler API server, workflow executor, and schedule executor
    agent     API server, agent
    all       API server, agent, workflow executor, and schedule executor
EOF
    exit 1
}

sigterm_handler() {
    if [[ ${pid} -ne 0 ]]; then
        kill -15 "${pid}"
        wait "${pid}"
    fi
    exit $((128 + 15));
}

trap 'sigterm_handler' SIGTERM


# TODO この辺でconfig/server.conf作るなり


case "${1}" in
    web)
        java -jar /usr/local/bin/digdag server \
            --disable-executor-loop \
            --disable-local-agent \
            --disable-scheduler \
            -c /home/digdag/config/server.conf \
            &
        pid="$!"
        ;;
    scheduler)
        java -jar /usr/local/bin/digdag server \
            --disable-local-agent \
            -c /home/digdag/config/server.conf \
            &
        pid="$!"
        ;;
    agent)
        java -jar /usr/local/bin/digdag server \
            --disable-executor-loop \
            --disable-scheduler \
            -c /home/digdag/config/server.conf \
            &
        pid="$!"
        ;;
    all)
        java -jar /usr/local/bin/digdag server \
            -c /home/digdag/config/server.conf \
            &
        pid="$!"
        ;;
    *)
        usage
        ;;
esac

wait "${pid}"
exit "$?"
