package bot

import (
	"telegram_webcomic_bot/configs"
	"log"
	tb "gopkg.in/tucnak/telebot.v2"
	"time"
)

func StartTelegramBot() {
	config := configs.GetConfigs()
	if len(config.Token) <= 0 {
		log.Fatal("Missing Telegram Token, please edit config.json file first!")
	}

	log.Print("Initializing telegram bot...")

	b, err := tb.NewBot(tb.Settings{
		Token: config.Token,
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		log.Fatal(err)
	}

	log.Print("Registering command handlers...")

	setupCommands(b)

	log.Print("Start Telegram bot.")

	b.Start()

	log.Print("Telegram bot terminated.")
}