#!/bin/bash
set -euo pipefail

# shellcheck disable=SC2016
envsubst '$$DIGDAG_HOST' < /tmp/nginx.conf > /etc/nginx/nginx.conf

exec nginx -g "daemon off;"
