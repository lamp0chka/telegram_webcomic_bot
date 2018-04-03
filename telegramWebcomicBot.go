package main

import (
	"telegram_webcomic_bot/configs"
	"log"
	"gopkg.in/tucnak/telebot.v2"
	"time"
	"telegram_webcomic_bot/bot"
)

func main() {
	config := configs.GetConfigs()
	if len(config.GetToken()) <= 0 {
		log.Fatal("Missing Telegram Token, please edit config.json file first!")
	}

	log.Print("Initializing telegram bot...")

	b, err := telebot.NewBot(telebot.Settings{
		Token: config.GetToken(),
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
		Reporter: func(e error) {
			log.Printf("Telegram Bot Error: %s\n", e.Error())
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	log.Print("Registering command handlers...")

	bot.SetupCommands(b)

	log.Print("Starting prediodic jobs...")

	bot.StartTasks(b)

	log.Print("Start Telegram bot.")

	b.Start()

	log.Print("Telegram bot terminated.")
}
