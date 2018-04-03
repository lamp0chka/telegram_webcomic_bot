package sources

import (
	"gopkg.in/tucnak/telebot.v2"
	"time"
	"log"
	"telegram_webcomic_bot/configs"
	"github.com/mmcdole/gofeed"
	"runtime/debug"
	"fmt"
	"sort"
	"telegram_webcomic_bot/util"
)

type comicUpdate struct {
	source string
	title string
	url string
	alt string
	updated time.Time
}

func (u *comicUpdate) send(bot *telebot.Bot, uid int64) {
	bot.Send(
		&telebot.Chat{ID: uid},
		fmt.Sprintf("%s\n<a href=\"%s\">%s</a>", u.source, u.url, u.title),
		telebot.ModeHTML,
	)
	if len(u.alt) > 0 {
		bot.Send(
			&telebot.Chat{ID: uid},
			u.alt,
		)
	}
}

type feedSrc struct {
	url string
	name string
	scrapeTime func(*gofeed.Item)(time.Time, error)
	scrapeContent func(*gofeed.Item, feedSrc)([]comicUpdate, error)
}

func timeParserUpdated(timespec string) (func(*gofeed.Item)(time.Time, error)){
	return func(item *gofeed.Item) (time.Time, error){
		return time.Parse(timespec, item.Updated)
	}
}

func timeParserPublished(timespec string) (func(*gofeed.Item)(time.Time, error)){
	return func(item *gofeed.Item) (time.Time, error){
		return time.Parse(timespec, item.Published)
	}
}

func parseError(source string, err error) {
	log.Printf("Error parsing %s: %s", source, err.Error())
	debug.PrintStack()
}

var feedScrapers = []feedSrc{
	{
		"http://dilbert.com/feed",
		"Dilbert",
		timeParserUpdated(time.RFC3339),
		scrapeDilbert,
	},
	{
		"https://www.xkcd.com/rss.xml",
		"xkcd",
		timeParserPublished(time.RFC1123Z),
		scrapeXkcd,
	},
	{
		"https://what-if.xkcd.com/feed.atom",
		"What if?",
		timeParserPublished(time.RFC3339),
		scrapeXkcdWhatIf,
	},
	{
		"https://www.oglaf.com/feeds/rss/",
		"Oglaf",
		timeParserPublished(time.RFC1123Z),
		scrapeOglaf,
	},
	{
		"http://www.menagea3.net/comic.rss",
		"Ma3",
		timeParserPublished(time.RFC1123Z),
		scrapeMa3,
	},
}

func notifyComics(bot *telebot.Bot, source string, comics []comicUpdate) {
	conf := configs.GetConfigs()

	sort.Slice(comics, func(i, j int) bool {
		return comics[i].updated.Before(comics[j].updated)
	})

	for _, id := range conf.GetUsers(source) {
		uid := int64(id)
		for _, c := range comics {
			c.send(bot, uid)
		}
	}
}

func notifyNewSources(bot *telebot.Bot, sources []string) {
	conf := configs.GetConfigs()

	for _, id := range conf.GetAllUsers() {
		kbd := util.CreateInlineKbd(bot, id, sources)
		bot.Send(&telebot.Chat{ID: int64(id)}, "New webomics available!", &telebot.ReplyMarkup{
			InlineKeyboard: kbd,
		})
	}
}

func updateFeeds(bot *telebot.Bot) {
	conf := configs.GetConfigs()
	fp := gofeed.NewParser()

	for _, f := range(feedScrapers) {
		if len(conf.GetUsers(f.name)) > 0 {
			feed, err := fp.ParseURL(f.url)
			if err != nil {
				parseError(f.name, err)
				continue
			}

			comics := make([]comicUpdate, 0, len(feed.Items))
			var lastItemTime time.Time

			for _, item := range(feed.Items){
				updated, err := f.scrapeTime(item)
				if err != nil {
					parseError(f.name, err)
					continue
				}

				if lastItemTime.Before(updated) {
					lastItemTime = updated
				}

				if conf.IsItemNew(f.name, updated) {
					c, err := f.scrapeContent(item, f)
					if err != nil {
						parseError(f.name, err)
						continue
					}

					for i := range(c){
						c[i].updated = updated
					}

					comics = append(comics, c...)
				}
			}

			if len(comics) > 0 {
				notifyComics(bot, f.name, comics)
				conf.UpdateFeed(f.name, lastItemTime)
			}
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

