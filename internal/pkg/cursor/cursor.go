package cursor

import (
	"encoding/base64"
	"encoding/json"
)

func ToString(cursor interface{}) (string, error) {
	str, err := json.Marshal(cursor)
	return base64.StdEncoding.EncodeToString(str), err
}

func FromString(str string, cursor interface{}) error {
	jsonStr, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		return err
	}
	return json.Unmarshal(jsonStr, &cursor)
}
