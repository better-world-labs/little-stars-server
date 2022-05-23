package asserts

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

var resources = map[string][]byte{}

func LoadResourceDir(dir string) error {
	files, err := os.ReadDir(dir)

	if err != nil {
		return err
	}

	for _, file := range files {
		filePath := fmt.Sprintf("%s/%s", dir, file.Name())
		if file.IsDir() {
			err := LoadResourceDir(filePath)
			if err != nil {
				return err
			}

			continue
		}

		bytes, err := ioutil.ReadFile(filePath)
		if err != nil {
			return err
		}

		resources[file.Name()] = bytes
	}

	return nil
}

func GetResource(key string) ([]byte, bool) {
	v, exists := resources[key]
	return v, exists
}

func GetResourceFromUrl(url string) ([]byte, error) {
	resposne, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	if resposne.StatusCode != 200 {
		return nil, errors.New("invalid request for url")
	}

	defer resposne.Body.Close()
	b, err := ioutil.ReadAll(resposne.Body)
	return b, err
}
