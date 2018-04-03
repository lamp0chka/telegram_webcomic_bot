package sources

import (
	"gopkg.in/tucnak/telebot.v2"
	"telegram_webcomic_bot/configs"
	"telegram_webcomic_bot/util"
	"fmt"
)

type comicUpdate struct {
	source string
	title string
	url string
	alt string
	postLinkOnly bool
}

func (u *comicUpdate) send(bot *telebot.Bot, uid int64) {
	bot.Send(
		&telebot.Chat{ID: uid},
		fmt.Sprintf("%s\n<a href=\"%s\">%s</a>", u.source, u.url, u.title),
		telebot.ModeHTML,
	)
	if !u.postLinkOnly {
		bot.Send(
			&telebot.Chat{ID: uid},
			u.alt,
		)
	}
}

func notifyComics(bot *telebot.Bot, source string, comics []comicUpdate) {
	conf := configs.GetConfigs()
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