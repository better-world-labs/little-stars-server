package db

import (
	"github.com/go-xorm/xorm"
	"github.com/jtolds/gls"
	log "github.com/sirupsen/logrus"
	"sync"
)

//============================================================================
func sessionRollback(session *xorm.Session) {
	err := session.Rollback()
	if err != nil {
		log.Errorf("session rollback err:%v", err)
		panic(err)
	}
}

var sessionMap = sync.Map{}

//var sessionMap = make(map[uint]*xorm.Session) fatal error: concurrent map writes

func getTransaction(id uint) (*xorm.Session, bool) {
	session, suc := sessionMap.Load(id)
	if suc {
		return session.(*xorm.Session), false
	}
	session = engine.NewSession()
	sessionMap.Store(id, session)
	return session.(*xorm.Session), true
}

func delTransaction(id uint, session *xorm.Session) {
	session.Close()
	sessionMap.Delete(id)
}

//Transaction 事物处理 不允许在事物中新开协程，否则事物会失效
func Transaction(fn func(session *xorm.Session) error) error {
	var err error
	gls.EnsureGoroutineId(func(gid uint) {
		session, isNew := getTransaction(gid)

		if isNew {
			defer delTransaction(gid, session)
			defer func() {
				if info := recover(); info != nil {
					log.Errorf("transaction has panic:%v, transaction will rollback", info)
					sessionRollback(session)
				}
			}()

			err = session.Begin()
			if err != nil {
				return
			}
		}

		err = fn(session)
		if err == nil && isNew {
			err = session.Commit()
		}

		if err != nil {
			if isNew {
				sessionRollback(session)
			}
		}
	})
	return err
}
