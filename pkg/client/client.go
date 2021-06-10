package client

import (
	"encoding/json"
	"github.com/glory-go/glory/common"
	"github.com/glory-go/glory/config"
	"github.com/glory-go/glory/log"
	"github.com/glory-go/glory/plugin"
	_ "github.com/glory-go/glory/registry/k8s"
	"github.com/online-im/go-sdk/internal/constant"
	"github.com/online-im/go-sdk/internal/message"
	internalPkg "github.com/online-im/go-sdk/internal/pkg"
	pmessage "github.com/online-im/go-sdk/pkg/message"
	"github.com/online-im/online-im/pkg/client"
	perrors "github.com/pkg/errors"
	"io/ioutil"
	"net/http"
)

type IMGoClient struct {
	userRecvChan  chan pmessage.UserMessage
	msgController *message.MsgController
	userID        string
}

func NewIMGoClientWithInstanceAddr(instanceAddr string, userID string) (*IMGoClient, error) {
	coreClient, err := client.NewImCoreClient(instanceAddr, userID)
	if err != nil {
		log.Errorf("NewIMGoClientWithInstanceAddr = %s err = %v", instanceAddr, err)
		return nil, err
	}
	ch := make(chan pmessage.UserMessage, 0)
	return &IMGoClient{
		userRecvChan:  ch,
		msgController: message.NewMsgController(ch, coreClient),
		userID:        userID,
	}, nil
}

func NewIMGoClientWithK8sServiceName(serviceName, userID string) (*IMGoClient, error) {
	return NewIMGoClientWithInstanceAddr(serviceName+":8080", userID)
}

func NewIMGoClientWithGatewayAddr(gatewayAddr string, userID string) (*IMGoClient, error) {
	// todo add timeout
	instanceAddr, err := getInstanceAddr(gatewayAddr)
	if err != nil {
		return nil, err
	}
	return NewIMGoClientWithInstanceAddr(instanceAddr, userID)
}

// NewIMGoClientWithoutGatewayAddrInCluster deprecated
func NewIMGoClientWithoutGatewayAddrInCluster(userID string) (*IMGoClient, error) {
	// todo add timeout
	gatewayAddr, err := discFistGateway()
	if err != nil {
		return nil, err
	}
	return NewIMGoClientWithGatewayAddr(gatewayAddr, userID)
}

func discFistGateway() (string, error) {
	reg := plugin.GetRegistry(&config.RegistryConfig{Service: constant.GloryK8SPluginName})
	eventCh, err := reg.Subscribe(constant.OnlineIMGatewayID)
	if err != nil {
		panic(err)
	}

	for v := range eventCh {
		addr := v.Addr.GetUrl()
		switch v.Opt {
		case common.RegistryAddEvent, common.RegistryUpdateEvent:
			return addr, nil
		}
	}
	return "", perrors.Errorf("disc gateway by " + constant.OnlineIMGatewayID + " using k8s failed")
}

func getInstanceAddr(gatewayAddr string) (string, error) {
	rsp, err := http.Get("http://" + gatewayAddr + "/connect")
	if err != nil {
		log.Errorf("get instance addr from gateway %s err = %s", gatewayAddr, err)
		return "", err
	}
	body, err := ioutil.ReadAll(rsp.Body)
	getInstanceRsp := &internalPkg.ConnRsp{}
	if err := json.Unmarshal(body, getInstanceRsp); err != nil {
		return "", err
	}
	if getInstanceRsp.Ok != true {
		return "", perrors.Errorf("get instance addr from gateway %s not ok", gatewayAddr)
	}
	return getInstanceRsp.Address, nil
}

func (c *IMGoClient) SendMsg(message *pmessage.UserMessage) error {
	return c.msgController.SendMsg(message)
}

func (c *IMGoClient) GetRecvChan() chan pmessage.UserMessage {
	return c.userRecvChan
}

func (c *IMGoClient) Close() {
	c.msgController.Close()
}

func (c *IMGoClient) CreateGroup(groupID string, users []string) error {
	return c.msgController.CreateGroup(c.userID, groupID, users)
}

func (c *IMGoClient) DeleteGroup(groupID string) error {
	return c.msgController.DeleteGroup(groupID)
}
func (mc *IMGoClient) AddUserTOGroup(groupID, userID string) error {
	return mc.msgController.AddUserTOGroup(groupID, userID)
}

func (mc *IMGoClient) RemoveUserFromGroup(groupID, userID string) error {
	return mc.msgController.RemoveUserFromGroup(groupID, userID)
}
