#!/bin/sh

BuildModule() {
    go build -o backend
}

RunModule(){
    ./backend $@
    exit_status=$?
    if [ $exit_status -eq 1 ]; then
        exit $exit_status
    fi
    exit $exit_status
}


Start(){
    BuildModule
    while true
    do
        if ping -c 1 db_migrate &> /dev/null
        then
            echo "Still waiting db_migrate"
        else
            echo "db_migrate done,let's fvcking goooooo!!!!!"
            RunModule
        fi
        sleep 30
    done
}

Start