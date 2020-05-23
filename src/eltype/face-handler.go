package eltype

// FaceHandler 判断是否命中和表情有关的配置
// @property	configList		[]Config			要判断的配置列表
// @property	messageList		[]Message			要判断的消息列表
// @property	operationList	[]Operation			要判断的配置列表
// @property	configHitList	[]Config			命中的配置列表
// @property	preDefVarMap	map[string]string	预定义变量 Map
type FaceHandler struct {
	configList    []Config
	messageList   []Message
	operationList []Operation
	configHitList []Config
	preDefVarMap  *map[string]string
}

// NewFaceHandler 构造一个 FaceHandler
// @param	configList		[]Config			要判断的配置列表
// @param	messageList		[]Message			要判断的消息列表
// @param	operationList	[]Operation			要判断的配置列表
// @param	preDefVarMap	map[string]string	预定义变量 Map
func NewFaceHandler(configList []Config, messageList []Message, operationList []Operation,
	preDefVarMap *map[string]string) (IHandler, error) {
	var handler FaceHandler
	handler.configList = configList
	handler.messageList = messageList
	handler.preDefVarMap = preDefVarMap
	handler.searchHitConfig()
	return handler, nil
}

func (handler *FaceHandler) searchHitConfig() {
	for _, config := range handler.configList {
		goto SECOND_LOOP
	TOP_LOOP:
		continue
	SECOND_LOOP:
		for _, message := range handler.messageList {
			for _, whenMessage := range config.WhenMessageList {
				if message.Type == MessageTypeFace &&
					whenMessage.Type == MessageTypeFace &&
					whenMessage.Value["name"] == message.Value["name"] {
					handler.configHitList = append(handler.configHitList, config)
					goto TOP_LOOP
				}
			}
		}
	}
}

// GetConfigHitList 获取命中的配置列表
func (handler FaceHandler) GetConfigHitList() []Config {
	return handler.configHitList
}
