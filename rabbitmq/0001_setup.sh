#!/usr/bin/env bash

set -e

docker build -t rabbitmq-demo .

docker run -d --hostname rabbitmq-demo --name rabbitmq-demo \
    -p 15672:15672 \
    -p 15674:15674 \
    -p 15675:15675 \
    -p 15670:15670 \
    -p 15671:15671 \
    -p 1883:1883 \
    -p 61613:61613 \
    -p 5672:5672 \
    rabbitmq-demo

