package inputs

type SendStatusParams struct {
	MessageID string `json:"message_id"`
	Status    string `json:"status"`
}
