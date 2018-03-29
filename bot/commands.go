package bot

import (
	"gopkg.in/tucnak/telebot.v2"
	"telegram_webcomic_bot/configs"
	"telegram_webcomic_bot/util"
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

func setup(bot *telebot.Bot, m *telebot.Message) {
	if !m.Private() {
		return
	}

	conf := configs.GetConfigs()
	sources := conf.GetFeedSources()

	inlineKeys := util.CreateInlineKbd(bot, m.Sender.ID, sources)

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
		conf.CreateUser(m.Sender.ID)
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