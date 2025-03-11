package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"

	"github.com/nathakusuma/elevateu-backend/domain/ctxkey"
	"github.com/nathakusuma/elevateu-backend/domain/enum"
	"github.com/nathakusuma/elevateu-backend/domain/errorpkg"
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

	validateResp, err := m.jwt.Validate(token)
	if err != nil {
		return err
	}

	ctx.Locals(ctxkey.UserID, validateResp.UserID)
	ctx.Locals(ctxkey.UserRole, validateResp.Role)
	ctx.Locals(ctxkey.IsSubscribedBoost, validateResp.IsSubscribedBoost)
	ctx.Locals(ctxkey.IsSubscribedChallenge, validateResp.IsSubscribedChallenge)

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
