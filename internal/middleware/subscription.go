package middleware

import (
	"github.com/gofiber/fiber/v2"

	"github.com/nathakusuma/elevateu-backend/domain/ctxkey"
	"github.com/nathakusuma/elevateu-backend/domain/enum"
	"github.com/nathakusuma/elevateu-backend/domain/errorpkg"
)

func RequireSubscription(subscriptionType enum.PaymentType) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		userRole, ok := ctx.Locals(ctxkey.UserRole).(enum.UserRole)
		if !ok {
			return errorpkg.ErrInvalidBearerToken
		}

		switch subscriptionType {
		case enum.PaymentTypeBoost:
			if userRole == enum.UserRoleAdmin {
				return ctx.Next()
			}
			isSubscribedBoost, ok := ctx.Locals(ctxkey.IsSubscribedBoost).(bool)
			if !ok || !isSubscribedBoost {
				return errorpkg.ErrNotSubscribed.Build().
					WithDetail("You need to subscribe to Skill Boost to access this.")
			}
		case enum.PaymentTypeChallenge:
			if userRole == enum.UserRoleAdmin || userRole == enum.UserRoleMentor {
				return ctx.Next()
			}
			isSubscribedChallenge, ok := ctx.Locals(ctxkey.IsSubscribedChallenge).(bool)
			if !ok || !isSubscribedChallenge {
				return errorpkg.ErrNotSubscribed.Build().
					WithDetail("You need to subscribe to Skill Challenge to access this.")
			}
		}

		return ctx.Next()
	}
}
