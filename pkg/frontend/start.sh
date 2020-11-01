#!/bin/sh

while true
do
    if ping -c 1 db_migrate &> /dev/null
    then
        echo "Still waiting db_migrate"
    else
        echo "db_migrate done,let's fvcking goooooo!!!!!"
        go run .
    fi
    sleep 30
done