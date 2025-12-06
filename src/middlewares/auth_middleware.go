package middlewares

import (
	"context"
	"fmt"
	"go_chat/src/types"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

const userContextKey types.ContextJWTClaimKey = "userID"

func Protect(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		secret := []byte(os.Getenv("JWT_SECRET"))

		if authHeader := r.Header.Get("authorization"); strings.HasPrefix(authHeader, "Bearer") {
			tokenString := strings.Split(authHeader, " ")[1]

			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
				return secret, nil
			})

			if err != nil {
				log.Println(err)
				http.Error(w, "Not authorized, token failed!", http.StatusUnauthorized)
				return
			}

			if claims, ok := token.Claims.(jwt.MapClaims); ok {
				ctx := context.WithValue(r.Context(), userContextKey, claims["id"].(string))
				r = r.WithContext(ctx)
				next.ServeHTTP(w, r)
				return
			} else {
				fmt.Println(err)
				http.Error(w, "Not authorized, token failed!", http.StatusUnauthorized)
				return
			}
		}

		http.Error(w, "Not authorized, missing authorization header!", http.StatusUnauthorized)
	})
}
