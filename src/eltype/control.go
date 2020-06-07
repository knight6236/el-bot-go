package eltype

import "strings"

type ControlType int

const (
	// ControlTypeSuspend 挂起机器人
	ControlTypeSuspend ControlType = iota
	// ControlTypeActive 激活机器人
	ControlTypeActive
	// ControlTypeDestory 终止机器人
	ControlTypeDestory
	// ControlTypeRestart 重启机器人
	ControlTypeRestart
	// ControlTypeBlock 屏蔽指定消息接收者
	ControlTypeBlock
	// ControlTypeUnblock 取消被屏蔽的一些消息接收者
	ControlTypeUnblock
	ControlTypeEnterConfig
	ControlTypeBackToPrevConfig
)

type Control struct {
	InnerType   ControlType `json:"-"`
	Folder      string      `yaml:"folder" json:"folder"`
	Type        string      `yaml:"type" json:"type"`
	GroupIDList []string    `yaml:"group" json:"group"`
	UserIDList  []string    `yaml:"user" json:"user"`
}

func (control *Control) CompleteType() {
	if control.Type != "" {
		switch control.Type {
		case "Suspend":
			control.InnerType = ControlTypeSuspend
		case "Active":
			control.InnerType = ControlTypeActive
		case "Destory":
			control.InnerType = ControlTypeDestory
		case "Restart":
			control.InnerType = ControlTypeRestart
		case "Block":
			control.InnerType = ControlTypeBlock
		case "Unblock":
			control.InnerType = ControlTypeUnblock
		case "EnterConfig":
			control.InnerType = ControlTypeEnterConfig
		case "BackToPrevConfig":
			control.InnerType = ControlTypeBackToPrevConfig
		}
	}

	switch control.InnerType {
	case ControlTypeSuspend:
		control.Type = "Suspend"
	case ControlTypeActive:
		control.Type = "Active"
	case ControlTypeDestory:
		control.Type = "Destory"
	case ControlTypeRestart:
		control.Type = "Restart"
	case ControlTypeBlock:
		control.Type = "Block"
	case ControlTypeUnblock:
		control.Type = "Unblock"
	case ControlTypeEnterConfig:
		control.Type = "EnterConfig"
	case ControlTypeBackToPrevConfig:
		control.Type = "BackToPrevConfig"
	}
}

func (control *Control) CompleteContent(event Event) {
	for varName, value := range event.PreDefVarMap {
		for i := 0; i < len(control.GroupIDList); i++ {
			control.GroupIDList[i] = strings.ReplaceAll(control.GroupIDList[i], varName, value)
		}
		for i := 0; i < len(control.UserIDList); i++ {
			control.UserIDList[i] = strings.ReplaceAll(control.UserIDList[i], varName, value)
		}
	}
}

func (control *Control) DeepCopy() Control {
	var newControl Control
	newControl.InnerType = control.InnerType
	newControl.Type = control.Type
	newControl.Folder = control.Folder
	for _, groupID := range control.GroupIDList {
		newControl.GroupIDList = append(newControl.GroupIDList, groupID)
	}
	newControl.UserIDList = make([]string, len(control.UserIDList))
	for _, userID := range control.UserIDList {
		newControl.UserIDList = append(newControl.UserIDList, userID)
	}
	return newControl
}
