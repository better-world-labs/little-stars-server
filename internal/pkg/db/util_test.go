package db

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestParamPlaceHolder(t *testing.T) {
	holder := ParamPlaceHolder(3)
	require.Equal(t, "(?,?,?)", holder)
}

func TestTuple(t *testing.T) {
	tuple := TupleOf([]int64{1, 2, 3}, 4, 5, 6, "7")
	for _, t := range tuple {
		fmt.Printf("%v, ", t)
	}
}
