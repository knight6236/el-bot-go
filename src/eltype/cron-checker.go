package eltype

import (
	"fmt"

	"github.com/robfig/cron"
)

type CronChecker struct {
	cronConfigList   []Config
	WillBeSentConfig chan Config
}

type CronJob struct {
	config           Config
	willBeSentConfig chan Config
}

func NewCronChecker(cronConfigList []Config) (*CronChecker, error) {
	checker := new(CronChecker)
	checker.cronConfigList = cronConfigList
	checker.WillBeSentConfig = make(chan Config, 10)
	return checker, nil
}

func (checker *CronChecker) Start() {
	go checker.start()
}

func (checker *CronChecker) start() {
	c := cron.New()
	for _, config := range checker.cronConfigList {
		err := c.AddJob(config.Cron, CronJob{config: config, willBeSentConfig: checker.WillBeSentConfig})
		if err != nil {
			fmt.Println(err)
		}
	}
	c.Start()
	select {}
}

// Run TODO
func (job CronJob) Run() {
	job.willBeSentConfig <- job.config
}
