#!/usr/bin/env bash

git clone https://github.com/edenhill/librdkafka.git
cd librdkafka
./configure --prefix /usr
make
sudo make install
