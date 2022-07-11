package credential

import (
	"aed-api-server/internal/service/evidence/credential/claim"
	"fmt"
	"github.com/stretchr/testify/assert"
	openapi "gitlab.openviewtech.com/openview-pub/gopkg/open-api"
	"testing"
)

//测试开发环境数据库的可用性
func TestForDev(t *testing.T) {
	doTest(t, "http://192.168.164.223", 2000003)
}

//测试生产环境数据库的可用性
func TestForProd(t *testing.T) {
	doTest(t, "https://base.openviewtech.com", 2076869)
}

type TestClaim struct {
	claim.AedCert
	cptId int
}

func (t *TestClaim) CptID() int {
	return t.cptId
}

func (t *TestClaim) SetCptID(id int) {
	t.cptId = id
}

func doTest(t *testing.T, url string, cptId int) {
	config := openapi.Config{
		Url:       url,
		AppSecret: "PN5fQILpp6qdejOE0SrqvToZ2C",
		AppKey:    "f7ec97332f0e4016a64961e77bbc8f36",
	}

	service := NewServiceWith(config)
	cert := TestClaim{
		AedCert: claim.AedCert{
			User:   "test-xx",
			Detail: "是否可用测试",
		},
	}
	cert.SetCptID(cptId)
	rst, err := service.CreateCredential(&cert)
	assert.Nil(t, err)
	fmt.Printf("rst:%v", rst)
}
