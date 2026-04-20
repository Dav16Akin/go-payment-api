package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/Dav16Akin/payment-api/internal/utils"
)

type contextKey string

const UserIDKey contextKey = "user_id"

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		authHeader := r.Header.Get("Authorization")

		if authHeader == "" {
			utils.JSONResponse(w, http.StatusUnauthorized, nil, "Missing Authorization header")
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			utils.JSONResponse(w, http.StatusUnauthorized, nil, "Invalid Authorization format")
			return
		}

		tokenString := parts[1]

		claims, err := utils.ValidateJwt(tokenString)
		if err != nil {
			utils.JSONResponse(w, http.StatusUnauthorized, nil, "Invalid or expired token")
			return
		}

		ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
