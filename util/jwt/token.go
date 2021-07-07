package auth

import (
	"fmt"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

var jwtKey = []byte("my_secret_key")

type Claims struct {
	Email string `json:"email"`
	jwt.StandardClaims
}

func CreateToken(w http.ResponseWriter, email string) (string, error) {

	expirationTime := time.Now().Add(1 * time.Hour)
	// Create the JWT claims, which includes the user email and expiry time
	claims := &Claims{
		Email: email,
		StandardClaims: jwt.StandardClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: expirationTime.Unix(),
		}}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)

	if err != nil {
		fmt.Println(err)
		return "", err
	}

	return tokenString, nil
}
