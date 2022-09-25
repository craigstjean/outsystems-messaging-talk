#!/usr/bin/env bash

set -e

docker-compose -f bridge/docker-compose.yml up -d

sleep 60

docker exec broker2 \
    kafka-topics --bootstrap-server broker2:9092 \
        --create \
        --partitions 3 \
        --topic quickstart

curl -X POST \
  http://localhost:9080/topics/quickstart \
  -H 'content-type: application/vnd.kafka.json.v2+json' \
  -d '{
    "records": [
        {
            "key": "my-key",
            "value": "sales-lead-0001"
        },
        {
            "value": "sales-lead-0002",
            "partition": 2
        },
        {
            "value": "sales-lead-0003"
        }
    ]
}'

