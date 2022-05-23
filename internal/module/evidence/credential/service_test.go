package credential

import (
	"aed-api-server/internal/module/evidence/credential/claim"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"gitlab.openviewtech.com/openview-pub/gopkg/log"
	openapi "gitlab.openviewtech.com/openview-pub/gopkg/open-api"
	"testing"
)

func TestService_CreateCredential(t *testing.T) {
	service := NewService(openapi.Config{
		Url:       "http://dev002:8081",
		AppSecret: "IzpT0o%#]a&]=#/\\TO`;VUl!s!={?LE&",
		AppKey:    "71bf5ff0147942aaa0b9732a107b7ecb",
	})

	id, err := service.CreateCredential(&claim.AedCert{
		User:   "15548720906",
		Detail: "AED 证书",
	})

	require.Nil(t, err)

	fmt.Println(id)
}

func TestCreateCredential(t *testing.T) {
	log.Init(log.LogConfig{
		Level:  logrus.DebugLevel,
		Output: "stdout",
	})
	service := NewService(openapi.Config{
		Url: "http://192.168.164.223",
		//Url: "http://47.97.231.172:8081",

		//Url: "http://localhost:8080",
		AppSecret: "PN5fQILpp6qdejOE0SrqvToZ2C",
		AppKey:    "f7ec97332f0e4016a64961e77bbc8f36",
	})

	id, err := service.CreateCredential(&claim.AedCert{
		User:   "dapeng",
		Detail: "AED 证书-测试",
	})

	require.Nil(t, err)

	fmt.Println(id)
}
