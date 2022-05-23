package utils

import (
	"errors"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"gitlab.openviewtech.com/openview-pub/gopkg/log"
	"testing"
	"time"
)

func initLog() {
	log.Init(log.LogConfig{
		Level:        logrus.DebugLevel,
		Output:       "stdout",
		ReportCaller: true,
	})
}

func Test_PromiseAll(t *testing.T) {
	initLog()
	t.Run("error", func(t *testing.T) {
		myErr := errors.New("my error")

		all, err := PromiseAll(func() (interface{}, error) {
			time.Sleep(100 * time.Millisecond)
			println("deal 1")
			return nil, errors.New("error 1")
		}, func() (interface{}, error) {
			time.Sleep(20 * time.Millisecond)
			println("deal 2")
			return nil, errors.New("error 2")
		}, func() (interface{}, error) {
			time.Sleep(4 * time.Millisecond)
			println("deal 3")
			return nil, myErr
		}, func() (interface{}, error) {
			time.Sleep(1 * time.Millisecond)
			println("deal 4")
			return "test", nil
		})
		assert.Equal(t, err, myErr)
		assert.Nil(t, all)
	})

	t.Run("panic", func(t *testing.T) {
		all, err := PromiseAll(func() (interface{}, error) {
			time.Sleep(100 * time.Millisecond)
			println("deal 1")
			return nil, errors.New("error 1")
		}, func() (interface{}, error) {
			time.Sleep(20 * time.Millisecond)
			println("deal 2")
			return nil, errors.New("error 2")
		}, func() (interface{}, error) {
			time.Sleep(4 * time.Millisecond)
			println("deal 3")
			panic("make a panic")
			return nil, nil
		}, func() (interface{}, error) {
			time.Sleep(1 * time.Millisecond)
			println("deal 4")
			return "test", nil
		})

		assert.NotNil(t, err)
		assert.Nil(t, all)
	})

	t.Run("success", func(t *testing.T) {
		type Object struct {
			Test  int
			chars string
		}

		object := Object{
			Test:  100,
			chars: "this is a test",
		}

		all, err := PromiseAll(func() (interface{}, error) {
			time.Sleep(100 * time.Millisecond)
			return 100, nil
		}, func() (interface{}, error) {
			time.Sleep(20 * time.Millisecond)
			return "string", nil
		}, func() (interface{}, error) {
			time.Sleep(4 * time.Millisecond)
			return object, nil
		}, func() (interface{}, error) {
			time.Sleep(1 * time.Millisecond)
			return &object, nil
		})

		assert.Nil(t, err)
		assert.Equal(t, 100, all[0])
		assert.Equal(t, "string", all[1])
		assert.Equal(t, object, all[2])
		assert.Equal(t, &object, all[3])
	})
}

func Test_PromiseAllArr(t *testing.T) {
	initLog()
	var fn1 = func() (interface{}, error) {
		time.Sleep(100 * time.Millisecond)
		return 100, nil
	}

	var fn2 = func() (interface{}, error) {
		time.Sleep(300 * time.Millisecond)
		return "string", nil
	}

	arr, err := PromiseAllArr([]PromiseProcessor{fn1, fn2})
	assert.Nil(t, err)
	assert.Equal(t, 100, arr[0])
	assert.Equal(t, "string", arr[1])
}
