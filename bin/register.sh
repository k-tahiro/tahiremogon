#!/bin/bash

: ${SERVER_HOST:="192.168.10.81"}

if [ $# -ne 2 ]; then
  echo "register.sh <冷房 or 暖房> <希望温度> の形式で呼んでね(・∀・)" 2>&1
  exit 1
fi

CMD_NAME="$1_$2"
if [ "$1" == "冷房" ]; then
  CMD_ID="c$2"
elif [ "$1" == "暖房" ]; then
  CMD_ID="w$2"
else
  echo "冷房なのか暖房なのかはっきりしてね(・∀・)" 1>&2
  exit 1
fi

curl "${SERVER_HOST}:1323/commands/receive" \
  -H 'Content-Type: application/json' \
  -d "{\"id\": \"${CMD_ID}\", \"name\": \"${CMD_NAME}\"}"
