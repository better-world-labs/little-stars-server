package db

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestParamPlaceHolder(t *testing.T) {
	holder := ParamPlaceHolder(3)
	require.Equal(t, "(?,?,?)", holder)
}
