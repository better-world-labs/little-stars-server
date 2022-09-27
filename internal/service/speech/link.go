package speech

import (
	"aed-api-server/internal/interfaces"
	"aed-api-server/internal/pkg/cache"
	"encoding/base64"
	"fmt"
	"github.com/google/uuid"
	"time"
)

const (
	LinkKeyPrefix = "LINK"
)

func CreateLink(url string) (string, error) {
	conn := cache.GetConn()
	defer conn.Close()

	code := genericLinkCode()
	_, err := conn.Do("SET", key(code), url, "EX", 30*24*3600)
	if err != nil {
		return "", err
	}

	return code, nil
}

func GetLink(code string) (string, error) {
	conn := cache.GetConn()
	defer conn.Close()

	reply, err := conn.Do("GET", key(code))
	if err != nil {
		return "", err
	}

	if reply == nil {
		return "", nil
	}

	return string(reply.([]byte)), nil
}

func key(code string) string {
	return fmt.Sprintf("%s-%s-%s", interfaces.GetConfig().Server.Env, LinkKeyPrefix, code)
}

func genericLinkCode() string {
	u := uuid.NewString()
	str := []byte(u)[:6]
	return base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s%d", str, time.Now().Unix())))
}
