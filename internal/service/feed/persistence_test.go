package feed

import (
	"aed-api-server/internal/interfaces"
	"aed-api-server/internal/pkg/db"
	"aed-api-server/internal/server/config"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_persistence_findUserMarkedInFeeds(t *testing.T) {
	c, err := config.LoadConfig("../../../")
	if err != nil {
		panic("get config error")
	}
	interfaces.InitConfig(c)
	db.InitEngine(c.Database)

	p := persistence{}

	feeds, err := p.findUserMarkedInFeeds(50, []int64{1, 2, 3, 4})
	assert.Nil(t, err)
	fmt.Printf("feeds:%v", feeds)
}

func Test_persistence_getMyFeedsCount(t *testing.T) {
	c, err := config.LoadConfig("../../../")
	if err != nil {
		panic("get config error")
	}
	interfaces.InitConfig(c)
	db.InitEngine(c.Database)

	p := persistence{}

	count, err := p.getMyFeedsCount(50)
	assert.Nil(t, err)
	assert.True(t, count > 0)
}
