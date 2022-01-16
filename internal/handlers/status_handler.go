package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	coreviews "euromoby.com/core/views"
	"euromoby.com/tgbot/internal/config"
	"euromoby.com/tgbot/internal/inputs"
	"euromoby.com/tgbot/internal/models"
	"euromoby.com/tgbot/internal/tgbot"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type StatusHandler struct {
	app *config.App
	bot *tgbot.SandboxBot
}

func NewStatusHandler(app *config.App, bot *tgbot.SandboxBot) *StatusHandler {
	return &StatusHandler{app, bot}
}

func (h *StatusHandler) SendStatus(c *gin.Context) {
	p, err := h.parseRequest(c)
	if err != nil {
		coreviews.ErrorJSON(c, http.StatusBadRequest, err)

		return
	}

	if err = h.bot.SendStatus(p.MessageID, p.Status); err != nil {
		h.handleError(c, err)

		return
	}

	c.JSON(http.StatusCreated, p)
}

func (h *StatusHandler) parseRequest(c *gin.Context) (*inputs.SendStatusParams, error) {
	var p inputs.SendStatusParams

	dec := json.NewDecoder(c.Request.Body)

	dec.DisallowUnknownFields()

	if err := dec.Decode(&p); err != nil {
		return nil, err
	}

	return &p, nil
}

func (h *StatusHandler) handleError(c *gin.Context, err error) {
	if errors.Is(err, models.ErrRecipientNotFound) {
		coreviews.ErrorJSON(c, http.StatusNotFound, err)
		return
	}

	var boterr *tgbotapi.Error
	if errors.As(err, &boterr) {
		coreviews.ErrorJSON(c, boterr.Code, boterr)
	} else {
		coreviews.ErrorJSON(c, http.StatusInternalServerError, err)
	}
}
