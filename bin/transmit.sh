#!/bin/bash

sudo /usr/local/bin/bto_ir_cmd -e -t "$1" >/dev/null 2>&1

touch "newer"
curl -s "${TASKER_WEBHOOK}"
while :
do
  FILENAME=$(find . -type f -newer newer)
  if [ $(echo "${FILENAME}" | wc -l) -eq 1 ]; then
    break
  fi
  sleep 1
done
rm -f "newer"

echo "${FILENAME}"
