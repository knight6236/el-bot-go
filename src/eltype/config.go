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
	//ConfigTypeCrontab 定时执行的配置
	ConfigTypeCrontab
)

type Config struct {
	innerID int64
	Type    ConfigType
	IsCount bool
	RssURL  string `yaml:"url"`
	CountID string `yaml:"countID"`
	Cron    string `yaml:"cron"`
	When    When   `yaml:"when"`
	Do      Do     `yaml:"do"`
}

func (config *Config) DeepCopy() Config {
	newConfig := Config{
		innerID: config.innerID,
		Type:    config.Type,
		IsCount: config.IsCount,
		Cron:    config.Cron,
		RssURL:  config.RssURL,
		When:    config.When.DeepCopy(),
		Do:      config.Do.DeepCopy(),
	}
	return newConfig
}

func (config *Config) CompleteType() {
	config.When.Message.CompleteType()
	config.Do.Message.CompleteType()

	for i := 0; i < len(config.When.OperationList); i++ {
		temp := config.When.OperationList[i]
		temp.CompleteType()
		config.When.OperationList[i] = temp
	}
	for i := 0; i < len(config.Do.OperationList); i++ {
		temp := config.Do.OperationList[i]
		temp.CompleteType()
		config.Do.OperationList[i] = temp
	}
}

func (config *Config) CompleteContent(event Event) {
	config.When.Message.CompleteContent(event)
	config.Do.Message.CompleteContent(event)

	for i := 0; i < len(config.When.OperationList); i++ {
		temp := config.When.OperationList[i]
		temp.CompleteContent(event.PreDefVarMap)
		config.When.OperationList[i] = temp
	}
	for i := 0; i < len(config.Do.OperationList); i++ {
		temp := config.Do.OperationList[i]
		temp.CompleteContent(event.PreDefVarMap)
		config.Do.OperationList[i] = temp
	}
}
