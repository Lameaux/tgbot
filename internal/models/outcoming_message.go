package models

import (
	coremodels "euromoby.com/core/models"
)

type OutcomingMessage struct {
	MessageID string            `json:"message_id"`
	Shortcode string            `json:"shortcode"`
	MSISDN    coremodels.MSISDN `json:"msisdn"`
	Body      string            `json:"body"`
}
