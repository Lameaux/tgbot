package main

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/Lameaux/core/logger"
	"github.com/Lameaux/tgbot/internal/config"
	"github.com/Lameaux/tgbot/internal/routes"
	"github.com/Lameaux/tgbot/internal/tgbot"
)

const serverShutdownTimeout = 5 * time.Second

func startAPIServer(app *config.App, bot *tgbot.SandboxBot) *http.Server {
	srv := &http.Server{
		Addr:    ":" + app.Config.Port,
		Handler: routes.Gin(app, bot),
	}

	logger.Infow("starting server", "port", app.Config.Port)

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Fatal("listen: %s\n", err)
		}
	}()

	return srv
}

func shutdownAPIServer(srv *http.Server) {
	logger.Infow("shutting down API server")

	ctx, cancel := context.WithTimeout(context.Background(), serverShutdownTimeout)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("api server forced to shutdown: ", err)
	}

	logger.Infow("api server exiting")
}
