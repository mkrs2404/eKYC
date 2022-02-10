package auth

import (
	"errors"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go/v4"
	"github.com/mkrs2404/eKYC/app/messages"
)

//GenerateToken generates a JWT key with the user Id in the payload
func GenerateToken(clientId uint) (string, error) {

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
		Issuer: strconv.Itoa(int(clientId)),
		ExpiresAt: &jwt.Time{
			Time: time.Now().Add(time.Hour * time.Duration(expiryDelay)),
		},
	})

	token, err := claims.SignedString([]byte(secretKey))

	return token, err
}

//ValidateToken takes the tokenString from the request and validates if it is authentic. If yes, client object is returned, else, 'unauthorized' error is returned
func ValidateToken(tokenString string) (int, error) {
	secretKey := os.Getenv("SECRET_KEY")
	if len(secretKey) == 0 {
		log.Fatal("Secret Key not found")
	}

	token, err := jwt.ParseWithClaims(tokenString, &jwt.StandardClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	if err != nil {
		return -1, errors.New(messages.UNAUTHORIZED)
	}

	claims := token.Claims.(*jwt.StandardClaims)
	clientId, _ := strconv.Atoi(claims.Issuer)
	return clientId, err
}
