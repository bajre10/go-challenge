package middleware

import (
	"context"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"os"
	"shortcut-challenge/api"
	"shortcut-challenge/handlers"
	"slices"
	"strings"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Get Bearer token from header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			api.ThrowRequestError(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		secret := []byte(os.Getenv("JWT_SECRET"))

		claims := &jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return secret, nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		// Add claims to context
		ctx := context.WithValue(r.Context(), "JWT_CLAIMS", claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func RestrictTo(allowedRoles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims := handlers.GetClaims(r)
			role, _ := (*claims)["role"].(string)

			// Check if user's role is allowed
			if slices.Contains(allowedRoles, role) {
				next.ServeHTTP(w, r)
				return
			}

			// Deny access
			api.ThrowRequestError(w, "Forbidden", http.StatusForbidden)
		})
	}
}
