package test

import (
	"aed-api-server/internal/pkg/domain/emitter"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestUtils(t *testing.T) {
	tick := &TimeTick{}
	structType := emitter.GetStructType(tick)
	require.Equal(t, "aed-api-server/internal/pkg/domain/emitter/test/TimeTick", structType)

	structType2 := emitter.GetStructType(*tick)
	require.Equal(t, "aed-api-server/internal/pkg/domain/emitter/test/TimeTick", structType2)
}
