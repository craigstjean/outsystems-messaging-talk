#!/usr/bin/env bash

set -e

docker-compose up -d -f restproxy/docker-compose.yml

sleep 60

docker exec broker \
    kafka-topics --bootstrap-server broker:9092 \
        --create \
        --partitions 3 \
        --topic quickstart

curl -X POST -H "Content-Type: application/vnd.kafka.json.v2+json" \
      --data '{"records":[{"value":{"foo":"bar"}}]}' "http://localhost:8086/topics/jsontest"
