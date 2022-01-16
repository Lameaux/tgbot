package main

import (
	"euromoby.com/core/logger"
	"euromoby.com/tgbot/internal/config"
	"euromoby.com/tgbot/internal/connectors"
	"euromoby.com/tgbot/internal/tgbot"
)

func startBot(app *config.App) *tgbot.SandboxBot {
	c := connectors.NewGatewayConnector()

	bot, err := tgbot.NewSandboxBot(app, c)
	if err != nil {
		logger.Fatal(err)
	}

	go bot.StartListener()

	return bot
}

func shutdownBot(bot *tgbot.SandboxBot) {
	bot.StopListener()
}
