package utils

import (
	"bytes"
	"encoding/json"
	"gitlab.openviewtech.com/openview-pub/gopkg/log"
	"io/ioutil"
	"net/http"
)

type DingTalkMsg struct {
	Msgtype    string     `json:"msgtype"`
	Title      string     `json:"title"`
	ActionCard ActionCard `json:"actionCard"`
	Markdown   Markdown   `json:"markdown"`
}

type ActionCard struct {
	Title          string `json:"title"`
	Text           string `json:"text"`
	BtnOrientation string `json:"btnOrientation"`
	SingleTitle    string `json:"singleTitle"`
	SingleURL      string `json:"singleURL"`
}

type Markdown struct {
	Title string `json:"title"`
	Text  string `json:"text"`
}

func SendDingTalkBot(url string, msg *DingTalkMsg) error {
	jsonStr, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	rst, err := http.Post(url, "application/json", bytes.NewBuffer(jsonStr))
	if err != nil {
		return err
	}

	defer rst.Body.Close()

	body, err := ioutil.ReadAll(rst.Body)

	log.Info("dingTalk send rst:", string(body))
	return nil
}
