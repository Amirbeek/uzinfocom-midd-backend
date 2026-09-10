package middleware

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/Amirbeek/uzinfocom-midd-backend/internal/auth"
	"github.com/Amirbeek/uzinfocom-midd-backend/internal/utils"
	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const claimsKey contextKey = "claims"

func RequireAuth(a *auth.JWTAuthenticator) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization")
			if header == "" {
				utils.UnauthorizedError(w, r, errors.New("authorization header is missing"))
				return
			}

			parts := strings.Split(header, " ")
			if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
				utils.UnauthorizedError(w, r, errors.New("authorization header is malformed"))
				return
			}
			token, err := a.ValidateToken(parts[1])
			if err != nil {
				utils.UnauthorizedError(w, r, err)
				return
			}

			ctx := context.WithValue(r.Context(), claimsKey, token.Claims)
			next.ServeHTTP(w, r.WithContext(ctx)) // Contextga solamiz va handlerlar buni o'qiy oladi

		})
	}
}

// Yordamchi Helper function handlerlar claimni olishi uchun  funksiya
func ClaimsFromContext(ctx context.Context) (jwt.Claims, bool) {
	claims, ok := ctx.Value(claimsKey).(jwt.Claims)
	return claims, ok
}
