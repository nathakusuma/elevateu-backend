package middleware

import (
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"github.com/nathakusuma/elevateu-backend/domain/ctxkey"
	"github.com/nathakusuma/elevateu-backend/domain/enum"
	"github.com/nathakusuma/elevateu-backend/domain/errorpkg"
	"github.com/nathakusuma/elevateu-backend/pkg/jwt"
)

func (m *Middleware) RequireAuthenticated(ctx *fiber.Ctx) error {
	header := ctx.Get("Authorization")
	if header == "" {
		return errorpkg.ErrNoBearerToken
	}

	headerSlice := strings.Split(header, " ")
	if len(headerSlice) != 2 && headerSlice[0] != "Bearer" {
		return errorpkg.ErrInvalidBearerToken
	}

	token := headerSlice[1]
	var claims jwt.Claims
	err := m.jwt.Decode(token, &claims)
	if err != nil {
		return errorpkg.ErrInvalidBearerToken
	}

	expirationTime, err := claims.GetExpirationTime()
	if err != nil {
		return errorpkg.ErrInvalidBearerToken
	}

	if expirationTime.Before(time.Now()) {
		return errorpkg.ErrInvalidBearerToken
	}

	ctx.Locals(ctxkey.UserID, uuid.MustParse(claims.Subject))
	ctx.Locals(ctxkey.UserRole, claims.Role)

	return ctx.Next()
}

// RequireOneOfRoles dependency: RequireAuthenticated
func (m *Middleware) RequireOneOfRoles(roles ...enum.UserRole) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		userRole, ok := ctx.Locals(ctxkey.UserRole).(enum.UserRole)
		if !ok {
			return errorpkg.ErrInvalidBearerToken
		}

		for _, role := range roles {
			if userRole == role {
				return ctx.Next()
			}
		}

		return errorpkg.ErrForbiddenRole
	}
}
