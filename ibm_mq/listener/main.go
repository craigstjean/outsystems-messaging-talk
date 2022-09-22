package main

import (
	"encoding/hex"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/ibm-messaging/mq-golang/v5/ibmmq"
)

// Hacked up version of https://github.com/ibm-messaging/mq-golang/blob/master/samples/amqsget.go

func main() {
	var qMgrObject ibmmq.MQQueueManager
	var requestQueueObject ibmmq.MQObject
	var responseQueueObject ibmmq.MQObject

	qMgrName := "QM1"
	requestQueueName := "REQ.HELLO"
	responseQueueName := "RES.HELLO"

	cno := ibmmq.NewMQCNO()
	cd := ibmmq.NewMQCD()
	cd.ChannelName = "DEV.APP.SVRCONN"
	cd.ConnectionName = "localhost(1414)"

	cno.ClientConn = cd
	cno.Options = ibmmq.MQCNO_CLIENT_BINDING
	cno.ApplName = "ibmmqlistener"

	csp := ibmmq.NewMQCSP()
	csp.AuthenticationType = ibmmq.MQCSP_AUTH_USER_ID_AND_PWD
	csp.UserId = "app"
	csp.Password = "passw0rd"
	cno.SecurityParms = csp

	qMgrObject, err := ibmmq.Connx(qMgrName, cno)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("Connected to queue manager %s\n", qMgrName)
	defer disc(qMgrObject)

	requestMqOD := ibmmq.NewMQOD()
	requestOpenOptions := ibmmq.MQOO_INPUT_EXCLUSIVE

	requestMqOD.ObjectType = ibmmq.MQOT_Q
	requestMqOD.ObjectName = requestQueueName

	requestQueueObject, err = qMgrObject.Open(requestMqOD, requestOpenOptions)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("Opened queue", requestQueueObject.Name)
	defer close(requestQueueObject)

	responseMqOD := ibmmq.NewMQOD()
	responseOpenOptions := ibmmq.MQOO_OUTPUT

	responseMqOD.ObjectType = ibmmq.MQOT_Q
	responseMqOD.ObjectName = responseQueueName

	responseQueueObject, err = qMgrObject.Open(responseMqOD, responseOpenOptions)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("Opened queue", responseQueueObject.Name)
	defer close(responseQueueObject)

	msgAvail := true
	for msgAvail == true {
		var datalen int

		getmqmd := ibmmq.NewMQMD()
		gmo := ibmmq.NewMQGMO()
		gmo.Options = ibmmq.MQGMO_NO_SYNCPOINT

		// Set options to wait for a maximum of 3 seconds for any new message to arrive
		gmo.Options |= ibmmq.MQGMO_WAIT
		gmo.WaitInterval = 3 * 1000

		buffer := make([]byte, 0, 1024)
		buffer, datalen, err = requestQueueObject.GetSlice(getmqmd, gmo, buffer)

		if err != nil {
			fmt.Println(err)
			mqret := err.(*ibmmq.MQReturn)
			if mqret.MQRC == ibmmq.MQRC_NO_MSG_AVAILABLE {
				// If there's no message available, then I won't treat that as a real error as
				// it's an expected situation
				err = nil
			}
		} else {
			fmt.Printf("Got message of length %d: ", datalen)
			fmt.Println("MsgId:" + hex.EncodeToString(getmqmd.MsgId))
			fmt.Println("CorrelId:" + hex.EncodeToString(getmqmd.CorrelId))
			fmt.Println(strings.TrimSpace(string(buffer)))

			putmqmd := ibmmq.NewMQMD()
			pmo := ibmmq.NewMQPMO()
			pmo.Options = ibmmq.MQPMO_NO_SYNCPOINT
			putmqmd.Format = ibmmq.MQFMT_STRING
			putmqmd.MsgId = getmqmd.CorrelId

			msgData := "Hello from Go at " + time.Now().Format(time.RFC3339)

			buffer := []byte(msgData)
			err = responseQueueObject.Put(putmqmd, pmo, buffer)

			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println("Put message to", strings.TrimSpace(responseQueueObject.Name))
				fmt.Println("MsgId:" + hex.EncodeToString(putmqmd.MsgId))
			}
		}
	}
}

func disc(qMgrObject ibmmq.MQQueueManager) error {
	err := qMgrObject.Disc()
	if err == nil {
		fmt.Println("Disconnected from queue manager")
	} else {
		fmt.Println(err)
	}

	return err
}

func close(object ibmmq.MQObject) error {
	err := object.Close(0)

	if err == nil {
		fmt.Println("Closed queue")
	} else {
		fmt.Println(err)
	}

	return err
}
