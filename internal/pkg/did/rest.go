package did

import (
	"errors"
	openapi "gitlab.openviewtech.com/openview-pub/gopkg/open-api"
	"go.uber.org/zap/buffer"
	"log"
)

type Rest struct {
	config *openapi.Config
}

func NewRest(config *openapi.Config) Rest {
	return Rest{
		config: config,
	}
}
func (r *Rest) Login(phone string) error {
	log.Printf("did-rest: login phone=%s", phone)
	json := buffer.Buffer{}
	json.AppendString(`{"phone": "`)
	json.AppendString(phone)
	json.AppendString(`"}`)
	res := make(map[string]interface{}, 0)
	err := openapi.RestRequest(*r.config, "/open/api/auth/login", json.String(), &res)
	if err != nil {
		return err
	}

	if res["code"].(float64) != float64(0) {
		log.Printf("login error : response=%v", res)
		return errors.New("login with error code")
	}

	return nil
}
