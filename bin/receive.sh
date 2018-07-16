#!/bin/bash

result="$(sudo /usr/local/bin/bto_ir_cmd -e -r)"
if [ $? -ne 0 ]; then
  exit 1
fi
echo "${result}" | tail -n 1 | cut -f 2 -d : | cut -b 2-
