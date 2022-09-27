package feed

import (
	"encoding/base64"
	"encoding/json"
)

type cursorS struct {
	BeforeId int64 `json:"BeforeId"`
}

func parseCursor(cursor string) (beforeId int64, err error) {
	decodeString, err := base64.StdEncoding.DecodeString(cursor)
	if err != nil {
		return 0, err
	}
	var c cursorS
	err = json.Unmarshal(decodeString, &c)
	if err != nil {
		return 0, err
	}
	return c.BeforeId, nil
}

func toCursor(beforeId int64) string {
	s := cursorS{BeforeId: beforeId}
	str, _ := json.Marshal(s)
	return base64.StdEncoding.EncodeToString(str)
}

func limitPageSize(size int) int {
	if size == 0 {
		return DefaultPageSize
	}
	if size > MaxPageSize {
		return MaxPageSize
	}
	return size
}

const MaxSubscribeMessageFieldLen = 12

func stringCutAndEclipse(str string) string {
	if len(str) > MaxSubscribeMessageFieldLen {
		rs := []rune(str)
		return string(rs[:MaxSubscribeMessageFieldLen]) + "..."
	}
	return str
}
