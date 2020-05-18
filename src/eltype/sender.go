package eltype

type SenderType int

const (
	SenderTypeFriend SenderType = iota
	SenderTypeGroup
	SenderTypeMember
)

type Sender struct {
	Type       SenderType
	ID         int64
	Name       string
	Permission string
}

func NewSender(senderType SenderType, ID int64, name string, permission string) (Sender, error) {
	var sender Sender
	sender.Type = senderType
	sender.ID = ID
	sender.Name = name
	sender.Permission = permission
	return sender, nil
}
