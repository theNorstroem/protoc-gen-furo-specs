#!/usr/bin/env bash

go build .

cd in

protoc --plugin ../protoc-gen-furo-specs \
-I./in \
-I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
-I/usr/local/include \
-I$GOPATH/src \
-I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
-I$GOPATH/src/github.com/googleapis/googleapis \
-I./ \
--furo-specs_out=\
Mhelloworld.proto=../helloworld,\
:../../out ./helloworld/*.proto
