package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	coremodels "github.com/Lameaux/core/models"
	coreviews "github.com/Lameaux/core/views"
	"github.com/Lameaux/tgbot/internal/config"
	"github.com/Lameaux/tgbot/internal/inputs"
	"github.com/Lameaux/tgbot/internal/models"
	"github.com/Lameaux/tgbot/internal/tgbot"
	"github.com/gin-gonic/gin"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type MessageHandler struct {
	app *config.App
	bot *tgbot.SandboxBot
}

func NewMessageHandler(app *config.App, bot *tgbot.SandboxBot) *MessageHandler {
	return &MessageHandler{app, bot}
}

func (h *MessageHandler) SendMessage(c *gin.Context) {
	p, err := h.parseRequest(c)
	if err != nil {
		coreviews.ErrorJSON(c, http.StatusBadRequest, err)

		return
	}

	msg, err := h.bot.SendMessage(p.Sender, p.Recipient, p.Body)
	if err != nil {
		h.handleError(c, err)
		return
	}

	p.MessageID = fmt.Sprintf("%d:%d", msg.Chat.ID, msg.MessageID)

	c.JSON(http.StatusCreated, p)
}

func (h *MessageHandler) parseRequest(c *gin.Context) (*inputs.SendMessageParams, error) {
	var p inputs.SendMessageParams

	dec := json.NewDecoder(c.Request.Body)

	dec.DisallowUnknownFields()

	if err := dec.Decode(&p); err != nil {
		return nil, err
	}

	msisdn, err := coremodels.NormalizeMSISDN(p.MSISDN)
	if err != nil {
		return nil, err
	}
	p.Recipient = msisdn

	return &p, nil
}

func (h *MessageHandler) handleError(c *gin.Context, err error) {
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
