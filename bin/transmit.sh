#!/bin/bash

readonly LOG_DIR="/var/log/tahiremogon"
readonly CMD_LOG_FILE="${LOG_DIR}/cmd.log"

sudo /usr/local/bin/bto_ir_cmd -e -t "$1" >>"${CMD_LOG_FILE}" 2>&1
