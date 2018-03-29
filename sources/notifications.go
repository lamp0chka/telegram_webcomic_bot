package sources

import (
	"gopkg.in/tucnak/telebot.v2"
	"telegram_webcomic_bot/configs"
	"telegram_webcomic_bot/util"
)

func notifyComic(bot *telebot.Bot, source, title, imageUrl, alt string) {

}

func notifyNewSources(bot *telebot.Bot, sources []string) {
	conf := configs.GetConfigs()

	for _, id := range conf.GetUsers() {
		kbd := util.CreateInlineKbd(bot, id, sources)
		bot.Send(&telebot.Chat{ID:int64(id)}, "New webomics available!", &telebot.ReplyMarkup{
			InlineKeyboard: kbd,
		})
	}
}