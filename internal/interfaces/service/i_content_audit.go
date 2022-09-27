package service

import (
	"aed-api-server/internal/interfaces/entities"
)

type IContentAudit interface {

	// ScanText 文本内容同步审核，审核通过返回 true
	ScanText(text string) (*entities.AuditResult, error)

	// ScanImage 图片内容同步审核，审核通过返回 true
	ScanImage(imgUrl string) (*entities.AuditResult, error)
}
