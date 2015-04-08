#!/bin/sh
exec 2>&1

cd /home/ubuntu
exec setuidgid ubuntu ./HBserver \
    --redir-host="https://filippo.io/Heartbleed" \
    --listen=":8080" --expiry="6h"
