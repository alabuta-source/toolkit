package bot

import (
	"fmt"
	"net/http"
	"os"

	"github.com/alabuta-source/toolkit/rest"
)

const (
	sendMessageURL    = "https://api.telegram.org/bot%s/sendMessage?chat_id=%s&text=%s&parse_mode=%s&disable_notification=%t"
	MarkDownParseMode = "Markdown"
	HTMLParseMode     = "HTML"
)

func SendTelegramMessage(chatID, message, parseMode string) error {
	url := fmt.Sprintf(
		sendMessageURL,
		os.Getenv("BOT_TOKEN"),
		chatID,
		message,
		parseMode,
		false,
	)

	return rest.
		NewHttpClient().
		BuildRequest(url, http.MethodPost).
		Execute()
}
