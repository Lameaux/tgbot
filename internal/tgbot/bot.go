package tgbot

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"euromoby.com/tgbot/internal/config"
	"euromoby.com/tgbot/internal/connectors"
	"euromoby.com/tgbot/internal/models"
	"euromoby.com/tgbot/internal/repos"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	coremodels "euromoby.com/core/models"
	coreutils "euromoby.com/core/utils"
)

type SandboxBot struct {
	app      *config.App
	gw       *connectors.GatewayConnector
	botAPI   *tgbotapi.BotAPI
	chatRepo *repos.ChatRepo
}

func NewSandboxBot(app *config.App, gw *connectors.GatewayConnector) (*SandboxBot, error) {
	botAPI, err := tgbotapi.NewBotAPI(app.BotToken)
	if err != nil {
		return nil, err
	}

	botAPI.Debug = true

	sandboxBot := SandboxBot{
		app:      app,
		gw:       gw,
		botAPI:   botAPI,
		chatRepo: repos.NewChatRepo(),
	}

	return &sandboxBot, nil
}

func (b *SandboxBot) StartListener() {
	log.Printf("Listening on account %s", b.botAPI.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.botAPI.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			log.Printf("Unsupported update: %v", update)
			continue
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		var responseMessage string

		switch update.Message.Command() {
		case "help", "start":
			responseMessage = b.commandHelp()
		case "msisdn":
			responseMessage = b.commandMSISDN(update.Message)
		case "sms":
			responseMessage = b.commandSMS(update.Message)
		case "ack":
			responseMessage = b.commandAck(update.Message)
		default:
			responseMessage = b.commandUnknown()
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, responseMessage)
		msg.ReplyToMessageID = update.Message.MessageID
		b.botAPI.Send(msg)
	}
}

func (b *SandboxBot) commandUnknown() string {
	return txtUnknownCmd
}

func (b *SandboxBot) commandHelp() string {
	return txtHelp
}

func (b *SandboxBot) commandMSISDN(message *tgbotapi.Message) string {
	chatID := message.Chat.ID
	rawMSISDN := message.CommandArguments()

	msisdn, err := coremodels.NormalizeMSISDN(rawMSISDN)
	if err != nil {
		return fmt.Sprintf(txtInvalidMSISDN, rawMSISDN)
	}

	b.chatRepo.LinkMSISDNAndChat(msisdn, models.ChatID(chatID))

	return fmt.Sprintf("%s was linked to this chat", msisdn)
}

func (b *SandboxBot) commandAck(message *tgbotapi.Message) string {
	if message.ReplyToMessage == nil {
		return txtInvalidAck
	}

	messageID := message.ReplyToMessage.MessageID

	d := models.MessageDelivery{
		MessageID: strconv.Itoa(messageID),
	}

	err := b.gw.AckMessage(&d)
	if err != nil {
		log.Println(err)

		return "Oops! Ack failed. Please try again later."
	}

	return fmt.Sprintf("Message #%d acked", messageID)
}

func (b *SandboxBot) commandSMS(message *tgbotapi.Message) string {
	shortcodeAndText := strings.SplitN(message.CommandArguments(), " ", 2)

	if len(shortcodeAndText) != 2 {
		return txtInvalidCmdFmtSMS
	}

	shortcode := shortcodeAndText[0]
	if !coremodels.ValidateShortcode(shortcode) {
		return fmt.Sprintf(txtInvalidShortcode, shortcode)
	}

	text := shortcodeAndText[1]

	chatID := models.ChatID(message.Chat.ID)
	msisdn, exists := b.chatRepo.GetMsisdnByChat(chatID)
	if !exists {
		return txtRegisterMSISDN
	}

	m := models.OutcomingMessage{
		MessageID: fmt.Sprintf("%d:%d", chatID, message.MessageID),
		Shortcode: shortcode,
		MSISDN:    msisdn,
		Body:      text,
	}

	err := b.gw.SendMessage(&m)
	if err != nil {
		log.Println(err)

		return "Oops! Sending message failed. Please try again later."
	}

	return fmt.Sprintf(txtMessage, msisdn, shortcode, text)
}

func (b *SandboxBot) StopListener() {
	b.botAPI.StopReceivingUpdates()
}

func (b *SandboxBot) SendMessage(sender string, msisdn coremodels.MSISDN, body string) (*tgbotapi.Message, error) {
	text := fmt.Sprintf(txtMessage, sender, msisdn, body)

	chatID, exists := b.chatRepo.GetChatByMSISDN(msisdn)
	if exists {
		return b.sendToChat(chatID, 0, text)
	}

	return b.sendToChannel(text)
}

func (b *SandboxBot) SendStatus(messageID string, status string) error {
	chatAndReplyTo := strings.Split(messageID, ":")

	if len(chatAndReplyTo) != 2 {
		return models.ErrRecipientNotFound
	}

	chatID, err := coreutils.ParseInt64(chatAndReplyTo[0])
	if err != nil {
		return err
	}

	replyToMessageID, err := strconv.Atoi(chatAndReplyTo[1])
	if err != nil {
		return err
	}

	text := "Status: " + status
	_, err = b.sendToChat(models.ChatID(chatID), replyToMessageID, text)
	return err
}

func (b *SandboxBot) sendToChat(chatID models.ChatID, replyToMessageID int, text string) (*tgbotapi.Message, error) {
	msgConfig := tgbotapi.NewMessage(int64(chatID), text)
	msgConfig.ReplyToMessageID = replyToMessageID

	msg, err := b.botAPI.Send(msgConfig)
	return &msg, err
}

func (b *SandboxBot) sendToChannel(text string) (*tgbotapi.Message, error) {
	msgConfig := tgbotapi.NewMessageToChannel(b.app.ChannelUsername, text)

	msg, err := b.botAPI.Send(msgConfig)
	return &msg, err
}
