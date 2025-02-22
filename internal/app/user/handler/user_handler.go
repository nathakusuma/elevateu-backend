package handler

import (
	"errors"
	"github.com/google/uuid"

	"github.com/gofiber/fiber/v2"

	"github.com/nathakusuma/elevateu-backend/domain/contract"
	"github.com/nathakusuma/elevateu-backend/domain/ctxkey"
	"github.com/nathakusuma/elevateu-backend/domain/dto"
	"github.com/nathakusuma/elevateu-backend/domain/errorpkg"
	"github.com/nathakusuma/elevateu-backend/internal/middleware"
	"github.com/nathakusuma/elevateu-backend/pkg/validator"
)

type userHandler struct {
	val validator.IValidator
	svc contract.IUserService
}

func InitUserHandler(
	router fiber.Router,
	midw *middleware.Middleware,
	validator validator.IValidator,
	userSvc contract.IUserService,
) {
	handler := userHandler{
		svc: userSvc,
		val: validator,
	}

	userGroup := router.Group("/users")
	userGroup.Get("/me",
		midw.RequireAuthenticated,
		handler.getUser("me"),
	)
	userGroup.Get("/:id",
		midw.RequireAuthenticated,
		handler.getUser("id"),
	)
	userGroup.Patch("/me",
		midw.RequireAuthenticated,
		handler.updateUser,
	)
}

func (c *userHandler) getUser(param string) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var userID uuid.UUID
		if param == "me" {
			var ok bool
			userID, ok = ctx.Locals(ctxkey.UserID).(uuid.UUID)
			if !ok {
				return errorpkg.ErrInvalidBearerToken
			}
		} else {
			var err error
			userID, err = uuid.Parse(ctx.Params("id"))
			if err != nil {
				return errorpkg.ErrFailParseRequest
			}
		}

		user, err := c.svc.GetUserByID(ctx.Context(), userID)
		if err != nil {
			return err
		}

		resp := dto.UserResponse{}
		if param == "me" {
			resp.PopulateFromEntity(user)
		} else {
			resp.PopulateMinimalFromEntity(user)
		}

		return ctx.Status(fiber.StatusOK).JSON(map[string]interface{}{
			"user": resp,
		})
	}
}

func (c *userHandler) updateUser(ctx *fiber.Ctx) error {
	var req dto.UpdateUserRequest
	if err := ctx.BodyParser(&req); err != nil {
		return errorpkg.ErrFailParseRequest
	}

	avatar, err := ctx.FormFile("avatar")
	if err != nil && !errors.Is(err, fiber.ErrNotFound) {
		return errorpkg.ErrFailParseRequest
	}
	req.Avatar = avatar

	if err := c.val.ValidateStruct(req); err != nil {
		return err
	}

	userID, ok := ctx.Locals(ctxkey.UserID).(uuid.UUID)
	if !ok {
		return errorpkg.ErrInvalidBearerToken
	}

	if err := c.svc.UpdateUser(ctx.Context(), userID, req); err != nil {
		return err
	}

	return ctx.SendStatus(fiber.StatusNoContent)
}
