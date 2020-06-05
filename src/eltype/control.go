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
)

type Control struct {
	innerType   ControlType
	Type        string   `yaml:"type"`
	GroupIDList []string `yaml:"group"`
	UserIDList  []string `yaml:"user"`
}

func (control *Control) CompleteType() {
	if control.Type != "" {
		switch control.Type {
		case "Suspend":
			control.innerType = ControlTypeSuspend
		case "Active":
			control.innerType = ControlTypeActive
		case "Destory":
			control.innerType = ControlTypeDestory
		case "Restart":
			control.innerType = ControlTypeRestart
		case "Block":
			control.innerType = ControlTypeBlock
		case "Unblock":
			control.innerType = ControlTypeUnblock
		}
	}

	switch control.innerType {
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
	newControl.innerType = control.innerType
	newControl.Type = control.Type
	for _, groupID := range control.GroupIDList {
		newControl.GroupIDList = append(newControl.GroupIDList, groupID)
	}
	newControl.UserIDList = make([]string, len(control.UserIDList))
	for _, userID := range control.UserIDList {
		newControl.UserIDList = append(newControl.UserIDList, userID)
	}
	return newControl
}
