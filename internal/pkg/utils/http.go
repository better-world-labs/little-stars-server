package utils

import (
	"bytes"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"net/http"
)

func Post(url string, data interface{}, rst interface{}) error {
	jsonStr, err := json.Marshal(data)
	if err != nil {
		return err
	}

	res, err := http.Post(url, "application/json", bytes.NewBuffer(jsonStr))
	if err != nil {
		return err
	}
	return readHttpResAsJson(res.Body, rst)
}

func Get(url string, rst interface{}) error {
	res, err := http.Get(url)
	if err != nil {
		return err
	}
	return readHttpResAsJson(res.Body, rst)
}

func GetWithClient(client http.Client, url string, rst interface{}) error {
	res, err := client.Get(url)
	if err != nil {
		return err
	}
	return readHttpResAsJson(res.Body, rst)
}

func readHttpResAsJson(reader io.ReadCloser, rst interface{}) error {
	defer func() {
		err := reader.Close()
		if err != nil {
			log.Error("readHttpRes error:", err)
		}
	}()

	content, err := ioutil.ReadAll(reader)

	if err != nil {
		return err
	}
	return json.Unmarshal(content, rst)
}
