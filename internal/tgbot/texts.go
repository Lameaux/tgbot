package tgbot

const txtMessage = `
From: %s
To: %s

%s
`

const txtHelp = `
Welcome to Euromoby Sandbox!

To register your MSISDN:
/msisdn <msisdn>
Example: /msisdn 420123456789

To send SMS:
/sms <shortcode> <text>
Example: /sms 9999 Hello World

To confirm SMS delivery:
reply with /ack
`

const (
	txtUnknownCmd       = `Unknown command. Type /help to get the list of supported commands.`
	txtInvalidCmdFmtSMS = `
Invalid format.
Please use "/sms <shortcode> <text>"
Example: /sms 9999 Hello World
`
)

const txtRegisterMSISDN = `
To send SMS you need to register MSISDN first.
Please use "/msisdn <msisdn>"
Example: /msisdn 420123456789
`

const (
	txtInvalidMSISDN    = `%s is invalid. Please use international format, e.g. 420123456789.`
	txtInvalidShortcode = `%s is invalid. Please use digits only, e.g. 9999.`
	txtInvalidAck       = `You need to reply to incoming message with /ack`
)
