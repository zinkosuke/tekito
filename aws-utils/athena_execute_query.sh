#!/bin/bash
set -euo pipefail
cd "$(dirname "${0}")"

usage() {
    cat << EOF
Usage: ./athena_download_ddl.sh query

Athenaのクエリ実行

Environments:
    SLEEP_SEC    get-query-executionのインターバル(秒) default: 5
    WORK_GROUP   WorkGroup
    SHOW_RESULTS 結果をstdout表示(Y/n) default: Y
EOF
    exit 1
}


##################################################
# Args / Environments
##################################################
[[ ${WORK_GROUP:-} = "" ]] && usage
SLEEP_SEC=${SLEEP_SEC:-5}

QUERY=${1:-}
[[ ${QUERY} = "" ]] && usage


##################################################
# Signal
##################################################
signal_handler() {
    if [[ ${query_id:-} != "" ]]; then
        echo "StopQueryExecution" 1>&2
        aws athena stop-query-execution --query-execution-id "${query_id}"
    fi
}
trap 'signal_handler' SIGINT
trap 'signal_handler' SIGTERM


##################################################
# Main
##################################################
echo "StartQueryExecution" 1>&2
query_id=$(
    aws athena start-query-execution \
    --work-group "${WORK_GROUP}" \
    --query-string "${QUERY}" \
    --query QueryExecutionId --output text
)

echo "QueryId: ${query_id}" 1>&2
while :
do
    sleep "${SLEEP_SEC}"

    query_execution=$(aws athena get-query-execution --query-execution-id "${query_id}")
    query_state=$(echo "${query_execution}" | jq -r ".QueryExecution.Status.State")

    # 'QUEUED'|'RUNNING'|'SUCCEEDED'|'FAILED'|'CANCELLED'
    case ${query_state} in
        CANCELLED|FAILED)
            echo "${query_state}" 1>&2
            break
            ;;
        SUCCEEDED)
            query_output=$(
                echo "${query_execution}" \
                | jq -r ".QueryExecution.ResultConfiguration.OutputLocation"
            )
            aws s3 cp "${query_output}" -
            break
            ;;
        *) ;;
    esac
done
