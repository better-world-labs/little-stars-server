package user

import (
	"aed-api-server/internal/interfaces/entities"
	"aed-api-server/internal/pkg/cache"
	"encoding/json"
	"fmt"
)

type CryptKeyCache struct {
	Env string `conf:"server.env"`
}

const (
	Key = "UserEncryptKey"
)

//go:inject-component
func NewCryptKeyCache() *CryptKeyCache {
	return &CryptKeyCache{}
}

func (c *CryptKeyCache) parseKey(userId int64, version int) string {
	return fmt.Sprintf("%s_%s_%d_%d", c.Env, Key, userId, version)
}

func (c *CryptKeyCache) getKeyFromCache(userId int64, version int) (*entities.WechatEncryptKey, error) {
	conn := cache.GetConn()
	defer conn.Close()

	reply, err := conn.Do("GET", c.parseKey(userId, version))
	if err != nil {
		return nil, err
	}

	if reply == nil {
		return nil, nil
	}

	uint8s := reply.([]uint8)

	var res entities.WechatEncryptKey
	return &res, json.Unmarshal(uint8s, &res)

}

func (c *CryptKeyCache) PutKeys(userId int64, keys []*entities.WechatEncryptKey) error {
	conn := cache.GetConn()
	defer conn.Close()

	for _, k := range keys {
		j, err := json.Marshal(k)
		if err != nil {
			return err
		}

		expiresIn := k.ExpiredInSecond() + 1800
		_, err = conn.Do("SET", c.parseKey(userId, k.Version), string(j), "EX", int(expiresIn)) // 延时一秒，以应对前端拿到刚好过期的key的情况
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *CryptKeyCache) GetKey(userId int64, version int) (*entities.WechatEncryptKey, error) {
	return c.getKeyFromCache(userId, version)
}
