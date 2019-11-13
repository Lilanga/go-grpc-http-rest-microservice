#!/bin/bash
protoc --proto_path=./ --proto_path=../api/proto/v1/ --go_out=plugins=grpc:../pkg/api/v1/ ../api/proto/v1/todo-service.proto
protoc --proto_path=./ --proto_path=../api/proto/v1/ --grpc-gateway_out=logtostderr=true:../pkg/api/v1/ ../api/proto/v1/todo-service.proto
protoc --proto_path=./ --proto_path=../api/proto/v1/ --swagger_out=logtostderr=true:../dist/ ../api/proto/v1/todo-service.proto