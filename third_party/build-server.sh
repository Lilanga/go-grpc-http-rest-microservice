#!/bin/bash
if [[ $1 == *"http"* ]]
then
        rm -rfv ../dist/
        mkdir -p ../dist/

        sh protoc-gen.sh
        cp -R ./swagger-ui/. ../dist/
    
        cd ../cmd/server-http/
        go build -o ../../dist/ .
else
    if [[ $1 == *"grpc"* ]]
    then
        rm -rfv ../dist/
        mkdir -p ../dist/
        
        cd ../cmd/server-grpc/
        go build -o ../../dist/ .
    else
        rm -rfv ../dist/
        mkdir -p ../dist/

        sh protoc-gen.sh
        cp -R ./swagger-ui/. ../dist/

        cd ../cmd/server/
        go build -o ../../dist/ .
    fi
fi