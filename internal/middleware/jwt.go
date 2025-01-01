package middleware

import (
	"fmt"
	"mentoref-webapp/db"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateToken(entity interface{}) string {
	claims := jwt.MapClaims{}

	switch v := entity.(type) {
	case *db.User:
		claims["userId"] = v.ID
	case *db.Company:
		claims["companyId"] = v.ID
	default:
		fmt.Println("invalid entity type")
		return ""
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(os.Getenv("SECRET_KEY")))
	if err != nil {
		fmt.Printf("Error generating token: %v\n", err)
		return ""
	}

	return signedToken
}

func GetClaims(cookie *http.Cookie) (jwt.MapClaims, error) {
	token, err := jwt.Parse(cookie.Value, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return []byte(os.Getenv("SECRET_KEY")), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	return claims, nil
}
