package services

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go/v4"
)

//This method generates a JWT key with the user Id in the payload
func GenerateToken(userId int, expiryDelay int) (string, error) {

	secretKey := os.Getenv("SECRET_KEY")
	if len(secretKey) == 0 {
		log.Fatal("Secret Key not found")
	}
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer: strconv.Itoa(userId),
		ExpiresAt: &jwt.Time{
			Time: time.Now().Add(time.Hour * time.Duration(expiryDelay)),
		},
	})

	token, err := claims.SignedString([]byte(secretKey))

	return token, err
}
