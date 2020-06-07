package eltype

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/mmcdole/gofeed"
	"github.com/robfig/cron"
	"gopkg.in/yaml.v2"
)

type RssListener struct {
	rssConfigList    []Config
	cron             *cron.Cron
	rssDataMap       map[string]string
	monthsMap        map[string]string
	WillBeSentConfig chan Config
	WillBeUsedEvent  chan Event
	Signal           chan SingalType
}

func NewRssListener(rssConfigList []Config) (*RssListener, error) {
	listener := new(RssListener)
	monthsMap := map[string]string{"January": "01", "February": "02", "March": "03",
		"April": "04", "May": "05", "June": "06", "July": "07", "August": "08", "September": "09",
		"October": "10", "November": "11", "December": "12"}
	listener.rssDataMap = make(map[string]string)
	listener.monthsMap = monthsMap
	listener.rssConfigList = make([]Config, len(rssConfigList))
	for _, config := range rssConfigList {
		listener.rssConfigList = append(listener.rssConfigList, config)
	}
	listener.WillBeSentConfig = make(chan Config, 10)
	listener.WillBeUsedEvent = make(chan Event, 10)
	listener.Signal = make(chan SingalType, 2)
	buf, err := ioutil.ReadFile(fmt.Sprintf("%s/%s", DataRoot, RssDataFileName))
	if err != nil {
		os.Create(fmt.Sprintf("%s/%s", DataRoot, RssDataFileName))
		return listener, err
	}
	yaml.Unmarshal(buf, &listener.rssDataMap)
	return listener, nil
}

func (listener *RssListener) Start() {
	go listener.start()
}

func (listener *RssListener) Destory() {
	listener.cron.Stop()
	listener.Signal <- Destory
	listener.Signal <- Destory
}

func (listener *RssListener) start() {
	listener.cron = cron.New()
	err := listener.cron.AddFunc("0 0/1 * * * *", func() {
		for _, rssConfig := range listener.rssConfigList {
			temp := listener.checkUpdate(rssConfig.RssURL)
			if temp != nil {
				event := Event{
					PreDefVarMap: temp,
				}
				listener.WillBeSentConfig <- rssConfig
				listener.WillBeUsedEvent <- event
			}
		}
	})
	if err != nil {
		log.Printf("RssListener.start: %s", err.Error())
		return
	}
	listener.cron.Start()
	select {
	case signalType := <-listener.Signal:
		if signalType == Destory {
			return
		}
	}
}

func (listener *RssListener) checkUpdate(url string) map[string]string {
	ret := make(map[string]string)
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	fp := gofeed.NewParser()
	feed, _ := fp.ParseURLWithContext(url, ctx)
	if len(feed.Items) == 0 || listener.rssDataMap[url] == feed.Items[0].Title ||
		listener.rssDataMap[url] == "" {
		return nil
	}
	year := fmt.Sprintf("%02d", feed.Items[0].UpdatedParsed.Year())
	month := listener.monthsMap[feed.Items[0].UpdatedParsed.Month().String()]
	day := fmt.Sprintf("%02d", feed.Items[0].UpdatedParsed.Day())
	hour := fmt.Sprintf("%02d", feed.Items[0].UpdatedParsed.Hour())
	minute := fmt.Sprintf("%02d", feed.Items[0].UpdatedParsed.Minute())
	second := fmt.Sprintf("%02d", feed.Items[0].UpdatedParsed.Second())
	ret["el-rss-author"] = feed.Author.Name
	ret["el-rss-title"] = feed.Items[0].Title
	ret["el-rss-year"] = year
	ret["el-rss-month"] = month
	ret["el-rss-day"] = day
	ret["el-rss-hour"] = hour
	ret["el-rss-minute"] = minute
	ret["el-rss-second"] = second
	ret["el-rss-link"] = feed.Items[0].Link
	ret["\\n"] = "\n"
	listener.rssDataMap[url] = feed.Items[0].Title

	ymlStr, err := yaml.Marshal(listener.rssDataMap)
	if err != nil {
		log.Printf("RssListener.checkUpdate: %s", err.Error())
	}
	err = ioutil.WriteFile(fmt.Sprintf("%s/%s", DataRoot, RssDataFileName), ymlStr, 0777)
	if err != nil {
		log.Printf("RssListener.checkUpdate: %s", err.Error())
	}
	return ret
}
