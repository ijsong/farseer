#!/usr/bin/env bash

scriptdir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" > /dev/null && pwd )"
source $scriptdir/common.sh

#{CASSANDRA_HOME}/bin/cassandra
