package entities

import "aed-api-server/internal/pkg/global"

type UserCertEntity struct {
	Id          int
	AccountId   int64
	ProjectId   int
	ProjectName string
	Uid         string
	Img         map[string]interface{}
	Created     global.FormattedTime `json:"time,omitempty"`
}

type DtoCert struct {
	Uid         string               `json:"certId,omitempty"`
	ProjectId   int                  `json:"projectId"`
	ProjectName string               `json:"projectName"`
	Origin      string               `json:"origin"`
	Thumbnail   string               `json:"thumbnail"`
	Created     global.FormattedTime `json:"time,omitempty"`
}

type UserCert struct {
	Id          int
	AccountId   int64
	ProjectId   int
	ProjectName string
	Uid         string
	Img         string
	Created     global.FormattedTime `json:"time,omitempty"`
}
