package imageprocessing

import (
	"errors"
	"net/http"
)

func FromUrl(url string) ([]byte, error) {
	resposne, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	if resposne.StatusCode != 200 {
		return nil, errors.New("invalid request for url")
	}

	var b []byte
	_, err = resposne.Body.Read(b)
	return b, err
}
