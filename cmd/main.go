package main

import (
	"context"
	"os/signal"
	"syscall"

	"euromoby.com/core/logger"

	"euromoby.com/tgbot/internal/config"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	app := config.NewApp()
	logger.Infow("Starting", "app", config.AppName, "version", config.AppVersion)

	bot := startBot(app)
	srv := startAPIServer(app, bot)

	// Listen for the interrupt signal.
	<-ctx.Done()
	logger.Infow("the interrupt received, shutting down gracefully, press Ctrl+C again to force")
	stop()

	shutdownBot(bot)
	shutdownAPIServer(srv)

	app.Config.Shutdown()

	logger.Infow("bye")
}
