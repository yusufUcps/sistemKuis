package helper

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
)

type JWTInterface interface {
	GenerateJWT(userID uint) string
	GenerateToken(id uint) string
	ExtractToken(token *jwt.Token) uint
}

type JWT struct {
	signKey string
}

func New(signKey string) JWTInterface {
	return &JWT{
		signKey: signKey,
	}
}

func (j *JWT) GenerateJWT(userID uint) string {
	var accessToken = j.GenerateToken(userID)
	if accessToken == "" {
		return ""
	}
	return accessToken
}

func (j *JWT) GenerateToken(id uint) string {
	var claims = jwt.MapClaims{}
	claims["id"] = id
	claims["iat"] = time.Now().Unix()
	claims["exp"] = time.Now().Add(time.Hour* 1).Unix()

	var sign = jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	validToken, err := sign.SignedString([]byte(j.signKey))

	if err != nil {
		return ""
	}

	return validToken
}

func (j *JWT) ExtractToken(token *jwt.Token) uint {
    if token != nil {
        claims, ok := token.Claims.(jwt.MapClaims)
        if !ok {
            logrus.Error("Invalid token claims")
            return 0
        }

        id, exists := claims["id"].(float64)
        if !exists {
            logrus.Error("ID not found in token claims")
            return 0
        }

        return uint(id)
    }

    logrus.Error("Invalid token")
    return 0
}

