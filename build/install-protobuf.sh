#!/usr/bin/env sh
wget https://github.com/google/protobuf/archive/v$PROTOBUF_VERSION.tar.gz
tar -xzvf v$PROTOBUF_VERSION.tar.gz
cd protobuf-$PROTOBUF_VERSION && ./configure --prefix=/usr && make && sudo make install
