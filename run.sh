#!/bin/bash
set -x

rm -rf staging/src
cp -r /home/xd/go1.13.9/src/github.com/xlab-uiuc/sonar-lib/ staging/src
cd docker/cassandra-operator
make
docker tag cassandra-operator:latest xudongs/cassandra-operator:latest
docker push xudongs/cassandra-operator:latest
