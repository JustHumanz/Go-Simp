#!/bin/sh

for (( ; ; ))
do
    if ping -c 1 db_migrate &> /dev/null
    then
        echo "Still waiting db_migrate"
    else
        exit 0
    fi
    sleep 1
done