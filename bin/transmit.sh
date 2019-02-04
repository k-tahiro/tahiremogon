#!/bin/bash

readonly LOG_DIR="/var/log/tahiremogon"
readonly CMD_LOG_FILE="${LOG_DIR}/cmd.log"

sudo /usr/local/bin/bto_ir_cmd -e -t "$1" >>"${CMD_LOG_FILE}" 2>&1 >>"${CMD_LOG_FILE}"
FILE=$(camera.sh "${IMG_DIR}")
echo "${FILE}" >>"${CMD_LOG_FILE}"
mv "${FILE}" "/mnt/nasne/share2/tahiremocon/" >>"${CMD_LOG_FILE}"
curl "192.168.10.43:5042/detect/$(basename "${FILE}")" | python -c 'import json; import sys; d = json.loads(raw_input()); print(d["status"])' | tr -d '\n'
