package db

import (
	"fmt"
	"testing"
)

func TestTuple(t *testing.T) {
	tuple := TupleOf([]int64{1, 2, 3}, 4, 5, 6, "7")
	for _, t := range tuple {
		fmt.Printf("%v, ", t)
	}
}
