package asserts

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGetResource(t *testing.T) {
	err := LoadResourceDir("../../../assert")
	require.Nil(t, err)

	resource, exist := GetResource("evidence_background.jpg")
	require.True(t, exist)
	fmt.Println(len(resource))
}
