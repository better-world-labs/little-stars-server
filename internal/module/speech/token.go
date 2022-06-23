package speech

import (
	"aed-api-server/internal/pkg/cache"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

var AidSmsPathTokenKeyPrefix = "AID_SMS_PATH_TOKEN"

type TokenGenerator interface {
	Generate() string
}

type TokenValidator interface {
	ValidateToken(token, value string) (bool, error)
}

type TokenStore interface {
	PutToken(token, value string) error
	RemoveToken(token string) (int64, error)
}

type TokenService interface {
	TokenGenerator
	TokenValidator
	TokenStore
}

type tokenService struct {
}

func NewTokenService() TokenService {
	return &tokenService{}
}

func (t tokenService) Generate() string {
	u := uuid.NewString()
	str := []byte(u)[:6]
	return base64.StdEncoding.EncodeToString(str)
}

func (t tokenService) getKey(token string) string {
	return fmt.Sprintf("%s_%s", AidSmsPathTokenKeyPrefix, token)
}

func (t tokenService) ValidateToken(token, value string) (bool, error) {
	conn := cache.GetConn()
	defer conn.Close()

	reply, err := conn.Do("GET", t.getKey(token))
	if err != nil {
		return false, err
	}

	if reply == nil {
		return false, nil
	}

	replyStr := string(reply.([]byte))
	if replyStr != value {
		return false, nil
	}

	return true, nil
}

func (t tokenService) PutToken(token, value string) error {
	conn := cache.GetConn()
	defer conn.Close()

	log.Infof("PutToken: token=%s, value=%s", token, value)
	reply, err := conn.Do("SET", t.getKey(token), value, "PX", 3600000*24)
	if err != nil {
		return err
	}

	if reply != "OK" {
		return errors.New(fmt.Sprintf("PutToken error: reply=%s", reply))
	}

	return nil
}

func (t tokenService) RemoveToken(token string) (int64, error) {
	conn := cache.GetConn()
	defer conn.Close()

	log.Infof("RomoveToken: token=%s", token)
	reply, err := conn.Do("DEL", t.getKey(token))
	if err != nil {
		return 0, err
	}

	return reply.(int64), nil
}
