package connectors

import (
	"bytes"
	"encoding/json"
	"net/http"

	"euromoby.com/tgbot/internal/models"
)

const (
	apiBaseURL = "http://0.0.0.0:8080/v1/sms/providers/sandbox"
)

type GatewayConnector struct{}

func NewGatewayConnector() *GatewayConnector {
	return &GatewayConnector{}
}

func (c *GatewayConnector) SendMessage(m *models.OutcomingMessage) error {
	jsonBody, err := json.Marshal(m)
	if err != nil {
		return err
	}

	httpResp, err := http.Post(apiBaseURL+"/inbound/message", "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		return err
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode != 201 {
		return models.ErrSendFailed
	}

	return nil
}

func (c *GatewayConnector) AckMessage(d *models.MessageDelivery) error {
	jsonBody, err := json.Marshal(d)
	if err != nil {
		return err
	}

	httpResp, err := http.Post(apiBaseURL+"/outbound/ack", "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		return err
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode != 200 {
		return models.ErrSendFailed
	}

	return nil
}
