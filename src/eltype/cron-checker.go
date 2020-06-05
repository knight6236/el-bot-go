package eltype

import (
	"fmt"

	"github.com/robfig/cron"
)

type CronChecker struct {
	cronConfigList   []Config
	cron             *cron.Cron
	WillBeSentConfig chan Config
	Signal           chan SingalType
}

type CronJob struct {
	config           Config
	willBeSentConfig chan Config
}

func NewCronChecker(cronConfigList []Config) (*CronChecker, error) {
	checker := new(CronChecker)
	for _, config := range cronConfigList {
		checker.cronConfigList = append(checker.cronConfigList, config)
	}
	checker.WillBeSentConfig = make(chan Config, 10)
	checker.Signal = make(chan SingalType, 2)
	return checker, nil
}

func (checker *CronChecker) Start() {
	go checker.start()
}

func (checker *CronChecker) Stop() {
	checker.cron.Stop()
	checker.Signal <- SingalTypeStop
	checker.Signal <- SingalTypeStop
}

func (checker *CronChecker) start() {
	checker.cron = cron.New()
	for _, config := range checker.cronConfigList {
		err := checker.cron.AddJob(config.Cron, CronJob{config: config, willBeSentConfig: checker.WillBeSentConfig})
		if err != nil {
			fmt.Println(err)
		}
	}
	checker.cron.Start()
	select {
	case sigalType := <-checker.Signal:
		if sigalType == SingalTypeStop {
			return
		}
	}
}

// Run TODO
func (job CronJob) Run() {
	job.willBeSentConfig <- job.config
}
