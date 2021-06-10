package message

import "github.com/online-im/go-sdk/pkg/constant"

type UserMessage struct {
	Data        string                          `json:"data"`
	TargetID    string                          `json:"target_id"`
	PublishType constant.UserMessagePublishType `json:"publish_type"`
	Type        constant.UserMessageType        `json:"type"`
	ErrMessage  string                          `json:"err_message"`
	ErrCode  uint32                          `json:"err_code"`
}
