package middleware

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/roystondz/students-api/internal/utils/response"
)

type contextKey string

const (
	UserIDKey   contextKey = "userId"
	UsernameKey contextKey = "username"
)

func AuthMiddleware(jwtKey []byte) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var tokenString string

			// 1. Try to read token from cookies
			cookie, err := r.Cookie("token")
			if err == nil {
				tokenString = cookie.Value
			} else {
				// 2. Fallback: check Authorization header (Bearer <token>)
				authHeader := r.Header.Get("Authorization")
				if strings.HasPrefix(authHeader, "Bearer ") {
					tokenString = strings.TrimPrefix(authHeader, "Bearer ")
				}
			}

			if tokenString == "" {
				response.WriteJson(w, http.StatusUnauthorized, response.CommonError(errors.New("unauthorized: missing token")))
				return
			}

			// 3. Parse and validate the JWT token
			claims := jwt.MapClaims{}
			token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
				return jwtKey, nil
			})

			if err != nil || !token.Valid {
				response.WriteJson(w, http.StatusUnauthorized, response.CommonError(errors.New("unauthorized: invalid or expired token")))
				return
			}

			// 4. Inject claims into request context for down-stream handlers
			ctx := r.Context()
			if id, ok := claims["id"].(float64); ok {
				ctx = context.WithValue(ctx, UserIDKey, int64(id))
			}
			if username, ok := claims["username"].(string); ok {
				ctx = context.WithValue(ctx, UsernameKey, username)
			}

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
