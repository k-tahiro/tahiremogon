#!/bin/bash

nohup ngrok http 1323 -region=ap -log=stdout >ngrok.log &
sleep 60
eval echo "$(curl -s localhost:4040/status | grep ngrok.io | cut -f 2,3,4 -d '(' | cut -f 1,2,3 -d ')')" | python3 -c 'import json; print(json.loads(input())["Session"]["Tunnels"]["command_line"]["URL"])' | tr -d '\n'
