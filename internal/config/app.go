package config

import (
	"os"

	coreconfig "github.com/Lameaux/core/config"
	coreutils "github.com/Lameaux/core/utils"
)

type App struct {
	Config coreconfig.AppConfig

	BotToken        string
	ChannelUsername string
}

const (
	AppName    = "tgbot"
	AppVersion = "0.2"
)

func defaultApp(env string) *App {
	return &App{
		Config:          *coreconfig.NewAppConfig(env),
		BotToken:        coreutils.GetEnv("BOT_TOKEN"),
		ChannelUsername: coreutils.GetEnv("CHANNEL_USERNAME"),
	}
}

func NewApp() *App {
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = coreconfig.EnvDevelopment
	}

	return defaultApp(env)
}
