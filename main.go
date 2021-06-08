package go_sdk

import (
	"encoding/json"
	"github.com/glory-go/glory/common"
	"github.com/glory-go/glory/config"
	"github.com/glory-go/glory/log"
	"github.com/glory-go/glory/plugin"
	_ "github.com/glory-go/glory/registry/k8s"
	"github.com/online-im/go-sdk/internal/constant"
	internalPkg "github.com/online-im/go-sdk/internal/pkg"
	"github.com/online-im/go-sdk/pkg"
	perrors "github.com/pkg/errors"
	"golang.org/x/net/websocket"
	"io/ioutil"
	"net/http"
)

type IMGoClient struct {
	ws       *websocket.Conn
	recvChan chan pkg.Message
}

func NewIMGoClientWithGatewayAddr(gatewayAddr string, userID string) (*IMGoClient, error) {
	// todo add timeout
	instanceAddr, err := getInstanceAddr(gatewayAddr)
	if err != nil {
		return nil, err
	}
	ws, err := websocket.Dial("ws://"+instanceAddr+"/online?fromid="+userID, "", "http://"+instanceAddr+"/")
	if err != nil {
		log.Errorf("websocket dial instance addr = %s failed!, with err = %s", instanceAddr, err)
		return nil, err
	}
	ch := make(chan pkg.Message, 0)
	go func() {
		defer func() {
			if e := recover(); e != nil {
				log.Errorf("IMGoClient recv message err = %s", e)
			}
		}()
		for {
			recv := pkg.Message{}
			if err := websocket.JSON.Receive(ws, &recv); err != nil {
				panic(err)
			}
			ch <- recv
		}
	}()
	return &IMGoClient{
		ws:       ws,
		recvChan: ch,
	}, nil
}

func NewIMGoClient(userID string) (*IMGoClient, error) {
	// todo add timeout
	gatewayAddr, err := discFistGateway()
	if err != nil {
		return nil, err
	}
	return NewIMGoClientWithGatewayAddr(gatewayAddr, userID)
}

func (c *IMGoClient) SendMsg(message *pkg.Message) error {
	err := websocket.JSON.Send(c.ws, *message)
	if err != nil {
		log.Errorf("IMGoClient Send Message %s failed with error = %s", message, err)
		return err
	}
	return nil
}

func (c *IMGoClient) GetRecvChan() chan pkg.Message {
	return c.recvChan
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
	getInstanceRsp := &internalPkg.WrappedGloryHttpRsp{}
	if err := json.Unmarshal(body, getInstanceRsp); err != nil {
		return "", err
	}
	if getInstanceRsp.Result.Ok != true {
		return "", perrors.Errorf("get instance addr from gateway %s not ok", gatewayAddr)
	}
	return getInstanceRsp.Result.Address, nil
}
