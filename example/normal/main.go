package main

import (
	"git.go-online.org.cn/Glory/glory/log"
	go_sdk "github.com/online-im/go-sdk"
	"github.com/online-im/go-sdk/pkg"
	"time"
)

func main() {
	go func() {
		client, err := go_sdk.NewIMGoClientWithGatewayAddr("localhost:8086", "123")
		if err != nil {
			panic(err)
		}
		ch := client.GetRecvChan()
		if err := client.SendMsg(&pkg.Message{
			Data:     "hello laurence1",
			FromID:   "123",
			TargetID: "123",
			Type:     1,
		}); err != nil {
			panic(err)
		}
		log.Info("send message success")

		select {
		case msg := <-ch:
			log.Info("recv message data = ", msg.Data)
		}
	}()

	time.Sleep(time.Second)
	go func() {
		client, err := go_sdk.NewIMGoClientWithGatewayAddr("localhost:8086", "456")
		if err != nil {
			panic(err)
		}
		ch := client.GetRecvChan()
		if err := client.SendMsg(&pkg.Message{
			Data:     "hello laurence1",
			FromID:   "456",
			TargetID: "123",
			Type:     1,
		}); err != nil {
			panic(err)
		}
		log.Info("send message success")

		select {
		case msg := <-ch:
			log.Info("recv message data = ", msg.Data)
		}
	}()

	time.Sleep(time.Second)
	client, err := go_sdk.NewIMGoClientWithGatewayAddr("localhost:8086", "789")
	if err != nil {
		panic(err)
	}
	if err := client.SendMsg(&pkg.Message{
		Data:     "hello laurence1",
		FromID:   "789",
		TargetID: "123",
		Type:     1,
	}); err != nil {
		panic(err)
	}
	log.Info("send message success")

	if err := client.SendMsg(&pkg.Message{
		Data:     "hello laurence2",
		FromID:   "789",
		TargetID: "456",
		Type:     1,
	}); err != nil {
		panic(err)
	}
	log.Info("send message success")
	select {}

}