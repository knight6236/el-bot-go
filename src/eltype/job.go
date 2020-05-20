package eltype

type Job struct {
	controller *Controller
	config     Config
}

func (job Job) Run() {
	var event Event
	var configHitList []Config
	configHitList = append(configHitList, job.config)
	sendedGoMiraiMessageList := job.controller.getSendedGoMiraiMessageList(event, configHitList)
	job.controller.sendMessage(event, configHitList, sendedGoMiraiMessageList)
}
