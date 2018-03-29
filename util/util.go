package util

import (
	"gopkg.in/tucnak/telebot.v2"
	"fmt"
	"telegram_webcomic_bot/configs"
)

func handleSetupBtn(bot *telebot.Bot,
		dst int,
		btn *telebot.InlineButton,
		kbd [][]telebot.InlineButton,
		source string) func(c *telebot.Callback) {
	return func(c *telebot.Callback) {
		conf := configs.GetConfigs()
		enabled := conf.UserToggleSource(dst, source)
		if enabled {
			btn.Text = fmt.Sprintf("%s (enabled)", source)
		} else {
			btn.Text = fmt.Sprintf("%s (disabled)", source)
		}
		bot.Edit(c.Message, c.Message.Text, &telebot.ReplyMarkup{
			InlineKeyboard: kbd,
		})
		bot.Respond(c, &telebot.CallbackResponse{})
	}
}

func CreateInlineKbd(bot *telebot.Bot, dst int, sources []string) ([][]telebot.InlineButton) {
	inlineKeys := make([][]telebot.InlineButton, len(sources))
	conf := configs.GetConfigs()
	for i, s := range(sources) {
		var text string
		if conf.UserSourceEnabled(dst, s) {
			text = fmt.Sprintf("%s (enabled)", s)
		} else {
			text = fmt.Sprintf("%s (disabled)", s)
		}
		btn := telebot.InlineButton {
			Unique: s,
			Text: text,
		}
		inlineKeys[i] = []telebot.InlineButton{
			btn,
		}
		bot.Handle(&btn, handleSetupBtn(bot, dst, &inlineKeys[i][0], inlineKeys, s))
	}
	return inlineKeys
}
