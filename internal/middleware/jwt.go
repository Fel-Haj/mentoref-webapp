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
		claims["accId"] = v.ID
		claims["accType"] = "user"
	case *db.Company:
		claims["accId"] = v.ID
		claims["accType"] = "company"
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

func GetClaims(r *http.Request) (jwt.MapClaims, error) {
	cookie, err := r.Cookie("session")
	if err != nil {
		return nil, err
	}
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

func GetAccountID(r *http.Request) (uint32, error) {
	claims, err := GetClaims(r)
	if err != nil {
		return 0, err
	}
	if claims["accId"] == nil {
		return 0, fmt.Errorf("accId not found in claims")
	}

	id, ok := claims["accId"].(float64)
	if !ok {
		return 0, fmt.Errorf("accId is not float64 type in claims")
	}
	return uint32(id), nil
}

func GetAccountType(r *http.Request) (string, error) {
	claims, err := GetClaims(r)
	if err != nil {
		return "", err
	}
	if claims["accType"] == nil {
		return "", fmt.Errorf("accType not found in claims")
	}
	accType, ok := claims["accType"].(string)
	if !ok {
		return "", fmt.Errorf("invalid accType type in claims")
	}
	return accType, nil
}
