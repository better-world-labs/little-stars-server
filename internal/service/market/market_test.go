package market

import (
	"fmt"
	"github.com/google/uuid"
	"math/rand"
	"strconv"
	"testing"
	"time"
)

func TestGenerate(t *testing.T) {
	u := uuid.New()
	formated := strconv.FormatUint(uint64(u.ID()), 10)
	l := len(formated)
	if l < 7 {
		for i := 0; i < 10-i; i++ {
			formated += strconv.Itoa(rand.Int())
		}
	}
	now := time.Now().UnixMilli()
	fmt.Printf("%d%s", now, formated[:7])
}
