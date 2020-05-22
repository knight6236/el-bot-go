package eltype

// ReceiverType Sender 的类型
type ReceiverType int

const (
	// ReceiverTypeFriend Sender 是好友
	ReceiverTypeFriend ReceiverType = iota
	// ReceiverTypeGroup Sender 是群
	ReceiverTypeGroup
	// ReceiverTypeMember Sender 是群成员
	ReceiverTypeMember
	// ReceiverTypeUser Sender 是QQ用户
	ReceiverTypeUser
)

// Receiver ...
// @property	Type		SenderType		Sender 的类型
// @property	ID			int64			Sender 的QQ号
// @property	Name		string			Sender 的名称
// @property	Permission	string			Sender 其它信息
type Receiver struct {
	Type       ReceiverType
	ID         int64
	Name       string
	Permission string
}

// NewReceiver 构造一个 Receiver
// @param	senderType		SenderType		Sender 的类型
// @param	ID				int64			Sender 的QQ号
// @param	Name			string			Sender 的名称
// @param	Permission		string			Sender 其它信息
func NewReceiver(receiveType ReceiverType, ID int64, name string, permission string) (Receiver, error) {
	var receiver Receiver
	receiver.Type = receiveType
	receiver.ID = ID
	receiver.Name = name
	receiver.Permission = permission
	return receiver, nil
}
