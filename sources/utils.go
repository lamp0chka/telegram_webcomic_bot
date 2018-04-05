package sources

import (
	"time"
	"gopkg.in/tucnak/telebot.v2"
	"github.com/mmcdole/gofeed"
	"log"
	"runtime/debug"
	"sort"
	"telegram_webcomic_bot/util"
	"telegram_webcomic_bot/configs"
	"telegram_webcomic_bot/sources/scrapers"
)


type feedSrc struct {
	url string
	name string
	scrapeTime func(*gofeed.Item, *gofeed.Feed, string)(time.Time, string, error)
	scrapeContent func(*gofeed.Item, string)([]scrapers.ComicUpdate, error)
}

func timeParserUpdated(timespec string) (func(*gofeed.Item, *gofeed.Feed, string)(time.Time, string, error)){
	return func(item *gofeed.Item, _ *gofeed.Feed, _ string) (time.Time, string, error){
		t, err := time.Parse(timespec, item.Updated)
		return t, "", err
	}
}

func timeParserPublished(timespec string) (func(*gofeed.Item, *gofeed.Feed, string)(time.Time, string, error)){
	return func(item *gofeed.Item, _ *gofeed.Feed, _ string) (time.Time, string, error){
		t, err := time.Parse(timespec, item.Published)
		return t, "", err
	}
}

func parseError(source string, err error) {
	log.Printf("Error parsing %s: %s", source, err.Error())
	debug.PrintStack()
}

func notifyComics(bot *telebot.Bot, source string, comics []scrapers.ComicUpdate) {
	conf := configs.GetConfigs()

	sort.Slice(comics, func(i, j int) bool {
		return comics[i].Updated.Before(comics[j].Updated)
	})

	for _, id := range conf.GetUsers(source) {
		uid := int64(id)
		for _, c := range comics {
			c.Send(bot, uid)
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

