package message

import (
	"github.com/glory-go/glory/log"
	"github.com/online-im/go-sdk/pkg/constant"
	pmessage "github.com/online-im/go-sdk/pkg/message"
	"github.com/online-im/online-im/pkg/client"
	cconstant "github.com/online-im/online-im/pkg/constant"
)

// MsgController aims at dealing the transfer between user and online-im core
type MsgController struct {
	uch       chan pmessage.UserMessage
	cclient   *client.IMCoreClient
	closeChan chan struct{}
}

func NewMsgController(ch chan pmessage.UserMessage, cclient *client.IMCoreClient) *MsgController {
	newCtrler := &MsgController{
		uch:       ch,
		cclient:   cclient,
		closeChan: make(chan struct{}),
	}
	go func() {
		defer func() {
			close(newCtrler.uch)
			if e := recover(); e != nil {
				log.Errorf("IMGoClient recv message err = %s", e)
			}
		}()
		for {
			select {
			case <-newCtrler.closeChan:
				return
			default:
				recv, err := newCtrler.cclient.RecvMessage()
				if err != nil {
					return
				}
				switch recv.Type {
				case cconstant.CoreMessageType_Message:
					// deal with target message
					userMsg := pmessage.UserMessage{
						Data:        string(recv.MessagePayload.Data),
						TargetID:    recv.MessagePayload.TargetID,
						PublishType: constant.UserMessagePublishType_User,
						Type:        constant.UserMessageType_Receive,
					}
					ch <- userMsg
				case cconstant.CoreMessageType_ErrorMessage:
					var userMsg pmessage.UserMessage
					switch recv.Err.Code {
					case cconstant.CoreErrorCode_UserOffLine:
						userMsg = pmessage.UserMessage{
							ErrMessage: "user offline",
							Type:       constant.UserMessageType_Notify,
						}
					default:
						userMsg = pmessage.UserMessage{
							ErrMessage: recv.Err.Message,
							Type:       constant.UserMessageType_Error,
						}
					}
					ch <- userMsg
				}
			}

		}
	}()
	return newCtrler
}

func (mc *MsgController) SendMsg(msg *pmessage.UserMessage) error {
	// todo add  protocol
	err := mc.cclient.SendMessage([]byte(msg.Data), msg.TargetID, cconstant.PublishType(msg.PublishType))
	if err != nil {
		log.Errorf("IMGoClient Send UserMessage %s failed with error = %s", msg, err)
		return err
	}
	return nil
}

func (mc *MsgController) Close() {
	mc.cclient.Close()
	close(mc.closeChan)
}

func (mc *MsgController) CreateGroup(userID, groupID string, users []string) error {
	return mc.cclient.CreateGroup(userID, groupID, users)
}

func (mc *MsgController) DeleteGroup(groupID string) error {
	return mc.cclient.DeleteGroup(groupID)
}

func (mc *MsgController) AddUserTOGroup(groupID, userID string) error {
	return mc.cclient.AddUserToGroup(userID, groupID)
}

func (mc *MsgController) RemoveUserFromGroup(groupID, userID string) error {
	return mc.cclient.RemoveUserFromGroup(userID, groupID)
}
