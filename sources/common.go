package sources

import (
	"gopkg.in/tucnak/telebot.v2"
	"time"
	"log"
	"telegram_webcomic_bot/configs"
	"github.com/mmcdole/gofeed"
)

type feedSrc struct {
	url string
	name string
	scraper func([]*gofeed.Item, *configs.Configs, feedSrc)([]comicUpdate, time.Time)
}

var feedScrapers = []feedSrc{
	{
		"http://dilbert.com/feed",
		"Dilbert",
		scrapeDilbert,
	},
	{
		"https://www.xkcd.com/rss.xml",
		"xkcd",
		scrapeXkcd,
	},
}

func updateFeeds(bot *telebot.Bot) {
	conf := configs.GetConfigs()
	fp := gofeed.NewParser()

	for _, f := range(feedScrapers) {
		feed, err := fp.ParseURL(f.url)
		if err != nil {
			parseError(f.name, err)
			continue
		}

		comics, lastTime := f.scraper(feed.Items, conf, f)

		if len(comics) > 0 {
			notifyComics(bot, f.name, comics)
			conf.UpdateFeed(f.name, lastTime)
		}
	}

	newSources := conf.GetNewFeedSources()
	if len(newSources) > 0 {
		notifyNewSources(bot, newSources)
	}
}

var feedUpdateDelay = "1h"

func KeepFeedsUpdated(bot *telebot.Bot) {
	sleepDuration, err := time.ParseDuration(feedUpdateDelay)
	if err != nil {
		sleepDuration = time.Hour
	}
	for true {
		log.Print("Updating feeds...")
		updateFeeds(bot)
		log.Print("Feeds updated.")

		time.Sleep(sleepDuration)
	}
}

func parseError(source string, err error) {
	log.Printf("Error parsing %s: %s", source, err.Error())
}