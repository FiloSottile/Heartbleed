#!/bin/sh

cd /home/ubuntu

exec sudo ./HBserver --redir-host="https://filippo.io/Heartbleed" --listen=":80" --expiry="1h" 2>&1 | tee -a ./heartbleed.log

