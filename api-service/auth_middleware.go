package main

import (
	"context"
	"net/http"
	"strings"
)

type contextKey string

const userContextKey contextKey = "authenticated_user"

func authMiddleware(authService *AuthService, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Yetkisiz erişim: Giriş yapmalısınız", http.StatusUnauthorized)
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			http.Error(w, "Geçersiz kimlik doğrulama başlığı biçimi", http.StatusUnauthorized)
			return
		}

		token := strings.TrimSpace(parts[1])
		user, ok := authService.Authenticate(token)
		if !ok {
			http.Error(w, "Geçersiz veya süresi dolmuş oturum", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), userContextKey, user)
		next(w, r.WithContext(ctx))
	}
}

func getUserFromContext(ctx context.Context) (User, bool) {
	user, ok := ctx.Value(userContextKey).(User)
	return user, ok
}

