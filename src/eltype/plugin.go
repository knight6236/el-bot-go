package eltype

// PluginType 插件类型
type PluginType int

const (
	// PluginTypeBinary 二进制文件
	PluginTypeBinary PluginType = iota
	// PluginTypeJava jar 包
	PluginTypeJava
	// PluginTypeJavaScript JS源代码
	PluginTypeJavaScript
	// PluginTypePython .py
	PluginTypePython
)

// Plugin 插件
type Plugin struct {
	// Type 插件类型
	Type PluginType
	// IsProcMsg 是否为消息处理插件
	IsProcMsg bool
	// RandKey 会话密钥
	RandKey string
	// Path 插件路径
	Path string
	// ConfigKeyword 对应配置的顶级关键字
	ConfigKeyword string
}
