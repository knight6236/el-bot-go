package eltype

// XMLHandler 判断是否命中和表情有关的配置
// @property	configList		[]Config			要判断的配置列表
// @property	messageList		[]Message			要判断的消息列表
// @property	operationList	[]Operation			要判断的配置列表
// @property	configHitList	[]Config			命中的配置列表
// @property	preDefVarMap	map[string]string	预定义变量 Map
type XMLHandler struct {
	configList    []Config
	messageList   []Message
	configHitList []Config
	operationList []Operation
	preDefVarMap  *map[string]string
}

// NewXMLHandler 构造一个 XMLHandler
// @param	configList		[]Config			要判断的配置列表
// @param	messageList		[]Message			要判断的消息列表
// @param	operationList	[]Operation			要判断的配置列表
// @param	preDefVarMap	map[string]string	预定义变量 Map
func NewXMLHandler(configList []Config, messageList []Message, operationList []Operation,
	preDefVarMap *map[string]string) (IHandler, error) {
	var handler XMLHandler
	handler.configList = configList
	handler.messageList = messageList
	handler.preDefVarMap = preDefVarMap
	handler.searchHitConfig()
	return handler, nil
}

func (handler *XMLHandler) searchHitConfig() {
	for _, config := range handler.configList {
		for _, whenMessage := range config.WhenMessageList {
			for _, message := range handler.messageList {
				if whenMessage.Type == MessageTypeXML &&
					message.Type == MessageTypeXML {
					handler.configHitList = append(handler.configHitList, config)
				}
			}
		}
	}
}

// GetConfigHitList 获取命中的配置列表
func (handler XMLHandler) GetConfigHitList() []Config {
	return handler.configHitList
}
