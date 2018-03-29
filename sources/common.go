package sources

import (
	"gopkg.in/tucnak/telebot.v2"
	"time"
	"log"
	"telegram_webcomic_bot/configs"
)

var feedScrapers = []func(*telebot.Bot){
	scrapeDilbert,
	scrapeXkcd,
}

func updateFeeds(bot *telebot.Bot) {
	for _, f := range(feedScrapers) {
		f(bot)
	}
	conf := configs.GetConfigs()
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