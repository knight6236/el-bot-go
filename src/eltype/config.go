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
	innerID int
	Type    ConfigType
	IsCount bool
	CountID string `yaml:"countID"`
	Cron    string `yaml:"cron"`
	When    When   `yaml:"when"`
	Do      Do     `yaml:"do"`
}

func (config *Config) Init() {
	config.When.Message.Complete()
	config.Do.Message.Complete()
}
