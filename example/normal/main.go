package main

import (
	"fmt"
	"github.com/glory-go/glory/log"
	pClient "github.com/online-im/go-sdk/pkg/client"
	"github.com/online-im/go-sdk/pkg/constant"
	pMessage "github.com/online-im/go-sdk/pkg/message"
	"time"
)

func main() {
	go func() {
		for {
			client, err := pClient.NewIMGoClientWithInstanceAddr("localhost:8080", "123")
			if err != nil {
				panic(err)
			}
			ch := client.GetRecvChan()
			//if err := client.SendMsg(&pMessage.UserMessage{
			//	Data:        "hello laurence1",
			//	TargetID:    "123",
			//	PublishType: 1,
			//}); err != nil {
			//	panic(err)
			//}
			//log.Info("send message success")
			for {
				select {
				case msg := <-ch:
					log.Info("client 1 recv message data = ", msg.Data)
				}

			}

			//
			//select {
			//case msg := <-ch:
			//	log.Info("client 1 recv message data = ", msg.Data)
			//}
			//
			//select {
			//case msg := <-ch:
			//	log.Info("client 1 recv message data = ", msg.Data)
			//}
			//
			//select {
			//case msg := <-ch:
			//	log.Info("client 1 recv message data = ", msg.Data)
			//}
			//
			////client.Close()
		}

	}()

	time.Sleep(time.Second)
	go func() {
		client, err := pClient.NewIMGoClientWithInstanceAddr("localhost:8080", "456")
		if err != nil {
			panic(err)
		}
		ch := client.GetRecvChan()
		//if err := client.SendMsg(&pMessage.UserMessage{
		//	Data:     "hello laurence1",
		//	TargetID: "123",
		//	Type:     1,
		//}); err != nil {
		//	panic(err)
		//}
		//log.Info("send message success")

		for {
			select {
			case msg := <-ch:
				log.Info("client 2 recv message data = ", msg.Data)
			}
		}
	}()

	time.Sleep(time.Second)

	client, err := pClient.NewIMGoClientWithInstanceAddr("localhost:8080", "789")
	go func() {
		ch := client.GetRecvChan()
		for {
			select {
			case e := <-ch:
				if len(e.Data) == 0 {
					return
				}
				fmt.Printf("%+v\n", e)
			}
		}
	}()

	if err != nil {
		panic(err)
	}
	//if err := client.SendMsg(&pMessage.UserMessage{
	//	Data:     "hello laurence1",
	//	TargetID: "123",
	//	Type:     1,
	//}); err != nil {
	//	panic(err)
	//}
	////log.Info("send message success")
	//
	//if err := client.SendMsg(&pMessage.UserMessage{
	//	Data:     "hello laurence2",
	//	TargetID: "456",
	//	Type:     1,
	//}); err != nil {
	//	panic(err)
	//}
	if err := client.CreateGroup("group1", []string{"123", "456"}); err != nil {
		panic(err)
	}

	if err := client.SendMsg(&pMessage.UserMessage{
		Data:        "hello laurence1 and 2",
		TargetID:    "group1",
		PublishType: constant.UserMessagePublishType_Room,
	}); err != nil {
		panic(err)
	}

	//log.Info("send message success")
	time.Sleep(time.Second)
	fmt.Println("---------")

	if err := client.RemoveUserFromGroup("group1", "456"); err != nil {
		panic(err)
	}

	if err := client.SendMsg(&pMessage.UserMessage{
		Data:        "hello laurence1",
		TargetID:    "group1",
		PublishType: constant.UserMessagePublishType_Room,
	}); err != nil {
		panic(err)
	}

	time.Sleep(time.Second)
	fmt.Println("---------")

	if err := client.AddUserTOGroup("group1", "456"); err != nil {
		panic(err)
	}

	if err := client.SendMsg(&pMessage.UserMessage{
		Data:        "hello laurence1 & 2",
		TargetID:    "group1",
		PublishType: constant.UserMessagePublishType_Room,
	}); err != nil {
		panic(err)
	}

	if err := client.DeleteGroup("group1"); err != nil {
		panic(err)
	}

	time.Sleep(time.Second)
	client.Close()
	select {}

}
