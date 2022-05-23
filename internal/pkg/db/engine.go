package db

import (
	"aed-api-server/internal/pkg/base"
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"gitlab.openviewtech.com/openview-pub/gopkg/log"
	"strings"
	"time"
	"xorm.io/core"
)

var engine *xorm.Engine

func InitEngine(config MysqlConfig) {
	e, err := xorm.NewEngine(config.DriverName, config.Dsn)
	if err != nil {
		panic(err)
	}
	err = e.Ping() // 可以判断是否能连接
	if err != nil {
		panic(err)
	}

	engine = e
	engine.SetConnMaxLifetime(config.MaxLifetime * time.Millisecond)
	engine.SetMaxOpenConns(config.MaxOpen)
	engine.SetMaxIdleConns(config.MaxIdleCount)
	engine.ShowSQL(true)
	engine.Logger().SetLevel(core.LOG_DEBUG)
}

func GetSession() *xorm.Session {
	return engine.NewSession()
}

func GetEngine() *xorm.Engine {
	return engine
}

func GetById(table string, id interface{}, bean interface{}) (bool, error) {
	return Table(table).Where("id = ?", id).Get(bean)
}

func Insert(table string, beans ...interface{}) (int64, error) {
	return engine.Table(table).Insert(beans...)
}

func Table(table string) *xorm.Session {
	return engine.Table(table)
}

func SQL(query interface{}, args ...interface{}) *xorm.Session {
	return engine.SQL(query, args...)
}

func Exec(sqlOrArgs ...interface{}) (sql.Result, error) {
	return engine.Exec(sqlOrArgs...)
}

func Exist(query interface{}, args ...interface{}) (bool, error) {
	type Existed struct {
		Existed int8 `xorm:"existed"`
	}
	var existed = Existed{}
	return engine.SQL(query, args...).Get(&existed)
}

// WithTransaction 事务 API
// 鉴于手动维护事务操作带来的复杂性，提供此事务 API 进行事务方法编排调用，无需关注提交与回滚动作
// @param s 一个*xorm.Session，无需调用者控制 Begin 与 Commit
// @param f 一个匿名函数，调用方在此做事务接口调用编排
// @return err 错误
func WithTransaction(s *xorm.Session, f func() error) (err error) {
	errChan := make(chan interface{}, 1)
	resChan := make(chan interface{}, 1)

	go func() {
		defer func() {
			close(errChan)
			close(resChan)
		}()

		defer func() {
			err := recover()
			if err != nil {
				fmt.Printf("tx error : %v", err)
				errChan <- err
			}

		}()

		err = s.Begin()
		if err != nil {
			errChan <- err
			return
		}

		err = f()

		if err != nil {
			errChan <- err
			return
		}

		resChan <- struct{}{}
	}()

	for {
		select {
		case e := <-errChan:
			switch e.(type) {
			case error:
				err = e.(error)
			default:
				err = errors.New(fmt.Sprintf("tx error %v", e))
			}

			re := rollback(s)
			if re != nil {
				return base.WrapError("db", "tx rollback error", re)
			}

			return

		case res := <-resChan:
			if res == struct {
			}{} {
				err = commit(s)
				return
			}
		}
	}
}

func rollback(s *xorm.Session) error {
	if err := s.Rollback(); err != nil {
		log.DefaultLogger().Errorf("tx error: %v", err)
		return err
	}

	return nil
}

func commit(s *xorm.Session) error {
	if err := s.Commit(); err != nil {
		log.DefaultLogger().Errorf("tx commiet error: %v", err)
		return err
	}

	return nil
}

// Begin 事务 API (V2)
// 鉴于手动维护事务操作带来的复杂性，提供此事务 API 进行事务方法编排调用，无需关注提交与回滚动作
// @param f 一个匿名函数，调用方在此做事务接口调用编排,调用方需要保证编排的方法使用匿名函数中传入的 session
// @return err 错误
func Begin(f func(session *xorm.Session) error) (err error) {
	session := GetSession()
	defer session.Close()
	return WithTransaction(session, func() error {
		return f(session)
	})
}

func ArrayPlaceholder(n int) string {
	p := make([]string, n)
	for i := 0; i < n; i++ {
		p[i] = "?"
	}
	return strings.Join(p, ",")
}
