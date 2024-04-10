package toolkit

import (
	"fmt"
	"github.com/alabuta-source/toolkit/rest"
)

const (
	sendMessageURL  = "https://api.telegram.org/bot%s/sendMessage"
	JsonContentType = "application/json"
)

type SendMessageBody struct {
	chatID              string
	text                string
	parseMode           string
	disableNotification bool
	linkPreviewOptions  LinkPreviewOptions
	contentType         string
	botToken            string
}

type messageBody struct {
	ChatID              string             `json:"chat_id"`
	Text                string             `json:"text"`
	ParseMode           string             `json:"parse_mode"`
	DisableNotification bool               `json:"disable_notification"`
	LinkPreviewOptions  LinkPreviewOptions `json:"link_preview_options"`
}

type Option func(body *SendMessageBody)

func BotMessageWithChatID(id string) Option {
	return func(body *SendMessageBody) {
		body.chatID = id
	}
}

func BotMessageWithMessage(message string) Option {
	return func(body *SendMessageBody) {
		body.text = message
	}
}

func BotMessageWithParseMode(parseMode string) Option {
	return func(body *SendMessageBody) {
		body.parseMode = parseMode
	}
}

func BotMessageWithNotificationDisableOn(disable bool) Option {
	return func(body *SendMessageBody) {
		body.disableNotification = disable
	}
}

func BotMessageWithPreviewURLOn(disable bool) Option {
	return func(body *SendMessageBody) {
		body.linkPreviewOptions.Disabled = disable
	}
}

func BotMessageWithContentType(contentType string) Option {
	return func(body *SendMessageBody) {
		body.contentType = contentType
	}
}

func BotMessageWithBotToken(token string) Option {
	return func(body *SendMessageBody) {
		body.botToken = token
	}
}

type LinkPreviewOptions struct {
	Disabled bool `json:"is_disabled"`
}

func NewBotMessage(options ...Option) *SendMessageBody {
	var body SendMessageBody
	for _, opt := range options {
		opt(&body)
	}
	return &body
}

func (body *SendMessageBody) Send() error {
	url := fmt.Sprintf(sendMessageURL, body.botToken)

	contentType := JsonContentType
	if body.contentType != "" {
		contentType = body.contentType
	}
	headers := map[string]string{
		"Content-Type": contentType,
	}
	client := rest.NewRestClient(5)
	return client.POST(url, body.withBotParams(), headers, nil)
}

func (body *SendMessageBody) withBotParams() *messageBody {
	return &messageBody{
		ChatID:              body.chatID,
		Text:                body.text,
		ParseMode:           body.parseMode,
		DisableNotification: body.disableNotification,
		LinkPreviewOptions:  body.linkPreviewOptions,
	}
}
