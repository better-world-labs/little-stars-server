package achievement

import (
	"aed-api-server/internal/pkg/asserts"
	"os"
	"testing"
)

func TestDrawMedalShare(t *testing.T) {
	os.Chdir("../../..")
	err := asserts.LoadResourceDir("assert")
	if err != nil {
		panic(err)
	}

	asserts.GetResource("medal_1_share_background.jpg")
}
