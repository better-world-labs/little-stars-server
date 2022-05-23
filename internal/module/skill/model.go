package skill

import (
	"aed-api-server/internal/pkg/global"
	"time"
)

type Project struct {
	Id          int    `json:"id,omitempty"`             // "id": "1"
	Name        string `json:"name,omitempty"`           // "name": "AED高级认证",
	Template    string `json:"pdfTemplateURL,omitempty"` // "pdfTemplateURL": "http://xxx/xxx"
	Images      string
	Description string
	Title       string
	Pic         string
	GrayImg     string
}

type Exam struct {
	Id            int
	Title         string
	Option        string
	CorrectAnswer string
	ProjectId     int
	Sort          int
	Type          int
}

type UserProject struct {
	Id         int       // Id
	AccountId  int64     // 用户ID
	ProjectId  int       // 认证项目ID
	ExamAnswer string    // 考试答题答案
	Status     string    // 状态
	CertImg    string    // 通过后认证证书图片
	Points     float64   // 获得积分
	Updated    time.Time `xorm:"updated"`
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

type UserCertEntity struct {
	Id          int
	AccountId   int64
	ProjectId   int
	ProjectName string
	Uid         string
	Img         map[string]interface{}
	Created     global.FormattedTime `json:"time,omitempty"`
}
