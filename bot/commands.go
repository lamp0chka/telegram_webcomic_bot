package bot

import (
	"gopkg.in/tucnak/telebot.v2"
	"telegram_webcomic_bot/configs"
	"fmt"
)

func help(bot *telebot.Bot, m *telebot.Message) {
	msg := "I can keep you up to date with some webcomics that i know of.\n" +
			"Use /setup command to tell me which updates you want.\n" +
			"/help command shows this message.\n" +
			"\n" +
			"If you want me to keep track of a new webcomic leave a ticket on " +
			"https://github.com/mellotanica/telegram_webcomic_bot/issues"
	bot.Send(m.Sender, msg)
}

func handleSetupBtn(bot *telebot.Bot,
		m *telebot.Message,
		btn *telebot.InlineButton,
		kbd [][]telebot.InlineButton,
		source string) func(c *telebot.Callback) {
	return func(c *telebot.Callback) {
		conf := configs.GetConfigs()
		enabled := conf.UserToggleSource(m.Sender.ID, source)
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

func setup(bot *telebot.Bot, m *telebot.Message) {
	if !m.Private() {
		return
	}

	conf := configs.GetConfigs()
	sources := conf.GetFeedSources()

	inlineKeys := make([][]telebot.InlineButton, len(sources))
	for i, s := range(sources) {
		var text string
		if conf.UserSourceEnabled(m.Sender.ID, s) {
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
		bot.Handle(&btn, handleSetupBtn(bot, m, &inlineKeys[i][0], inlineKeys, s))
	}

	if len(sources) > 0 {
		bot.Send(m.Sender, "These are the available comics:", &telebot.ReplyMarkup{
			InlineKeyboard: inlineKeys,
		})
	} else {
		bot.Send(m.Sender, "I don't know of any webcomic at this time")
	}

}

func start(bot *telebot.Bot, m *telebot.Message) {
	conf := configs.GetConfigs()
	_, ok := conf.GetUser(m.Sender.ID)
	help(bot, m)
	if !ok {
		bot.Send(m.Sender, "First you need to set some comics to follow")
		setup(bot, m)
	}
}

func msg_handler(bot *telebot.Bot, f func(bot *telebot.Bot, m *telebot.Message)) (func(*telebot.Message)) {
	return func(m *telebot.Message) {
		f(bot, m)
	}
}

func setupCommands(bot *telebot.Bot) {
	bot.Handle("/start", msg_handler(bot, start))
	bot.Handle("/help", msg_handler(bot, help))
	bot.Handle("/setup", msg_handler(bot, setup))
}