package eltype

// IHandler 所有 Handler 均实现此接口
// @method	GetConfigHitList	[]Config	获取命中的配置列表
type IHandler interface {
	GetConfigHitList() []Config
}
