package cert

import (
	"io"
	"time"
)

type ImageCreator interface {

	// Create 创建证书图片
	// @param avatarUrl 用户头像URL
	// @param nickname 用户昵称
	// @param description 证书正文
	// @param t 时间
	// @param writer Writer
	// @return err 错误
	Create(avatarUrl string, nickname string, description string, t time.Time, closer io.Writer) error
}
