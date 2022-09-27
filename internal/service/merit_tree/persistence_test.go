package merit_tree

import (
	"aed-api-server/internal/pkg/db"
	"aed-api-server/internal/server/config"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_persistence_getLastExpiredChest(t *testing.T) {
	c, err := config.LoadConfig("../../../")
	if err != nil {
		panic("get config error")
	}
	db.InitEngine(c.Database)
	pe := &persistence{}
	gotDto, err := pe.getLastExpiredChest(3141)

	assert.Nil(t, err)
	println(gotDto)
}
