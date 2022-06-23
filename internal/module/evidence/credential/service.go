package credential

import (
	"aed-api-server/internal/interfaces/entities"
	"aed-api-server/internal/pkg/base"
	log "github.com/sirupsen/logrus"
	openapi "gitlab.openviewtech.com/openview-pub/gopkg/open-api"
	"go.uber.org/zap/buffer"
)

type service struct {
	config openapi.Config
}

func NewService(config openapi.Config) Service {
	return &service{config: config}
}

func (s *service) CreateCredential(claim entities.IClaim) (*Info, error) {
	log.Infof("CreateCredential: cptId=%d, claims=%v", claim.CptID(), claim)

	var resp = map[string]interface{}{}
	err := openapi.RestRequest(s.config, "/open/api/credential/create_evi", map[string]interface{}{
		"claim":  claim,
		"cpt_id": claim.CptID(),
		"expire": 0,
	}, &resp)
	if err != nil {
		return nil, err
	}

	code := int(resp["code"].(float64))
	if code != 0 {
		return nil, base.WrapError("credential.service", resp["msg"].(string), err)
	}

	data := resp["data"].(map[string]interface{})

	info, err := s.getCredentialInfo(data["vc_id"].(string))
	if err != nil {
		return nil, err
	}

	return info, nil
	//now := time.Now().UnixMilli()
	//openapi.SignReq(s.config.AppKey, s.config.AppSecret, now,)
	//
	//resp, err := req.Post(fmt.Sprintf("%s/%s", s.config.Server, "/open/api/credential/create_evi"), claims)
	//if err != nil {
	//	return "", base.WrapError("credential.service", "do request error", err)
	//}
	//
	//log.DefaultLogger().Infof("response: %s", resp.String())
	//
	//if resp.Response().StatusCode != 200 {
	//	log.DefaultLogger().Errorf("CreateCredential with error code %d", resp.Response().StatusCode)
	//	return "", base.NewError("credential.service", "do request error")
	//}
	//
	//var res = map[string]interfaces{}{}
	//err = resp.ToJSON(&res)
	//if err != nil {
	//	return "", base.WrapError("credential.service", "to json error", err)
	//}
	//
	//code := int(res["code"].(float64))
	//if code != 0 {
	//	return "", base.WrapError("credential.service", res["message"].(string), err)
	//}
	//
	//data := res["data"].(map[string]interfaces{})
	//
	//d := data["vc_id"].(string)
	//return d, nil
}

func (s *service) getCredentialInfo(id string) (*Info, error) {
	b := buffer.Buffer{}
	b.AppendString(`{"vc_id": "`)
	b.AppendString(id)
	b.AppendString(`"}`)

	var resp = map[string]interface{}{}
	err := openapi.RestRequest(s.config, "/open/api/credential/info", b.Bytes(), &resp)
	if err != nil {
		return nil, err
	}

	code := int(resp["code"].(float64))
	if code != 0 {
		return nil, base.WrapError("credential.service", resp["msg"].(string), err)
	}

	data := resp["data"].(map[string]interface{})
	fileUrl, exists := data["evidence_file_url"]

	if !exists {
		return nil, base.NewError("credential.service", "no file url found")
	}
	date := data["issuance_date"].(float64)
	return &Info{
		Claim:           data["claim"].(map[string]interface{}),
		VcId:            id,
		IssuanceDate:    int64(date),
		SignatureValue:  data["signature_value"].(string),
		EvidenceFileURL: fileUrl.(string),
	}, nil
}
