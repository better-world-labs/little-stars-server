package base

import (
	"log"
	"testing"
)

func TestError(t *testing.T) {
	e1 := NewError("Module1", "error text 1")
	log.Printf("%v", e1)

	e2 := WrapError("Module2", "error text 2", e1)
	log.Printf("%v", e2)

	e3 := WrapError("Module3", "error text 3", e2)
	log.Printf("%v", e3)
}
