package main

import (
	"github.com/glory-go/glory/log"
	"github.com/online-im/go-sdk/pkg/message"
	"github.com/online-im/go-sdk/pkg/client"
	"time"
)

func main() {
	go func() {
		client, err := client.NewIMGoClientWithK8sServiceName("online-im-instance", "123")
		if err != nil {
			panic(err)
		}
		ch := client.GetRecvChan()
		if err := client.SendMsg(&message.UserMessage{
			Data:     "hello laurence1",
			TargetID: "123",
			Type:     1,
		}); err != nil {
			panic(err)
		}
		log.Info("send message success")

		for {
			select {
			case msg := <-ch:
				log.Info("client 1 recv message data = ", msg.Data)
			}
		}
	}()

	time.Sleep(time.Second)
	go func() {
		client, err := client.NewIMGoClientWithK8sServiceName("online-im-instance", "456")
		if err != nil {
			panic(err)
		}
		ch := client.GetRecvChan()
		if err := client.SendMsg(&message.UserMessage{
			Data:     "hello laurence1",
			TargetID: "123",
			Type:     1,
		}); err != nil {
			panic(err)
		}
		log.Info("send message success")

		for {
			select {
			case msg := <-ch:
				log.Info("client 2 recv message data = ", msg.Data)
			}
		}
	}()

	time.Sleep(time.Second)

	client, err := client.NewIMGoClientWithK8sServiceName("online-im-instance", "789")
	if err != nil {
		panic(err)
	}
	for {
		if err := client.SendMsg(&message.UserMessage{
			Data:     "hello laurence1",
			TargetID: "123",
			Type:     1,
		}); err != nil {
			panic(err)
		}
		//log.Info("send message success")

		if err := client.SendMsg(&message.UserMessage{
			Data:     "hello laurence2",
			TargetID: "456",
			Type:     1,
		}); err != nil {
			panic(err)
		}
		time.Sleep(time.Second)
	}
}
