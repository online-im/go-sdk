package main

import (
	"git.go-online.org.cn/Glory/glory/log"
	go_sdk "github.com/online-im/go-sdk"
	"github.com/online-im/go-sdk/pkg"
)

func main() {
	client, err := go_sdk.NewIMGoClient("123")
	if err != nil {
		panic(err)
	}
	ch := client.GetRecvChan()
	if err := client.SendMsg(&pkg.Message{
		Data:     "hello laurence",
		FromID:   "123",
		TargetID: "456",
		Type:     1,
	}); err != nil {
		panic(err)
	}
	log.Info("send message success")

	select {
	case msg := <-ch:
		log.Info("recv message data = ", msg.Data)
	}
}
