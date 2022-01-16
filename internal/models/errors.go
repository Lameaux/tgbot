package models

import "errors"

var (
	ErrRecipientNotFound = errors.New("recipient not found")
	ErrSendFailed        = errors.New("failed to send")
	ErrInvalidMSISDN     = errors.New("invalid msisdn format")
)
