#!/usr/bin/env bash

scriptdir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" > /dev/null && pwd )"
source $scriptdir/common.sh

${KAFKA_HOME}/bin/zookeeper-server-start.sh \
        ${KAFKA_HOME}/config/zookeeper.properties
