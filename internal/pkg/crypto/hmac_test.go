package crypto_test

import (
	"aed-api-server/internal/pkg/crypto"
	"encoding/base64"
	"fmt"
	"testing"
)

func Test_HmacSHA256(t *testing.T) {
	hash := crypto.HmacSHA256([]byte("zX4YmvPuTIdu60reFNAIWA=="), []byte(""))
	fmt.Println(base64.StdEncoding.EncodeToString(hash))
}
