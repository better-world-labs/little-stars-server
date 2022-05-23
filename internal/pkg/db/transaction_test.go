package db_test

import (
	"aed-api-server/internal/interfaces"
	"aed-api-server/internal/pkg/config"
	"aed-api-server/internal/pkg/db"
	"errors"
	"fmt"
	"github.com/go-xorm/xorm"
	"github.com/stretchr/testify/assert"
	"gitlab.openviewtech.com/openview-pub/gopkg/log"
	"math/rand"
	"testing"
)

var tableName string

type testTable struct {
	Id    int64 `xorm:"id pk autoincr"`
	Count int   `xorm:"count"`
}

func InitDbAndConfig() func() {
	c, err := config.LoadConfig("../../../config-local.yaml")
	if err != nil {
		panic("get config error")
	}
	interfaces.InitConfig(c)
	db.InitEngine(c.Database)
	log.Init(c.Log)
	n := rand.Intn(10000)
	tableName = fmt.Sprintf("test_%v", n)

	_, err = db.Exec(`
		CREATE TABLE ` + tableName + `  (
		  id int(0) UNSIGNED NOT NULL AUTO_INCREMENT,
		  count int(0) UNSIGNED NULL,
		  PRIMARY KEY (id)
		)
	`)
	if err != nil {
		panic("create test table error:" + err.Error())
	}

	return func() {
		defer db.GetEngine().Close()
		_, err = db.Exec(`drop table if exists ` + tableName)
		if err != nil {
			panic("drop test table error:" + err.Error())
		}
	}
}

func TestTransaction(t *testing.T) {
	t.Cleanup(InitDbAndConfig())
	t.Run("success", func(t *testing.T) {
		id := int64(1)
		err := db.Transaction(func(session *xorm.Session) error {
			_, _ = session.Table(tableName).Insert(testTable{
				Id:    id,
				Count: 10,
			})
			return nil
		})

		assert.Nil(t, err)
		var test testTable
		_, err = db.Table(tableName).Where("id=?", id).Get(&test)
		assert.Nil(t, err)
		assert.Equal(t, 10, test.Count)
	})

	t.Run("rollback", func(t *testing.T) {
		var xErr = errors.New("rollback")
		id := int64(2)
		err := db.Transaction(func(session *xorm.Session) error {
			_, _ = session.Table(tableName).Insert(testTable{
				Id:    id,
				Count: 10,
			})
			return xErr
		})

		assert.Equal(t, err, xErr)

		var test testTable
		existed, err := db.Table(tableName).Where("id=?", id).Get(&test)
		assert.Nil(t, err)
		assert.False(t, existed)
		assert.Equal(t, 0, test.Count)
	})

	t.Run("nest-success", func(t *testing.T) {
		id := int64(3)
		err := db.Transaction(func(session *xorm.Session) error {
			_, _ = session.Table(tableName).Insert(testTable{
				Id:    id,
				Count: 10,
			})

			return db.Transaction(func(session *xorm.Session) error {
				exec, err := session.Exec(`update `+tableName+` set count = count + 1 where id = ?`, id)
				assert.Nil(t, err)
				affected, _ := exec.RowsAffected()
				assert.Equal(t, int64(1), affected)
				return err
			})
		})

		assert.Nil(t, err)
		var test testTable
		_, err = db.Table(tableName).Where("id=?", id).Get(&test)
		assert.Nil(t, err)
		assert.Equal(t, 11, test.Count)
	})

	t.Run("nest-rollback", func(t *testing.T) {
		id := int64(4)
		var xErr = errors.New("rollback")
		err := db.Transaction(func(session *xorm.Session) error {
			_, _ = session.Table(tableName).Insert(testTable{
				Id:    id,
				Count: 10,
			})

			return db.Transaction(func(session *xorm.Session) error {
				exec, err := session.Exec(`update `+tableName+` set count = count + 1 where id = ?`, id)
				assert.Nil(t, err)
				affected, _ := exec.RowsAffected()
				assert.Equal(t, int64(1), affected)
				return xErr
			})
		})

		var test testTable
		existed, err := db.Table(tableName).Where("id=?", id).Get(&test)
		assert.Nil(t, err)
		assert.False(t, existed)
		assert.Equal(t, 0, test.Count)
	})
}

//
