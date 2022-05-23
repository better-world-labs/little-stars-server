package evidence

import (
	"aed-api-server/internal/interfaces/entities"
	"aed-api-server/internal/pkg/base"
	"encoding/json"
	"fmt"
	openapi "gitlab.openviewtech.com/openview-pub/gopkg/open-api"
	"gitlab.openviewtech.com/openview-pub/gopkg/uuid2"
	"go.uber.org/zap/buffer"
	"io"
)

type Api interface {
	CreateFileEvidence(name string, r io.ReadCloser) (*entities.EvidenceCreateResponse, error)
	CreateTextEvidence(name string, json map[string]interface{}) (*entities.EvidenceCreateResponse, error)
	GetEvidenceInfo(id string) (*entities.EvidenceDetailInfo, error)
}

type api struct {
	config openapi.Config
}

func NewApi(config openapi.Config) Api {
	return &api{
		config: config,
	}
}

func (s *api) CreateFileEvidence(name string, r io.ReadCloser) (*entities.EvidenceCreateResponse, error) {
	uuid2.InitSnowFlake(1)
	id := uuid2.SnowflakeUUID()
	var resp = map[string]interface{}{}
	upload := openapi.FileUpload{
		FileName:  fmt.Sprintf("%s.pdf", id),
		FieldName: "file",
		File:      r,
	}
	err := openapi.RestRequestWithFile(s.config, "/open/api/evidence/create/file", map[string]interface{}{
		"name": name,
		"id":   id,
	}, upload, &resp)

	if err != nil {
		return nil, err
	}

	code := int(resp["code"].(float64))
	if code != 0 {
		return nil, base.WrapError("credential.api", resp["msg"].(string), err)
	}

	data := resp["data"].(map[string]interface{})

	response, err := decodeCreateResponse(data)
	if err != nil {
		return nil, base.WrapError("credential.api", "invalid response", err)
	}

	return response, nil
}

func (s *api) GetEvidenceInfo(id string) (*entities.EvidenceDetailInfo, error) {
	b := buffer.Buffer{}
	b.AppendString(`{"id":"`)
	b.AppendString(id)
	b.AppendString(`"}`)

	var resp = map[string]interface{}{}
	err := openapi.RestRequest(s.config, "/open/api/evidence/info", b.Bytes(), &resp)
	if err != nil {
		return nil, err
	}

	code := int(resp["code"].(float64))
	if code != 0 {
		return nil, base.WrapError("credential.api", resp["msg"].(string), err)
	}

	data := resp["data"].(map[string]interface{})

	info, err := decodeDetailInfo(data)
	if err != nil {
		return nil, base.WrapError("credential.api", "decode error", err)
	}

	return info, nil
}

func decodeCreateResponse(data map[string]interface{}) (*entities.EvidenceCreateResponse, error) {
	return &entities.EvidenceCreateResponse{
		ID:              data["id"].(string),
		TransactionHash: data["transaction_hash"].(string),
		Type:            int(data["type"].(float64)),
	}, nil
}

func decodeDetailInfo(data map[string]interface{}) (*entities.EvidenceDetailInfo, error) {
	response, err := decodeCreateResponse(data)
	if err != nil {
		return nil, err
	}

	return &entities.EvidenceDetailInfo{
		EvidenceCreateResponse: *response,
		Content:                data["content"].(string),
		Created:                int64(data["created"].(float64)),
	}, nil
}

func (s *api) CreateTextEvidence(name string, jsonPayload map[string]interface{}) (*entities.EvidenceCreateResponse, error) {
	uuid2.InitSnowFlake(1)
	var resp = map[string]interface{}{}
	b, err := json.Marshal(&jsonPayload)
	if err != nil {
		return nil, err
	}

	err = openapi.RestRequest(s.config, "/open/api/evidence/create/text", map[string]interface{}{
		"name":    name,
		"content": string(b),
	}, &resp)

	if err != nil {
		return nil, err
	}

	code := int(resp["code"].(float64))
	if code != 0 {
		return nil, base.WrapError("credential.api", resp["msg"].(string), err)
	}

	data := resp["data"].(map[string]interface{})

	response, err := decodeCreateResponse(data)
	if err != nil {
		return nil, base.WrapError("credential.api", "invalid response", err)
	}

	return response, nil

}
