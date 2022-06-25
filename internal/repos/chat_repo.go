package repos

import (
	coremodels "github.com/Lameaux/core/models"
	"github.com/Lameaux/tgbot/internal/models"
)

var msisdnToChat = map[coremodels.MSISDN]models.ChatID{
	420773468273: 211242258,
	//"35797675605":  1917869903,
	380501231488: 218306959,
}

var chatToMSISDN = map[models.ChatID]coremodels.MSISDN{
	211242258: 420773468273,
	// 1917869903: "35797675605",
	218306959: 380501231488,
}

type ChatRepo struct{}

func NewChatRepo() *ChatRepo {
	return &ChatRepo{}
}

func (r *ChatRepo) GetChatByMSISDN(msisdn coremodels.MSISDN) (models.ChatID, bool) {
	chatID, exists := msisdnToChat[msisdn]
	return chatID, exists
}

func (r *ChatRepo) GetMsisdnByChat(chatID models.ChatID) (coremodels.MSISDN, bool) {
	msisdn, exists := chatToMSISDN[chatID]
	return msisdn, exists
}

func (r *ChatRepo) LinkMSISDNAndChat(msisdn coremodels.MSISDN, chatID models.ChatID) {
	delete(msisdnToChat, msisdn)
	msisdnToChat[msisdn] = chatID

	delete(chatToMSISDN, chatID)
	chatToMSISDN[chatID] = msisdn
}
