package auth

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go/v4"
)

//GenerateToken generates a JWT key with the user Id in the payload
func GenerateToken(userId int) (string, error) {

	//Fetching the JWT token delay from env variables
	tokenExpiryDelay := os.Getenv("TOKEN_EXPIRY_DELAY")
	if len(tokenExpiryDelay) == 0 {
		log.Fatal("Token expiry delay not found")
	}

	expiryDelay, err := strconv.Atoi(tokenExpiryDelay)
	if err != nil {
		log.Fatal("Incorrect delay provided")
	}

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
