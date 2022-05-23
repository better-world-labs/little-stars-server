package utils

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestDistinctInt(t *testing.T) {
	distincted := DistinctInt([]int64{2, 2, 1, 4, 5, 6, 4, 4, 3, 4, 56})
	require.Equal(t, []int64{2, 1, 4, 5, 6, 3, 56}, distincted)
}
