#!/usr/bin/env sh
wget https://github.com/google/protobuf/archive/v$PROTOBUF_VERSION.tar.gz
tar -xzvf v$PROTOBUF_VERSION.tar.gz
cd protobuf-$PROTOBUF_VERSION && ./autogen.sh && ./configure && make && make check && sudo make install && sudo ldconfig
