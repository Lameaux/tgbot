package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/Lameaux/tgbot/internal/config"
	"github.com/Lameaux/tgbot/internal/handlers"
	"github.com/Lameaux/tgbot/internal/tgbot"

	coremiddlewares "github.com/Lameaux/core/middlewares"
)

func Gin(app *config.App, b *tgbot.SandboxBot) *gin.Engine {
	r := gin.Default()
	r.Use(coremiddlewares.Timeout(app.Config.WaitTimeout))

	i := handlers.NewIndexHandler()

	r.GET("/", i.Index)
	r.GET("/health", i.Index)

	m := handlers.NewMessageHandler(app, b)
	s := handlers.NewStatusHandler(app, b)

	sandbox := r.Group("/sandbox")
	{
		sandbox.POST("/message", m.SendMessage)
		sandbox.POST("/status", s.SendStatus)
	}

	return r
}
