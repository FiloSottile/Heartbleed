#!/bin/sh

cd /home/ubuntu

exec sudo ./heartbleed 2>&1 | tee -a ./heartbleed.log

