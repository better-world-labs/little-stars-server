package test

import (
	"aed-api-server/internal/pkg/did"
	"github.com/stretchr/testify/require"
	openapi "gitlab.openviewtech.com/openview-pub/gopkg/open-api"
	"testing"
)

var r = did.NewRest(&openapi.Config{
	"http://47.97.231.172:8081",
	"f7ec97332f0e4016a64961e77bbc8f36",
	"PN5fQILpp6qdejOE0SrqvToZ2C",
})

func TestLogin(t *testing.T) {
	err := r.Login("99911111131")
	require.Nil(t, err)
}
