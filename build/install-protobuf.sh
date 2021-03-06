#!/usr/bin/env bash

set -ex
case $OSTYPE in
        darwin*)
                arch="osx-$(uname -m)"
                ;;
        linux*)
                arch="linux-$(uname -m)"
                ;;
        msys*)
                arch="win64"
                ;;
esac

wget -q https://github.com/protocolbuffers/protobuf/releases/download/v$PROTOBUF_VERSION/protoc-$PROTOBUF_VERSION-$arch.zip
sudo unzip -d /usr protoc-$PROTOBUF_VERSION-$arch.zip
rm -f protoc-$PROTOBUF_VERSION-$arch.zip
