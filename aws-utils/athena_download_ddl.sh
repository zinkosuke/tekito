#!/bin/bash
#
# DDLをローカルに持ってくるやつ
#
set -euo pipefail
# ----- Environments -----
TABLE_LIST=${TABLE_LIST}
OUTPUT_LOCATION=${OUTPUT_LOCATION}
SLEEP_SEC=${SLEEP_SEC:-5}
# ----- Args -----
# ----- Main -----
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
                output_location=$(echo "${query_execution}" | jq -r '.QueryExecution.ResultConfiguration.OutputLocation')
                aws s3 cp "${output_location}" "${sql_name}"
                break
                ;;
            *) ;;
        esac
    done
done
