package auth

import (
	"context"
	"net/http"
	"strings"
)

type contextKey string

const (
	UserIDKey contextKey = "user_id"
	RoleKey   contextKey = "role"
	RawJWTKey contextKey = "raw_jwt"
)

func JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Отсутствует заголовок Authorization", http.StatusUnauthorized)
			return
		}

		parts := strings.Split(authHeader, "Bearer ")
		if len(parts) != 2 {
			http.Error(w, "Неверный формат токена", http.StatusUnauthorized)
			return
		}

		tokenStr := strings.TrimSpace(parts[1])
		claims, err := ParseJWT(tokenStr)
		if err != nil {
			http.Error(w, "Неверный токен: "+err.Error(), http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
		ctx = context.WithValue(ctx, RoleKey, claims.Role)
		ctx = context.WithValue(ctx, RawJWTKey, tokenStr)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetUserID(r *http.Request) string {
	if val := r.Context().Value(UserIDKey); val != nil {
		if id, ok := val.(string); ok {
			return id
		}
	}
	return ""
}

func GetUserRole(r *http.Request) string {
	if val := r.Context().Value(RoleKey); val != nil {
		if role, ok := val.(string); ok {
			return role
		}
	}
	return ""
}

func GetRawJWT(r *http.Request) string {
	if val := r.Context().Value(RawJWTKey); val != nil {
		if jwt, ok := val.(string); ok {
			return jwt
		}
	}
	return ""
}
