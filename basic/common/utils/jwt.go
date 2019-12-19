package utils

import (
	"dio/basic/config"
	"github.com/dgrijalva/jwt-go"
	"time"
)

type Claims struct {
	UserId   int32  `json:"id"`
	Username string `json:"username"`
	ChName   string `json:"ch_name"`
	Admin    bool   `json:"admin"`
	jwt.StandardClaims
}

func GenerateToken(userId int32, username, chName string, admin bool) (string, error) {
	nowTime := time.Now()
	// 过期时间，默认1小时
	expireTime := nowTime.Add(60 * time.Minute)
	claims := Claims{
		userId,
		username,
		chName,
		admin,
		jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			Issuer:    "dio",
		},
	}
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString(config.GetJwtConfig().GetSecret())
	return token, err
}

// ParseToken parsing token
func ParseToken(token string) (*Claims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return config.GetJwtConfig().GetSecret(), nil
	})
	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}
	return nil, err
}
