package easy

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/monaco-io/request"
	"gopkg.in/h2non/gentleman.v2"
)

type EasyDownloadFile struct {
	api *tgbotapi.BotAPI
}

// DownloadFileStringOld deprecated
func (d *EasyDownloadFile) DownloadFileStringOld(fileId string) (string, error) {
	fileDirectLink, err := d.api.GetFileDirectURL(fileId)
	if err != nil {
		return "", fmt.Errorf("Getting file direct link error: %s", err)
	}

	cli := gentleman.New()
	cli.URL(fileDirectLink)

	req := cli.Request()

	res, err := req.Send()
	if err != nil {
		return "", fmt.Errorf("Downloading sended file error: %s", err)
	}
	if !res.Ok {
		return "", fmt.Errorf("Downloading sended file error invalid server response: %d\n", res.StatusCode)
	}

	return res.String(), nil
}

func (d *EasyDownloadFile) DownloadFileString(fileId string) (string, error) {
	fileDirectLink, err := d.api.GetFileDirectURL(fileId)
	if err != nil {
		return "", fmt.Errorf("Getting file direct link error: %s", err)
	}

	c := request.Client{
		URL:    fileDirectLink,
		Method: "GET",
		//Query: map[string]string{"hello": "world"},
		//JSON:   body,
	}

	response := c.Send()

	if response.Error() != nil {
		return "", response.Error()
	}
	if !response.OK() {
		return "", fmt.Errorf("invalid server response status code while downloading file: %d", response.Code())
	}

	return response.String(), nil
}