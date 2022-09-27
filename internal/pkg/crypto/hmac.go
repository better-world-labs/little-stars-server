package crypto

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

func HmacSHA256(key, data []byte) []byte {
	hash := hmac.New(sha256.New, key)
	hash.Write(data)
	return hash.Sum(nil)
}

func HmacSHA256HexString(key, data []byte) string {
	return hex.EncodeToString(HmacSHA256(key, data))
}
