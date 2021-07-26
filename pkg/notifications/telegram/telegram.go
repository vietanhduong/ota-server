package telegram

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/vietanhduong/ota-server/pkg/cerrors"
	"github.com/vietanhduong/ota-server/pkg/logger"
	"io/ioutil"
	"net/http"
)

const BaseAPI = "https://api.telegram.org"
const DefaultParseMode = "Markdown"

type Telegram struct {
	token   string
	groupId string
	baseUrl string
}

func InitializeTelegram(botToken, groupId string) *Telegram {
	return &Telegram{
		token:   botToken,
		groupId: groupId,
		baseUrl: getBaseURL(botToken),
	}
}
// SendMessage to Telegram. Default parse_mode is `Markdown`
// to write a new line just add '\n' to your message
// refer: https://rdrr.io/cran/telegram.bot/man/sendMessage.html
func (t *Telegram) SendMessage(message string) error {
	// init url
	url := fmt.Sprintf("%s/sendMessage", t.baseUrl)

	p := RequestPayload{
		ChatId:    t.groupId,
		Text:      message,
		ParseMode: DefaultParseMode,
	}
	// parse payload struct to json
	payload, err := json.Marshal(p)
	if err != nil {
		return err
	}
	// init request to Telegram API
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(payload))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	// make a request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer cerrors.Close(resp.Body)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var result *ResponseBody
	err = json.Unmarshal(body, &result)
	if err != nil {
		return err
	}

	if !result.Ok {
		logger.Logger.Errorf("send message to Telegram failed response body: %s", body)
		return errors.New("send message to Telegram failed")
	}

	return nil
}

func getBaseURL(token string) string {
	return fmt.Sprintf("%s/bot%s", BaseAPI, token)
}
