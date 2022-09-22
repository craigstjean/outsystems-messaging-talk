#!/usr/bin/env bash

set -e

docker pull icr.io/ibm-messaging/mq:latest

docker volume create qm1data
docker run --env LICENSE=accept --env MQ_QMGR_NAME=QM1 --volume qm1data:$PWD/mqmdata --publish 1414:1414 --publish 9443:9443 --detach --env MQ_ADMIN_PASSWORD=passw0rd --env MQ_APP_PASSWORD=passw0rd --name QM1 icr.io/ibm-messaging/mq:latest

sleep 60

docker exec -ti QM1 dspmqver
docker exec -ti QM1 dspmqweb status

# Create queue MSGQ
curl -k https://localhost:9443/ibmmq/rest/v1/admin/qmgr/QM1/queue -X POST -u admin:passw0rd -H "ibm-mq-rest-csrf-token: value" -H "Content-Type: application/json" --data "{\"name\":\"MSGQ\"}"

# Set permissions to user 'app'
# echo "SET AUTHREC PROFILE(MSGQ) OBJTYPE(QUEUE) +\\nPRINCIPAL(app) AUTHADD(BROWSE, INQ, GET, PUT)\\nend" | socat EXEC:"docker exec -ti QM1 runmqsc QM1",pty STDIN
docker exec -ti QM1 setmqaut -m QM1 -n MSGQ -t queue -p app +inq +get +put +browse

# Write message to queue
curl -k https://localhost:9443/ibmmq/rest/v1/messaging/qmgr/QM1/queue/MSGQ/message -X POST -u app:passw0rd -H "ibm-mq-rest-csrf-token: value" -H "Content-Type: text/plain;charset=utf-8" --data "Hello World"

# Read and delete message from queue
curl -k https://localhost:9443/ibmmq/rest/v1/messaging/qmgr/QM1/queue/MSGQ/message -X DELETE -u app:passw0rd -H "ibm-mq-rest-csrf-token: value"

# Create request and response queues
curl -k https://localhost:9443/ibmmq/rest/v1/admin/qmgr/QM1/queue -X POST -u admin:passw0rd -H "ibm-mq-rest-csrf-token: value" -H "Content-Type: application/json" --data "{\"name\":\"REQ.HELLO\"}"
curl -k https://localhost:9443/ibmmq/rest/v1/admin/qmgr/QM1/queue -X POST -u admin:passw0rd -H "ibm-mq-rest-csrf-token: value" -H "Content-Type: application/json" --data "{\"name\":\"RES.HELLO\"}"

# Set permissions to user 'app'
docker exec -ti QM1 setmqaut -m QM1 -n REQ.HELLO -t queue -p app +inq +get +put +browse
docker exec -ti QM1 setmqaut -m QM1 -n RES.HELLO -t queue -p app +inq +get +put +browse

