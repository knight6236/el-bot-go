package eltype

import (
	"log"
	"sync"

	"github.com/robfig/cron"
)

type FreqMonitor struct {
	mute               sync.RWMutex
	cron               *cron.Cron
	freqUpperLimit     int64
	groupCountMap      map[int64]map[int64]int64
	userCountMap       map[int64]map[int64]int64
	CountMap           map[string]int64
	groupBlockedConfig map[int64]map[int64]bool
	userBlockedConfig  map[int64]map[int64]bool
	Singal             chan SingalType
}

func NewFreqMonitor(freqUpperLimit int64) (*FreqMonitor, error) {
	monitor := new(FreqMonitor)
	monitor.freqUpperLimit = freqUpperLimit
	monitor.groupCountMap = make(map[int64]map[int64]int64)
	monitor.userCountMap = make(map[int64]map[int64]int64)
	monitor.groupBlockedConfig = make(map[int64]map[int64]bool)
	monitor.userBlockedConfig = make(map[int64]map[int64]bool)
	monitor.CountMap = make(map[string]int64)
	monitor.Singal = make(chan SingalType, 2)
	return monitor, nil
}

func (monitor *FreqMonitor) Start() {
	go monitor.autoClear()
}

func (monitor *FreqMonitor) Stop() {
	monitor.cron.Stop()
	monitor.Singal <- SingalTypeStop
	monitor.Singal <- SingalTypeStop
}

func (monitor *FreqMonitor) Commit(configHit Config) {
	monitor.mute.Lock()
	for _, groupID := range configHit.When.Message.Receiver.GroupIDList {
		if monitor.groupCountMap[CastStringToInt64(groupID)] == nil {
			monitor.groupCountMap[CastStringToInt64(groupID)] = make(map[int64]int64)
		}
		monitor.groupCountMap[CastStringToInt64(groupID)][configHit.innerID]++
		// fmt.Printf("%v\n", monitor.groupCountMap)
	}
	for _, userID := range configHit.When.Message.Receiver.UserIDList {
		if monitor.groupCountMap[CastStringToInt64(userID)] == nil {
			monitor.userCountMap[CastStringToInt64(userID)] = make(map[int64]int64)
		}
		monitor.userCountMap[CastStringToInt64(userID)][configHit.innerID]++
		// fmt.Printf("%v\n", monitor.userCountMap)
	}
	monitor.CountMap[configHit.CountID]++
	// fmt.Printf("\n\n")
	monitor.mute.Unlock()
	monitor.check()
}

func (monitor *FreqMonitor) IsBlocked(configInnerID int64, receiverType ReceiverType, receiverID int64) bool {
	var isBlocked bool
	monitor.mute.RLock()
	switch receiverType {
	case ReceiverTypeGroup:
		isBlocked = monitor.groupBlockedConfig[receiverID][configInnerID]
	case ReceiverTypeUser:
		isBlocked = monitor.userBlockedConfig[receiverID][configInnerID]
	default:
		isBlocked = false
	}
	monitor.mute.RUnlock()
	return isBlocked
}

func (monitor *FreqMonitor) check() {
	if monitor.freqUpperLimit == 0 {
		return
	}
	monitor.mute.RLock()
	for groupID, innerMap := range monitor.groupCountMap {
		for innerID, freq := range innerMap {
			if freq > monitor.freqUpperLimit {
				if monitor.groupBlockedConfig[groupID] == nil {
					monitor.groupBlockedConfig[groupID] = make(map[int64]bool)
				}
				monitor.groupBlockedConfig[groupID][innerID] = true
			}
		}
	}
	for userID, innerMap := range monitor.userCountMap {
		for innerID, freq := range innerMap {
			if freq > monitor.freqUpperLimit {
				if monitor.userBlockedConfig[userID] == nil {
					monitor.userBlockedConfig[userID] = map[int64]bool{}
				}
				monitor.userBlockedConfig[userID][innerID] = true
			}
		}
	}
	monitor.mute.RUnlock()
}

func (monitor *FreqMonitor) autoClear() {
	monitor.cron = cron.New()
	err := monitor.cron.AddFunc("0 * * * * *", func() {
		monitor.mute.Lock()
		monitor.groupCountMap = make(map[int64]map[int64]int64)
		monitor.userCountMap = make(map[int64]map[int64]int64)
		monitor.groupBlockedConfig = make(map[int64]map[int64]bool)
		monitor.userBlockedConfig = make(map[int64]map[int64]bool)
		monitor.mute.Unlock()
	})
	if err != nil {
		log.Printf("Monitor.FreqMonitor: %s", err.Error())
	}
	monitor.cron.Start()
	select {
	case signalType := <-monitor.Singal:
		if signalType == SingalTypeStop {
			return
		}
	}
}
