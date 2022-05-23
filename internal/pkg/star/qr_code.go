package star

import (
	"aed-api-server/internal/pkg/config"
	"fmt"
)

var c config.MiniProgramQrcodeConfig

func Init(conf config.MiniProgramQrcodeConfig) {
	c = conf
}

func GetPlaceCardSharedQrCodeContent(sharer int64) string {
	return fmt.Sprintf("%s/share/cert?source=placard&sharer=%d", c.ContentRootPath, sharer)
}

func GetEssaySharedQrCodeContent(sharer int64, url, source string) string {
	return fmt.Sprintf("%s/share/cert?source=%s&url=%s&sharer=%d", c.ContentRootPath, source, url, sharer)
}
