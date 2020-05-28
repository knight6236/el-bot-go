package eltype

import (
	"fmt"
	"strings"
)

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
	GroupIDList []string `yaml:"group"`
	UserIDList  []string `yaml:"user"`
}

// NewSender 构造一个 Sender
// @param	senderType		SenderType		Sender 的类型
// @param	ID				int64			Sender 的QQ号
// @param	Name			string			Sender 的名称
// @param	Permission		string			Sender 其它信息
// func NewSender(senderType SenderType, ID int64, name string, permission string) (Sender, error) {
// 	var sender Sender
// 	sender.Type = senderType
// 	sender.ID = ID
// 	sender.Name = name
// 	sender.Permission = permission
// 	return sender, nil
// }

func (sender *Sender) DeepCopy() Sender {
	var newSender Sender
	for _, item := range sender.GroupIDList {
		newSender.GroupIDList = append(newSender.GroupIDList, item)
	}
	for _, item := range sender.UserIDList {
		newSender.UserIDList = append(newSender.UserIDList, item)
	}
	return newSender
}

func (sender *Sender) Complete(preDefVarMap map[string]string) {
	for key, value := range preDefVarMap {
		varName := fmt.Sprintf("{%s}", key)
		for i := 0; i < len(sender.GroupIDList); i++ {
			sender.GroupIDList[i] = strings.ReplaceAll(sender.GroupIDList[i], varName, value)
		}
		for i := 0; i < len(sender.UserIDList); i++ {
			sender.UserIDList[i] = strings.ReplaceAll(sender.UserIDList[i], varName, value)
		}
	}
}

func (sender *Sender) AddGroupID(groupID string) {
	sender.GroupIDList = append(sender.GroupIDList, groupID)
}

func (sender *Sender) AddUserID(userID string) {
	sender.UserIDList = append(sender.UserIDList, userID)
}
