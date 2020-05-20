package eltype

// ConfigType 配置类型
type ConfigType int

const (
	// ConfigTypeGlobal 全局生效的配置类型
	ConfigTypeGlobal ConfigType = iota
	// ConfigTypeFriend 仅对好友消息生效的配置类型
	ConfigTypeFriend
	// ConfigTypeGroup 仅对群消息生效的配置类型
	ConfigTypeGroup
)

// Config 配置对象
// @property	Type					ConfigType	配置类型
// @property	WhenMessageList			[]Message	作为触发条件的「消息」
// @property	WhenOperationList		[]Operation	作为触发条件的「事件/操作」
// @property	DoMessageList			[]Message	作为动作的「消息」
// @property	DoOperationList			[]Operation 作为动作的「操作」
type Config struct {
	Type              ConfigType
	WhenMessageList   []Message
	WhenOperationList []Operation
	DoMessageList     []Message
	DoOperationList   []Operation
	SenderList        []Sender
	Receiver          []Sender
}

// NewConfig 构造一个新的配置对象
// @param	configType				ConfigType	配置类型
// @param	WhenMessageList			[]Message	作为触发条件的「消息」
// @param	WhenOperationList		[]Operation	作为触发条件的「事件/操作」
// @param	DoMessageList			[]Message	作为动作的「消息」
// @param	DoOperationList			[]Operation 作为动作的「操作」
func NewConfig(configType ConfigType,
	whenMessageList []Message,
	whenOperationList []Operation,
	doMessageList []Message,
	doOperationList []Operation,
	SenderList []Sender,
	Receiver []Sender) (Config, error) {
	var config Config
	config.Type = configType
	config.WhenMessageList = whenMessageList
	config.WhenOperationList = whenOperationList
	config.DoMessageList = doMessageList
	config.DoOperationList = doOperationList
	config.SenderList = SenderList
	config.Receiver = Receiver
	return config, nil
}
