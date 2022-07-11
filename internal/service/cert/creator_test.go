package cert

import (
	"aed-api-server/internal/pkg/asserts"
	"aed-api-server/internal/service/img"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
	"time"
)

func TestCreatCert(t *testing.T) {
	err := asserts.LoadResourceDir("../../../assert")
	require.Nil(t, err)

	err = img.Init()
	require.Nil(t, err)

	creator, err := NewImageCreator()
	assert.Nil(t, err)

	create, err := os.Create("cert_out.png")
	defer create.Close()
	assert.Nil(t, err)

	err = creator.Create("https://thirdwx.qlogo.cn/mmopen/vi_32/Q0j4TwGTfTLtZvbOWWopA7216libKCVabh9EcLLmh3UWYZAIQ6XMaxibIZpicRPB7lyibJ5d2zLricwS4wuYfgMfPCA/132",
		"Souththth_Tae🍑",
		"\"茫茫人海之中，去挽救下一个倒地昏迷的人吧\"", time.Now(), create)
	assert.Nil(t, err)
}
