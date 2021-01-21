#!/bin/sh

BuildModule() {
    go build -o livestream
}

RunModule(){
    ./livestream $@
    exit_status=$?
    if [ $exit_status -eq 1 ]; then
        exit $exit_status
    fi
    exit $exit_status
}


Start(){
    export VERSION=$(git tag | tail -n1)
    BuildModule
    RunModule $@
}

Start $@