package sources

import (
	"gopkg.in/tucnak/telebot.v2"
	"time"
	"log"
	"telegram_webcomic_bot/configs"
	"github.com/mmcdole/gofeed"
	"telegram_webcomic_bot/sources/scrapers"
)

// new comics must be added to this array, make sure the "name" field is unique

var feedScrapers = []feedSrc{
	{
		"http://dilbert.com/feed",
		"Dilbert",
		timeParserUpdated(time.RFC3339),
		scrapers.ScrapeDilbert,
	},
	{
		"https://www.xkcd.com/rss.xml",
		"xkcd",
		timeParserPublished(time.RFC1123Z),
		scrapers.ScrapeXkcd,
	},
	{
		"https://what-if.xkcd.com/feed.atom",
		"What if?",
		timeParserPublished(time.RFC3339),
		scrapers.ScrapeXkcdWhatIf,
	},
	{
		"https://www.oglaf.com/feeds/rss/",
		"Oglaf",
		timeParserPublished(time.RFC1123Z),
		scrapers.ScrapeOglaf,
	},
	{
		"http://www.menagea3.net/comic.rss",
		"Ma3",
		timeParserPublished(time.RFC1123Z),
		scrapers.ScrapeMa3,
	},
	{
		"http://www.stickydillybuns.com/comic.rss",
		"SDB",
		timeParserPublished(time.RFC1123Z),
		scrapers.ScrapeMa3,
	},
	{
		"http://www.giantitp.com/comics/oots.rss",
		"OotS",
		scrapers.TimeParserOots,
		scrapers.ScrapeOots,
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

		var lastItemTime time.Time
		var lastItemId string
		var comics []scrapers.ComicUpdate

		// the heavy scraping routines and some memory usage are avoided if no users are waiting for
		// updates on this comic, only the update times are kept up to date
		hasReaders := len(conf.GetUsers(f.name)) > 0
		hasUpdates := false

		if hasReaders {
			comics = make([]scrapers.ComicUpdate, 0, len(feed.Items))
		}

		for _, item := range(feed.Items){
			updated, item_id, err := f.scrapeTime(item, feed, f.name)
			if err != nil {
				parseError(f.name, err)
				continue
			}

			if lastItemTime.Before(updated) {
				lastItemTime = updated
				lastItemId = item_id
				hasUpdates = true
			}

			if hasReaders && conf.IsItemNew(f.name, updated) {
				c, err := f.scrapeContent(item, f.name)
				if err != nil {
					parseError(f.name, err)
					continue
				}

				for i := range(c){
					c[i].Updated = updated
				}

				comics = append(comics, c...)
			}
		}

		if hasUpdates {
			if hasReaders {
				notifyComics(bot, f.name, comics)
			}
			if len(lastItemId) > 0 {
				conf.StoreLastItem(f.name, lastItemId)
			}
			conf.UpdateFeed(f.name, lastItemTime)
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

