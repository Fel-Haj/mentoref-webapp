package middleware

import (
	"context"
	"fmt"
	"mentoref-webapp/db"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type BoolContextKey bool
type UserMailContextKey string

var (
	AuthContextKey BoolContextKey
	UserContextKey UserMailContextKey
)

const (
	UserMailKey string = "userMail"
)

func JWTAuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("access_token")
		ctx := r.Context()
		if err == nil {
			token, err := jwt.Parse(cookie.Value, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method")
				}
				return []byte(os.Getenv("SECRET_KEY")), nil
			})

			cookie.Expires = time.Now().Add(30 * time.Minute)
			http.SetCookie(w, cookie)

			if err == nil && token.Valid {
				claims, ok := token.Claims.(jwt.MapClaims)
				if !ok {
					http.Error(w, "Invalid token claims", http.StatusUnauthorized)
					return
				} else {
					ctx = context.WithValue(ctx, AuthContextKey, true)
					ctx = context.WithValue(ctx, UserContextKey, claims[UserMailKey])
				}
			}
		}
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

func GenerateToken(user *db.User) (string, error) {
	claims := jwt.MapClaims{
		UserMailKey: user.Email,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(os.Getenv("SECRET_KEY")))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}
