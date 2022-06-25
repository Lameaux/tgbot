package inputs

import (
	coremodels "github.com/Lameaux/core/models"
)

type SendMessageParams struct {
	MessageID           string            `json:"message_id,omitempty"`
	Sender              string            `json:"sender"`
	MSISDN              string            `json:"msisdn"`
	Recipient           coremodels.MSISDN `json:"-"`
	Body                string            `json:"body"`
	ClientTransactionID string            `json:"client_transaction_id"`
}
