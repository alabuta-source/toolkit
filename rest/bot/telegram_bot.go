package bot

import (
	"fmt"
	"net/http"
	"os"

	"github.com/alabuta-source/toolkit/rest"
)

const (
	sendMessageURL    = "https://api.telegram.org/bot%s/sendMessage"
	MarkDownParseMode = "Markdown"
	HTMLParseMode     = "HTML"
)

type sendMessageBody struct {
	ChatID              string             `json:"chat_id"`
	Text                string             `json:"text"`
	ParseMode           string             `json:"parse_mode"`
	DisableNotification bool               `json:"disable_notification"`
	LinkPreviewOptions  LinkPreviewOptions `json:"link_preview_options"`
}

type LinkPreviewOptions struct {
	Disabled bool `json:"is_disabled"`
}

func SendTelegramMessage(chatID, message, parseMode string, disableNotification bool) error {
	token := os.Getenv("BOT_TOKEN")
	url := fmt.Sprintf(sendMessageURL, token)
	body := sendMessageBody{
		ChatID:              chatID,
		Text:                message,
		ParseMode:           parseMode,
		DisableNotification: disableNotification,
		LinkPreviewOptions: LinkPreviewOptions{
			Disabled: true,
		},
	}

	return rest.
		NewHttpClient().
		BuildRequest(
			url,
			http.MethodPost,
			rest.RequestWithBody(&body),
		).
		Execute()
}
