package eltype

// Job TODO
type Job struct {
	controller *Controller
	config     Config
}

// Run TODO
func (job Job) Run() {
	var event Event
	var configHitList []Config
	configHitList = append(configHitList, job.config)
	job.controller.sendMessageAndOperation(event, configHitList)
}
