#!/usr/bin/env sh
git clone https://github.com/protocolbuffers/protobuf.git
cd protobuf \
        git checkout v$PROTOBUF_VERSION && \
        git submodule update --init --recursive && \
        ./autogen.sh && \
        ./configure && \
        make && \
        make check && \
        sudo make install && \
        sudo ldconfig
