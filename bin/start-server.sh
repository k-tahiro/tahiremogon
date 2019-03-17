#!/bin/bash

SCRIPT_DIR="$(cd $(dirname $0) && pwd)"

. "${SCRIPT_DIR}/../config/conf.txt"
DB_FILE="${1:-"${SCRIPT_DIR}/../command.db"}"

pushd "${SCRIPT_DIR}/../server"
docker build -t tahiro/tahiremogon .
popd
docker run -itd \
           -p 1323:1323 \
           -v "${DB_FILE}:/command.db" \
           -e "MODE=${MODE}" \
           -e "HOSTNAME=${HOSTNAME}" \
           -e "USERNAME=${USERNAME}" \
           -e "PASSWORD=${PASSWORD}" \
           --name tahiremogon \
           tahiro/tahiremogon
