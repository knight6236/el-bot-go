package eltype

type ConfigType int

const (
	ConfigTypeGlobal ConfigType = iota
	ConfigTypeFriend
	ConfigTypeGroup
)

type Config struct {
	Type              ConfigType
	WhenMessageList   []Message
	WhenOperationList []Operation
	DoMessageList     []Message
	DoOperationList   []Operation
}

func NewConfig(configType ConfigType,
	whenMessageList []Message,
	whenOperationList []Operation,
	doMessageList []Message,
	doOperationList []Operation) (Config, error) {
	var config Config
	config.Type = configType
	config.WhenMessageList = whenMessageList
	config.WhenOperationList = whenOperationList
	config.DoMessageList = doMessageList
	config.DoOperationList = doOperationList
	return config, nil
}
