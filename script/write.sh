#!/usr/bin/env bash

scriptdir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" > /dev/null && pwd )"
source $scriptdir/common.sh

input=$1
cat $input | \
        grpcurl -plaintext \
        -import-path ${scriptdir}/../internal/service \
        -import-path $GOPATH/src \
        -import-path ${scriptdir}/../vendor \
        -proto datagather_service.proto \
        -d @ \
        localhost:9091 \
        farseer.datagather.DatagatherService/CreateEvent

#-d '{"user_id": "u1", "event": {"user_id": "u1", "item_id": "i1"}}' \
#list farseer.datagather.DatagatherService
#farseer.datagather.DatagatherService.CreateEvent
