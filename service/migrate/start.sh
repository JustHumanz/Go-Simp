#!/bin/sh

while true
do
    if ping -c 1 db_migrate &> /dev/null
    then
        echo "Main bot still up"
        exit 0 
    else
        echo "mainbot is down,let's fvcking goooooo!!!!!"
        go run .

    fi
done