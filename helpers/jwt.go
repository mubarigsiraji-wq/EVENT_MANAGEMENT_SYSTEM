package helpers

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/mubarik/EVENT_MANBAGEMENT_SYSTEM/infra"
	"github.com/mubarik/EVENT_MANBAGEMENT_SYSTEM/models"
)

func GenerateJwt(role models.Role, sub string, ExpireIn int64, isrefreshToken bool) (string, error) {

	config := infra.Configuration

	var jwtsecret []byte

	if isrefreshToken {
		jwtsecret = []byte(config.Refresh_jwt_token)
	} else {
		jwtsecret = []byte(config.Access_jwt_Token)
	}

	claims := jwt.MapClaims{
		"sub":            sub,
		"nbf":            time.Now().Unix(),                                            // ✅ FIXED
		"exp":            time.Now().Add(time.Duration(ExpireIn) * time.Second).Unix(), // ✅ FIXED
		"isrefreshToken": isrefreshToken,
		"role":           role,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(jwtsecret)
}
