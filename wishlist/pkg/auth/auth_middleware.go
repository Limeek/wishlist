package auth

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func JWTMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tokenStr := r.Header.Get("Authorization")
			tokenStr = strings.TrimPrefix(tokenStr, "Bearer ")

			if tokenStr == "" {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			token, err := validateAccessToken(tokenStr)

			if err != nil || !token.Valid {
				w.WriteHeader(http.StatusUnauthorized)
				_, err = w.Write([]byte("Token is not valid"))
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				return
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok || !token.Valid {
				w.WriteHeader(http.StatusUnauthorized)
				_, err = w.Write([]byte("Token is not valid"))
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				return
			}

			userId := claims["sub"].(string)

			userUuid, err := uuid.Parse(userId)

			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				_, err = w.Write([]byte("Token is not valid"))
			}

			ctx := context.WithValue(r.Context(), "userId", userUuid)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func validateAccessToken(accessToken string) (*jwt.Token, error) {
	return jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(os.Getenv("JWT_SECRET")), nil
	})
}
