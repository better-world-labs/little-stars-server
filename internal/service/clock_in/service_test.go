package clock_in

import (
	"aed-api-server/internal/interfaces"
	"aed-api-server/internal/pkg/config"
	"aed-api-server/internal/pkg/db"
	"fmt"
	"github.com/stretchr/testify/assert"
	"gitlab.openviewtech.com/openview-pub/gopkg/log"
	"testing"
)

func initDbAndConfig() func() {
	c, err := config.LoadConfig("../../../")
	if err != nil {
		panic("get config error")
	}
	interfaces.InitConfig(c)
	db.InitEngine(c.Database)
	log.Init(c.Log)

	return func() {
		println("close db")
		engine := db.GetEngine()
		if engine != nil {
			err2 := engine.Close()
			fmt.Printf("err:%v\n", err2)
		}
	}
}

func Test_findLastClockInItem(t *testing.T) {
	t.Cleanup(initDbAndConfig())
	in, err := findLastClockIn([]string{
		"0391aec8-100c-4f5b-9e91-383f049fdff7",
		"09a21935-602c-4d7c-b457-6e5bdd59dffc",
		"09a21935-602c-4d7c-b457-e579d6ebc10c",
		"0b36ff7a-34d9-43e5-ae82-8e6657f26464",
		"0da04070-0014-4a42-99ad-5930a1ae7daf",
		"21ffc346-a18d-4b4e-b92d-92fa3a486cf0",
		"3a5454de-122f-430e-9713-4bede118c6ae",
		"4af0cf76-a86f-4d97-8819-6e5bdd59dffc",
		"79258a82-1e7b-4476-a547-e69087e48fe3",
		"96e5615b-9283-4258-9dc5-20bbd9924c66",
		"a472dc08-c618-49c8-9ca6-c252365e1afd",
		"ae7d1210-5ccb-4051-bb21-380f7f4eb213",
	})

	assert.Nil(t, err)

	inits := []int{1, 2, 3}[0:3]
	assert.Len(t, inits, 3)

	//log.Info(in)
	log.DefaultLogger().Infof("map:%v\n", in)
}
