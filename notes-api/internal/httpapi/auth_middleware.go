package httpapi

import (
	"context"
	"net/http"
	"strings"

	"notes-api/internal/service"
)

const userIDKey contextKey = "userID"

func AuthMiddleware(authService *service.AuthService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")

			parts := strings.SplitN(authHeader, " ", 2)

			if len(parts) != 2 || parts[0] != "Bearer" {
				http.Error(w, "Missing or invalid token", http.StatusUnauthorized)
				return
			}

			userID, err := authService.ValidateToken(parts[1])
			if err != nil {
				http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), userIDKey, userID)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func authenticatedUserID(w http.ResponseWriter, r *http.Request) (int64, bool) {
	userID, ok := r.Context().Value(userIDKey).(int64)
	if !ok || userID <= 0 {
		http.Error(w, "Missing or invalid token", http.StatusUnauthorized)
		return 0, false
	}

	return userID, true
}

func GetUserID(r *http.Request) (int64, bool) {
	userID, ok := r.Context().Value(userIDKey).(int64)
	return userID, ok
}
