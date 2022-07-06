#!/bin/bash

CONFIG_FILE=/app/app.yml
until [ $# -eq 0 ]
do
 case "$1" in
 --redis-addr)
    sed -i "4s/host.docker.internal:6379/$2/g" $CONFIG_FILE
    shift 2;;
 --redis-password)
    sed -i "5s/none/$2/g" $CONFIG_FILE
    shift 2;;
 --redis-db)
    sed -i "6s/0/$2/g" $CONFIG_FILE
    shift 2;;
 --redis-pool)
    sed -i "7s/100/$2/g" $CONFIG_FILE
    shift 2;;
 *) echo " unknow prop $1";shift;;
 esac
done

echo "============app.js==============="
cat $CONFIG_FILE
echo "===================================="

./app