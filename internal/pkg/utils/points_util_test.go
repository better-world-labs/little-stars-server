package utils

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestPointsString(t *testing.T) {
	input := []int{12312854458333, 12, 123, 1234}
	except := []string{"12,312,854,458,333", "12", "123", "1,234"}
	for i := range input {
		require.Equal(t, except[i], PointsString(input[i]))
	}
}
