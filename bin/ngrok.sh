#!/bin/bash

nohup ngrok http 1323 -region=ap -log=stdout >ngrok.log &
sleep 60
curl localhost:4040/status | grep ngrok.io
