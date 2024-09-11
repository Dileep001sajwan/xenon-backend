package heleprs

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

type SignedDetails struct {
	Uid string
	jwt.StandardClaims
}

func GenerateAccessToken(uid string) (string, error) {
	ACCESS_TOKEN_SECRET := os.Getenv("ACCESS_TOKEN_SECRET")
	if ACCESS_TOKEN_SECRET == "" {
		return "", errors.New("access token secret is not set")
	}
	accessClaims := &SignedDetails{
		Uid: uid,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
		},
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims).SignedString([]byte(ACCESS_TOKEN_SECRET))
	if err != nil {
		return "", err
	}
	return token, nil
}

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func VerifyPassword(hashedPassword string, providedPassword string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(providedPassword))

	if err != nil {
		return false, "Password is incorrect"
	}
	return true, "Password is correct"
}
