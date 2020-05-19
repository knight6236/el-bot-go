package eltype

// OperationHandler 判断是否命中和表情有关的配置
// @property	configList		[]Config			要判断的配置列表
// @property	messageList		[]Message			要判断的消息列表
// @property	operationList	[]Operation			要判断的配置列表
// @property	configHitList	[]Config			命中的配置列表
// @property	preDefVarMap	map[string]string	预定义变量 Map
type OperationHandler struct {
	configList    []Config
	messageList   []Message
	operationList []Operation
	configHitList []Config
	preDefVarMap  map[string]string
}

// NewOperationHandler 构造一个 OperationHandler
// @param	configList		[]Config			要判断的配置列表
// @param	messageList		[]Message			要判断的消息列表
// @param	operationList	[]Operation			要判断的配置列表
// @param	preDefVarMap	map[string]string	预定义变量 Map
func NewOperationHandler(configList []Config, messageList []Message, operationList []Operation,
	preDefVarMap map[string]string) (IHandler, error) {
	var handler OperationHandler
	handler.configList = configList
	handler.operationList = operationList
	handler.preDefVarMap = preDefVarMap
	handler.searchHitConfig()
	return handler, nil
}

func (handler *OperationHandler) searchHitConfig() {
	for _, config := range handler.configList {
		goto SECOND_LOOP
	TOP_LOOP:
		continue
	SECOND_LOOP:
		for _, operation := range handler.operationList {
			for _, doOperation := range config.WhenOperationList {
				if operation.Type == doOperation.Type {
					handler.configHitList = append(handler.configHitList, config)
					goto TOP_LOOP
				}
			}
		}
	}
}

// GetConfigHitList 获取命中的配置列表
func (handler OperationHandler) GetConfigHitList() []Config {
	return handler.configHitList
}
