#!/usr/bin/env bash

set -e

docker-compose -f basic/docker-compose.yml up -d

sleep 60

docker exec broker \
    kafka-topics --bootstrap-server broker:9092 \
        --create \
        --partitions 3 \
        --topic quickstart

