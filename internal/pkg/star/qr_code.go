package star

import (
	"aed-api-server/internal/interfaces"
	"fmt"
)

func GetPlaceCardSharedQrCodeContent(sharer int64) string {
	host := interfaces.GetConfig().Server.Host
	return fmt.Sprintf("https://%s/share/cert?source=placard&sharer=%d", host, sharer)
}

func GetEssaySharedQrCodeContent(sharer int64, url, source string) string {
	host := interfaces.GetConfig().Server.Host
	return fmt.Sprintf("https://%s/share/cert?source=%s&url=%s&sharer=%d", host, source, url, sharer)
}
