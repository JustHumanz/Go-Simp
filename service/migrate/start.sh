#!/bin/sh

while true
do
    if ping -c 1 livestream &> /dev/null || ping -c 1 fanart &> /dev/null
    then
        echo "Some module still up"
        exit 0 
    else
        go run .

    fi
done