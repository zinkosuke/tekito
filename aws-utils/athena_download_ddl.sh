#!/bin/bash
set -euo pipefail
cd "$(dirname "${0}")"

function usage() {
    cat << EOF
Usage: ./athena_download_ddl.sh

DDLをローカルに持ってくるやつ

Environments:
    TABLE_LIST      対象テーブルが書いてあるリスト
    OUTPUT_LOCATION OutputLocation
    SLEEP_SEC       get-query-executionのインターバル(秒) default: 5
EOF
    exit 1
}

required_environments=("TABLE_LIST" "OUTPUT_LOCATION")
for i in "${required_environments[@]}"; do
    [[ -v ${i} ]] || usage
done

# ----- Main -----
SLEEP_SEC=${SLEEP_SEC:-5}

grep -v '^#' < "${TABLE_LIST}" | while IFS= read -r table; do
    query_id=$(
        aws athena start-query-execution \
        --result-configuration "OutputLocation=${OUTPUT_LOCATION}" \
        --query-string "SHOW CREATE TABLE \`${table}\`;" \
        --query QueryExecutionId --output text
    )
    echo "${table}: ${query_id}"
    while :
    do
        sleep "${SLEEP_SEC}"
        query_execution=$(aws athena get-query-execution --query-execution-id "${query_id}")
        query_state=$(echo "${query_execution}" | jq -r '.QueryExecution.Status.State')
        # 'QUEUED'|'RUNNING'|'SUCCEEDED'|'FAILED'|'CANCELLED'
        case ${query_state} in
            CANCELLED|FAILED)
                echo "ERROR!"
                break
                ;;
            SUCCEEDED)
                sql_name="${table//\.//}.sql"
                mkdir -p "$(dirname "${sql_name}")"
                query_output=$(echo "${query_execution}" | jq -r '.QueryExecution.ResultConfiguration.OutputLocation')
                aws s3 cp "${query_output}" "${sql_name}"
                break
                ;;
            *) ;;
        esac
    done
done
