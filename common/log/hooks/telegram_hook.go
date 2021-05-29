package hooks

import (
	"fmt"
	"net/http"
)

const (
	MethodSendMessage = "sendMessage"
	ParseModeHTML     = "HTML"
)

type TelegramHook struct {
	baseEndpoint string
	token        string
	channel      string
}

func NewTelegramHook(baseEndpoint, token, channel string) *TelegramHook {
	return &TelegramHook{
		baseEndpoint: baseEndpoint,
		token:        token,
		channel:      channel,
	}
}

func (a *TelegramHook) SendMessage(body string) error {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/bot%s/%s", a.baseEndpoint, a.token, MethodSendMessage), nil)
	if err != nil {
		return err
	}
	q := req.URL.Query()
	q.Set("chat_id", a.channel)
	q.Set("parse_mode", ParseModeHTML)
	q.Set("text", body)
	req.URL.RawQuery = q.Encode()
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad response: [code:%v]", resp.StatusCode)
	}
	return err
}
