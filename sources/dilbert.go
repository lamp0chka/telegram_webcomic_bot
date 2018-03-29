package sources

import (
	"gopkg.in/tucnak/telebot.v2"
	"telegram_webcomic_bot/configs"
	"time"
)

func scrapeDilbert(bot *telebot.Bot) {

	configs.GetConfigs().UpdateFeed("Dilbert", time.Now())
}
