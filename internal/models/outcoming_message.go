package models

import (
	coremodels "github.com/Lameaux/core/models"
)

type OutcomingMessage struct {
	MessageID string            `json:"message_id"`
	Shortcode string            `json:"shortcode"`
	MSISDN    coremodels.MSISDN `json:"msisdn"`
	Body      string            `json:"body"`
}
