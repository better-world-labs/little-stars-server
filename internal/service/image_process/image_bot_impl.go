package image_process

import (
	"errors"
	"github.com/imroc/req"
	log "github.com/sirupsen/logrus"
)

type ImageBot struct {
	Url string `conf:"img-bot-service"`
}

//go:inject-component
func NewImageBotService() *ImageBot {
	return &ImageBot{}
}

func (s ImageBot) Call(tplName string, args map[string]interface{}, saveAs string) (string, error) {
	json := req.BodyJSON(map[string]interface{}{
		"tplName": tplName,
		"args":    args,
		"save":    saveAs,
	})

	log.Info("request", s.Url, "body=", json)
	resp, err := req.Post(s.Url, req.Header{
		"Content-Type": "application/json; charset=utf8",
	}, json)

	if err != nil {
		return "", err
	}

	res := make(Response)
	err = resp.ToJSON(&res)
	if err != nil {
		return "", err
	}

	if !res.Succeed() {
		return "", errors.New(res.Message())
	}

	return res["url"].(string), nil
}

type (
	Response map[string]interface{}
)

func (r Response) Succeed() bool {
	code := r["code"].(float64)
	return code == 200
}

func (r Response) Message() string {
	return r["message"].(string)
}
