package user

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type TokenClaims struct {
	jwt.StandardClaims

	ID int64 `json:"id"`
}

var (
	expiresInSecond time.Duration
	secret          string
)

func InitJwt(s string, e int64) {
	expiresInSecond = time.Duration(e)
	secret = s
}

// SignToken JwtToken 签发
// @param id account id
// @return Jwt Token String & error
func SignToken(id int64) (string, error) {
	claims := &TokenClaims{
		ID: id,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Second * expiresInSecond).UnixMilli(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// ParseToken JwtToken 解析
// @return TokenClaims & error
func ParseToken(tokenStr string) (*TokenClaims, error) {
	if tokenStr == "" {
		return nil, errors.New("empty token")
	}
	token, err := jwt.ParseWithClaims(tokenStr, &TokenClaims{}, func(token *jwt.Token) (i interface{}, err error) {
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*TokenClaims); ok && token.Valid && claims.ExpiresAt > time.Now().UnixMilli() {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
