#!/usr/bin/env bash

grpcurl -plaintext \
        -import-path ./internal/service \
        -import-path $GOPATH/src \
        -import-path ./vendor \
        -proto datagather_service.proto \
        -d '{"user_id": "u1", "event": {"user_id": "u1", "item_id": "i1"}}' \
        localhost:9091 \
        farseer.datagather.DatagatherService/CreateEvent

#list farseer.datagather.DatagatherService
#farseer.datagather.DatagatherService.CreateEvent
