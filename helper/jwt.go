package helper

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
)

func GenerateJWT(signKey string, userID uint, name string) string {
	var accessToken = generateToken(signKey, userID, name)
	if accessToken == "" {
		return ""
	}
	return accessToken
}


func generateToken(signKey string, id uint, name string) string {
	var claims = jwt.MapClaims{}
	claims["id"] = id
	claims["name"] = name
	claims["iat"] = time.Now().Unix()
	claims["exp"] = time.Now().Add(time.Minute * 10).Unix()

	var sign = jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	validToken, err := sign.SignedString([]byte(signKey))

	if err != nil {
		return ""
	}

	return validToken
}


func ExtractToken(token *jwt.Token) (uint, string) {
	if token.Valid {
		var claims = token.Claims
		expTime, _ := claims.GetExpirationTime()
		fmt.Println(expTime.Time.Compare(time.Now()))
		if expTime.Time.Compare(time.Now()) > 0 {
			var mapClaim = claims.(jwt.MapClaims)
			return uint(mapClaim["id"].(float64)), mapClaim["name"].(string)
		}

		logrus.Error("Token expired")
	}
	return 0, ""
}

