package eltype

// SenderType Sender 的类型
type SenderType int

const (
	// SenderTypeFriend Sender 是好友
	SenderTypeFriend SenderType = iota
	// SenderTypeGroup Sender 是群
	SenderTypeGroup
	// SenderTypeMember Sender 是群成员
	SenderTypeMember
	// SenderTypeUser Sender 是QQ用户
	SenderTypeUser
)

// Sender ...
// @property	Type		SenderType		Sender 的类型
// @property	ID			int64			Sender 的QQ号
// @property	Name		string			Sender 的名称
// @property	Permission	string			Sender 其它信息
type Sender struct {
	Type       SenderType
	ID         int64
	Name       string
	Permission string
}

// NewSender 构造一个 Sender
// @param	senderType		SenderType		Sender 的类型
// @param	ID				int64			Sender 的QQ号
// @param	Name			string			Sender 的名称
// @param	Permission		string			Sender 其它信息
func NewSender(senderType SenderType, ID int64, name string, permission string) (Sender, error) {
	var sender Sender
	sender.Type = senderType
	sender.ID = ID
	sender.Name = name
	sender.Permission = permission
	return sender, nil
}
