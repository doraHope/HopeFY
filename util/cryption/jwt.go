package cryption

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"

	"github.com/doraHope/HopeFY/settting"
)

var jwtSecret = []byte(settting.AppSetting.JwSecret)

type Claims struct {
	UserName string                 `json:"username"`
	UserId   uint64                 `json:"userId"`
	Extra    map[string]interface{} `json:"extra, omitempty"`
	jwt.StandardClaims
}

func GenerateToken(userName string, userId uint64, extra map[string]interface{}) (string, error) {
	timeNow := time.Now()
	timeExpired := timeNow.Add(2 * time.Hour)
	claims := Claims{
		userName,
		userId,
		extra,
		jwt.StandardClaims{
			ExpiresAt: timeExpired.Unix(),
			Issuer:    "sWithY",
		},
	}
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString(jwtSecret)
	return token, err
}

func ParseToken(token string) (*Claims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil {
		//todo log
		return nil, err
	}
	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}
	return nil, fmt.Errorf("parse Token but not thing") //解析token成功, 但是空无一物
}
