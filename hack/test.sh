#!/bin/sh

mkdir -p ./bin
go build -o ./bin/ddns
./bin/ddns \
-token="24000,0c4f454d3c7bcdef0803abcdef" \
-domain="example.com" \
-record="test.ddns" \
-interval=10

rm -rf ./bin