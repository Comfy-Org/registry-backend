package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var userID = flag.String("user-id", "", "User ID that is represented by the JWT token")
var expire = flag.Duration("expire", 30*24*time.Hour, "Expiry time of the JWT token")

func main() {
	secret, ok := os.LookupEnv("JWT_SECRET")
	if !ok {
		log.Fatalf("Environment variablel `JWT_SECRET` must be defined.\n")
	}

	flag.Parse()
	if userID == nil || *userID == "" {
		flag.PrintDefaults()
		log.Fatalf("Flag `--user-id` must be set to non-empty string.\n")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Subject:   *userID,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(*expire)),
		NotBefore: jwt.NewNumericDate(time.Now()),
	})

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		log.Fatalf("Fail to create jwt token: %v .\n", err)
	}

	fmt.Println(tokenString)
}
