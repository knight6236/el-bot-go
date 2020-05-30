package eltype

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/mmcdole/gofeed"
	"gopkg.in/yaml.v2"
)

type RssListener struct {
	rssConfigList []Config
	rssDataMap    map[string]string
	monthsMap     map[string]string
}

func NewRssListener(rssConfigList []Config) (RssListener, error) {
	var listener RssListener
	monthsMap := map[string]string{"January": "01", "February": "02", "March": "03",
		"April": "04", "May": "05", "June": "06", "July": "07", "August": "08", "September": "09",
		"October": "10", "November": "11", "December": "12"}
	listener.rssDataMap = make(map[string]string)
	listener.monthsMap = monthsMap
	listener.rssConfigList = rssConfigList
	buf, err := ioutil.ReadFile(fmt.Sprintf("%s/%s", DataRoot, RssDataFileName))
	if err != nil {
		os.Create(fmt.Sprintf("%s/%s", DataRoot, RssDataFileName))
		return listener, err
	}
	yaml.Unmarshal(buf, &listener.rssDataMap)
	return listener, nil
}

func (listen *RssListener) checkUpdate(url string) map[string]string {
	ret := make(map[string]string)
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	fp := gofeed.NewParser()
	feed, _ := fp.ParseURLWithContext(url, ctx)
	if len(feed.Items) == 0 || listen.rssDataMap[url] == feed.Items[0].Title ||
		listen.rssDataMap[url] == "" {
		return nil
	}
	year := fmt.Sprintf("%02d", feed.Items[0].UpdatedParsed.Year())
	month := listen.monthsMap[feed.Items[0].UpdatedParsed.Month().String()]
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
	listen.rssDataMap[url] = feed.Items[0].Title

	ymlStr, err := yaml.Marshal(listen.rssDataMap)
	if err != nil {
		log.Println(err)
	}
	err = ioutil.WriteFile(fmt.Sprintf("%s/%s", DataRoot, RssDataFileName), ymlStr, 0777)
	if err != nil {
		log.Println(err)
	}
	return ret
}
