#!/bin/sh
mkdir ./main
chown -R ubuntu:ubuntu ./main

exec setuidgid ubuntu multilog \
    t s10485760 n500 '!tai64nlocal' '!gzip' ./main
