package routes

import (
	"github.com/gin-gonic/gin"

	"euromoby.com/tgbot/internal/config"
	"euromoby.com/tgbot/internal/handlers"
	"euromoby.com/tgbot/internal/tgbot"

	coremiddlewares "euromoby.com/core/middlewares"
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
