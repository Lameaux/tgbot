package main

import (
	"github.com/Lameaux/core/logger"
	"github.com/Lameaux/tgbot/internal/config"
	"github.com/Lameaux/tgbot/internal/connectors"
	"github.com/Lameaux/tgbot/internal/tgbot"
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
