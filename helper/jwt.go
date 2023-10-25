package helper

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
)

func GenerateJWT(signKey string, refreshKey string, userID string) map[string]any {
	var result = map[string]any{}
	var accessToken = generateToken(signKey, userID)
	if accessToken == "" {
		return nil
	}
	result["access_token"] = accessToken
	return result
}


func generateToken(signKey string, id string) string {
	var claims = jwt.MapClaims{}
	claims["id"] = id
	claims["iat"] = time.Now().Unix()
	claims["exp"] = time.Now().Add(time.Minute * 10).Unix()

	var sign = jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	validToken, err := sign.SignedString([]byte(signKey))

	if err != nil {
		return ""
	}

	return validToken
}


func ExtractToken(token *jwt.Token) any {
	if token.Valid {
		var claims = token.Claims
		expTime, _ := claims.GetExpirationTime()
		fmt.Println(expTime.Time.Compare(time.Now()))
		if expTime.Time.Compare(time.Now()) > 0 {
			var mapClaim = claims.(jwt.MapClaims)
			return mapClaim["id"]
		}

		logrus.Error("Token expired")
		return nil

	}
	return nil
}
