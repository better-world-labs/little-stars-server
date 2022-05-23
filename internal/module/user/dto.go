package user

import (
	"aed-api-server/internal/interfaces/entities"
)

type LoginCommand struct {
	MobileCode   string `json:"mobileCode"`
	Code         string `json:"code"`
	EncryptPhone string `json:"encryptedMobile"`
	Iv           string `json:"iv"`
	Nickname     string `json:"nickname"`
	Avatar       string `json:"avatarUrl"`
}

type SimpleLoginCommand struct {
	Code string `json:"code" binding:"required"`
}

type AccountDTOWithSessionKey struct {
	entities.UserDTO

	SessionKey string `json:"sessionKey"`
}

type MobileDTO struct {
	Mobile string `json:"mobile" binding:"required"`
}
