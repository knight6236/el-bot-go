package eltype

type PluginType int

const (
	// Binary 二进制文件
	Binary PluginType = iota
	// Java jar 包
	Java
	// JavaScript JS源代码
	JavaScript
	// Python .py
	Python
)

type Plugin struct {
	Type          PluginType
	Path          string
	ConfigKeyword string
}
