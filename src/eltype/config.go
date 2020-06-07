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
	InnerID int64      `json:"-"`
	Type    ConfigType `json:"-"`
	IsCount bool       `json:"-"`
	RssURL  string     `yaml:"url" json:"url"`
	CountID string     `yaml:"countID" json:"countID"`
	Cron    string     `yaml:"cron" json:"cron"`
	When    When       `yaml:"when" json:"when"`
	Do      Do         `yaml:"do" json:"do"`
}

func (config *Config) DeepCopy() Config {
	newConfig := Config{
		InnerID: config.InnerID,
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
		temp.CompleteContent(event)
		config.When.OperationList[i] = temp
	}
	for i := 0; i < len(config.Do.OperationList); i++ {
		temp := config.Do.OperationList[i]
		temp.CompleteContent(event)
		config.Do.OperationList[i] = temp
	}
}
