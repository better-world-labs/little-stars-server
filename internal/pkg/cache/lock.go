package cache

import (
	"github.com/google/uuid"
	"time"
)

type distributeLock struct {
	name                 string
	value                string
	expiresAt            time.Time
	expiresInMicroSecond int64
	suc                  bool
	watchdog             func(lock distributeLock)
}

var (
	watchdog = func(lock distributeLock) {
		//TODO
	}
)

const unlockLua = `if redis.call("get",KEYS[1]) == ARGV[1] then
    return redis.call("del",KEYS[1])
else
    return 0
end`

type Lock interface {
	Release() error
	Locked() bool
}

func GetDistributeLock(name string, expiresInMicroSecond int64) (Lock, error) {
	conn := GetConn()
	defer conn.Close()
	v := uuid.NewString()

	expiresAt := time.Now().Add(time.Millisecond * time.Duration(expiresInMicroSecond))
	reply, err := conn.Do("SET", name, v, "NX", "PX", expiresInMicroSecond)
	if err != nil {
		return nil, err
	}

	suc := reply == "OK"
	return &distributeLock{
		name:                 name,
		value:                v,
		expiresAt:            expiresAt,
		expiresInMicroSecond: expiresInMicroSecond,
		suc:                  suc,
	}, nil
}

func (r *distributeLock) Release() error {
	if !r.suc {
		return nil
	}
	conn := GetConn()
	defer conn.Close()

	_, err := conn.Do("EVAL", unlockLua, 1, r.name, r.value)
	if err != nil {
		return err
	}

	return nil
}
func (r distributeLock) Locked() bool {
	return r.suc
}
