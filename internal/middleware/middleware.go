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

const (
	claimsKey contextKey = "claims"
	userIDKey contextKey = "userID"
)

func RequireAuth(a *auth.JWTAuthenticator) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			header := r.Header.Get("Authorization")
			if header == "" {
				utils.UnauthorizedError(w, r, errors.New("authorization header is missing"))
				return
			}

			parts := strings.Fields(header)
			if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
				utils.UnauthorizedError(w, r, errors.New("authorization header is malformed"))
				return
			}

			token, err := a.ValidateToken(parts[1])
			if err != nil {
				utils.UnauthorizedError(w, r, err)
				return
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				utils.UnauthorizedError(w, r, errors.New("claims invalid "))
				return
			}

			sub, ok := claims["sub"].(float64)
			if !ok {
				utils.UnauthorizedError(w, r, errors.New("user id is missing"))
				return
			}

			userID := int64(sub)

			ctx := context.WithValue(r.Context(), userIDKey, userID)
			ctx = context.WithValue(ctx, claimsKey, token.Claims)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func UserIDFromContext(ctx context.Context) (int64, bool) {
	id, ok := ctx.Value(userIDKey).(int64)
	return id, ok
}

func ClaimsFromContext(ctx context.Context) (jwt.Claims, bool) {
	claims, ok := ctx.Value(claimsKey).(jwt.Claims)
	return claims, ok
}
