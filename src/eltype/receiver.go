package eltype

import (
	"fmt"
	"strings"
)

// ReceiverType Sender 的类型
type ReceiverType int

const (
	// ReceiverTypeGroup Sender 是群
	ReceiverTypeGroup ReceiverType = iota
	// ReceiverTypeUser Sender 是QQ用户
	ReceiverTypeUser
)

// Receiver ...
// @property	Type		SenderType		Sender 的类型
// @property	ID			int64			Sender 的QQ号
// @property	Name		string			Sender 的名称
// @property	Permission	string			Sender 其它信息
type Receiver struct {
	GroupIDList []string `yaml:"group"`
	UserIDList  []string `yaml:"user"`
}

// NewReceiver 构造一个 Receiver
// @param	senderType		SenderType		Sender 的类型
// @param	ID				int64			Sender 的QQ号
// @param	Name			string			Sender 的名称
// @param	Permission		string			Sender 其它信息
// func NewReceiver(receiveType ReceiverType, ID int64, name string, permission string) (Receiver, error) {
// 	var receiver Receiver
// 	receiver.Type = receiveType
// 	receiver.ID = ID
// 	receiver.Name = name
// 	receiver.Permission = permission
// 	return receiver, nil
// }

func (receiver *Receiver) DeepCopy() Receiver {
	var newReceiver Receiver
	for _, item := range receiver.GroupIDList {
		newReceiver.GroupIDList = append(newReceiver.GroupIDList, item)
	}
	for _, item := range receiver.UserIDList {
		newReceiver.UserIDList = append(newReceiver.UserIDList, item)
	}
	return newReceiver
}

func (receiver *Receiver) Complete(preDefVarMap map[string]string) {
	for key, value := range preDefVarMap {
		varName := fmt.Sprintf("{%s}", key)
		for i := 0; i < len(receiver.GroupIDList); i++ {
			receiver.GroupIDList[i] = strings.ReplaceAll(receiver.GroupIDList[i], varName, value)
		}
		for i := 0; i < len(receiver.UserIDList); i++ {
			receiver.UserIDList[i] = strings.ReplaceAll(receiver.UserIDList[i], varName, value)
		}
	}
}

func (receiver *Receiver) AddGroupID(groupID int64) {
	receiver.GroupIDList = append(receiver.GroupIDList, CastInt64ToString(groupID))
}

func (receiver *Receiver) AddUserID(userID int64) {
	receiver.UserIDList = append(receiver.UserIDList, CastInt64ToString(userID))
}
