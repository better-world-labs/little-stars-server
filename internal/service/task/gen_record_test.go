package task

import (
	"testing"
	"time"
)

func Test_findLastedByUserId(t *testing.T) {
	t.Cleanup(InitDbAndConfig())

	id, err := findLastedByUserId(53, time.Now().Add(-recentTime))
	println("lookup =", id, err)
}
