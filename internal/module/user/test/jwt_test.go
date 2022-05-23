package test

import (
	"aed-api-server/internal/module/user"
	"encoding/base64"
	"fmt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gitlab.openviewtech.com/openview-pub/gopkg/crypto2"
	"strings"
	"testing"
	"time"
)

func TestJwt(t *testing.T) {
	u := uuid.New()
	fmt.Println(u.ID())
	expiresIn := int64(120)
	id := int64(49)
	user.InitJwt("this is a testing key", expiresIn)
	fmt.Println(strings.ReplaceAll(uuid.NewString(), "-", ""))
	newString := uuid.NewString()
	fmt.Println(base64.StdEncoding.EncodeToString(crypto2.Sha3Hash([]byte(newString))))
	token, err := user.SignToken(id)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%s\n", token)
	claim, err := user.ParseToken(token)
	assert.Nil(t, err, "token should valid")
	assert.Equal(t, id, claim.ID, "id not equals")
	time.Sleep(40 * time.Second)
	claim, err = user.ParseToken(token)
	assert.NotNilf(t, err, "toke should expired")
}
