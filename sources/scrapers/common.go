package scrapers

import (
	"time"
	"gopkg.in/tucnak/telebot.v2"
	"fmt"
)

type ComicUpdate struct {
	Source string
	Title string
	Url string
	Alt string
	Updated time.Time
}

func (u *ComicUpdate) Send(bot *telebot.Bot, uid int64) {
	bot.Send(
		&telebot.Chat{ID: uid},
		fmt.Sprintf("%s\n<a href=\"%s\">%s</a>", u.Source, u.Url, u.Title),
		telebot.ModeHTML,
	)
	if len(u.Alt) > 0 {
		bot.Send(
			&telebot.Chat{ID: uid},
			u.Alt,
		)
	}
}
