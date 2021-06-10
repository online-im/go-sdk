package constant

type UserMessageType uint32

const UserMessageType_Receive = UserMessageType(0)
const UserMessageType_Notify = UserMessageType(1)
const UserMessageType_Error = UserMessageType(2)

type UserMessagePublishType uint32

const UserMessagePublishType_User = UserMessagePublishType(1)
const UserMessagePublishType_Room = UserMessagePublishType(2)
const UserMessagePublishType_Global = UserMessagePublishType(3)
