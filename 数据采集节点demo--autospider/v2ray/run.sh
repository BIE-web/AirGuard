#!/bin/sh
nohup /usr/bin/v2ray -c /config.json &
/gocron-node -allow-root