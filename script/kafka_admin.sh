#!/usr/bin/env bash

scriptdir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" > /dev/null && pwd )"
source $scriptdir/common.sh

TOPIC_NAME="datagather"
${KAFKA_HOME}/bin/kafka-topics.sh \
        --bootstrap-server localhost:9092 \
        --list | \
        grep -q ${TOPIC_NAME}

if [[ $? -ne 0 ]]; then
        echo "creating topic: ${TOPIC_NAME}"
        ${KAFKA_HOME}/bin/kafka-topics.sh \
                --bootstrap-server localhost:9092 \
                --create \
                --replication-factor 1 \
                --partitions 1 \
                --topic ${TOPIC_NAME}
fi
