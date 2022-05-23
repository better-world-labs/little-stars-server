package cache

import (
	"github.com/google/uuid"
)

type distributeLock struct {
	name      string
	value     string
	expiresIn int64
}

const unlockLua = `if cache.call("get",KEYS[1]) == ARGV[1] then
    return cache.call("del",KEYS[1])
else
    return 0
end`

type Lock interface {
	Release() error
}

func GetDistributeLock(name string, expiresIn int64) (bool, Lock, error) {
	conn := GetConn()
	defer conn.Close()

	v := uuid.NewString()
	reply, err := conn.Do("SET", name, v, "NX", "PX", expiresIn)
	if err != nil {
		return false, nil, err
	}

	return reply == "OK", &distributeLock{
		name:      name,
		value:     v,
		expiresIn: expiresIn,
	}, nil
}

func (r *distributeLock) Release() error {
	conn := GetConn()
	defer conn.Close()

	_, err := conn.Do("EVAL", unlockLua, 1, r.name, r.value)
	if err != nil {
		return err
	}

	return nil
}
