#!/usr/bin/env bash

scriptdir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" > /dev/null && pwd )"
source $scriptdir/common.sh

grpcurl -plaintext \
        -import-path ${scriptdir}/../internal/service \
        -import-path $GOPATH/src \
        -import-path ${scriptdir}/../vendor \
        -proto datagather_service.proto \
        -d '{"user_id": "u1", "event": {"user_id": "u1", "item_id": "i1"}}' \
        localhost:9091 \
        farseer.datagather.DatagatherService/CreateEvent

#list farseer.datagather.DatagatherService
#farseer.datagather.DatagatherService.CreateEvent
