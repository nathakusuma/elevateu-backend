package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/nathakusuma/elevateu-backend/pkg/log"

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

	userGroup.Get("/leaderboards",
		midw.RequireAuthenticated,
		handler.getLeaderboard,
	)
	userGroup.Get("/me",
		midw.RequireAuthenticated,
		handler.getUser("me"),
	)
	userGroup.Patch("/me",
		midw.RequireAuthenticated,
		handler.updateUser,
	)
	userGroup.Delete("/me",
		midw.RequireAuthenticated,
		handler.deleteUser,
	)
	userGroup.Put("/me/avatar",
		midw.RequireAuthenticated,
		handler.updateUserAvatar)
	userGroup.Delete("/me/avatar",
		midw.RequireAuthenticated,
		handler.deleteUserAvatar)
	userGroup.Get("/:id",
		midw.RequireAuthenticated,
		handler.getUser("id"),
	)
}

func (c *userHandler) getUser(param string) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var userID uuid.UUID
		if param == "me" {
			var ok bool
			userID, ok = ctx.Locals(ctxkey.UserID).(uuid.UUID)
			if !ok {
				traceID := log.ErrorWithTraceID(ctx.Context(), nil, "Failed to get user ID from context")
				return errorpkg.ErrInternalServer().WithTraceID(traceID)
			}
		} else {
			var err error
			userID, err = uuid.Parse(ctx.Params("id"))
			if err != nil {
				return errorpkg.ErrFailParseRequest()
			}
		}

		var isMinimal bool
		if param == "me" {
			isMinimal = false
		}

		resp, err := c.svc.GetUserByID(ctx.Context(), userID, isMinimal)
		if err != nil {
			return err
		}

		return ctx.Status(fiber.StatusOK).JSON(map[string]interface{}{
			"user": resp,
		})
	}
}

func (c *userHandler) updateUser(ctx *fiber.Ctx) error {
	var req dto.UpdateUserRequest

	if err := ctx.BodyParser(&req); err != nil {
		return errorpkg.ErrFailParseRequest()
	}

	userID, ok := ctx.Locals(ctxkey.UserID).(uuid.UUID)
	if !ok {
		traceID := log.ErrorWithTraceID(ctx.Context(), nil, "Failed to get user ID from context")
		return errorpkg.ErrInternalServer().WithTraceID(traceID)
	}

	// Validate role-specific data
	if req.Student != nil {
		if err := c.val.ValidateStruct(req.Student); err != nil {
			return err
		}
	}
	if req.Mentor != nil {
		if err := c.val.ValidateStruct(req.Mentor); err != nil {
			return err
		}
	}

	if err := c.val.ValidateStruct(req); err != nil {
		return err
	}

	if err := c.svc.UpdateUser(ctx.Context(), userID, req); err != nil {
		return err
	}

	return ctx.SendStatus(fiber.StatusNoContent)
}

func (c *userHandler) deleteUser(ctx *fiber.Ctx) error {
	userID, ok := ctx.Locals(ctxkey.UserID).(uuid.UUID)
	if !ok {
		traceID := log.ErrorWithTraceID(ctx.Context(), nil, "Failed to get user ID from context")
		return errorpkg.ErrInternalServer().WithTraceID(traceID)
	}

	if err := c.svc.DeleteUser(ctx.Context(), userID); err != nil {
		return err
	}

	return ctx.SendStatus(fiber.StatusNoContent)
}

// Handler specifically for avatar upload (multipart form)
func (c *userHandler) updateUserAvatar(ctx *fiber.Ctx) error {
	file, err := ctx.FormFile("avatar")
	if err != nil {
		return errorpkg.ErrFailParseRequest()
	}

	userID, ok := ctx.Locals(ctxkey.UserID).(uuid.UUID)
	if !ok {
		traceID := log.ErrorWithTraceID(ctx.Context(), nil, "Failed to get user ID from context")
		return errorpkg.ErrInternalServer().WithTraceID(traceID)
	}

	if err2 := c.svc.UpdateUserAvatar(ctx.Context(), userID, file); err2 != nil {
		return err2
	}

	return ctx.SendStatus(fiber.StatusNoContent)
}

func (c *userHandler) deleteUserAvatar(ctx *fiber.Ctx) error {
	userID, ok := ctx.Locals(ctxkey.UserID).(uuid.UUID)
	if !ok {
		traceID := log.ErrorWithTraceID(ctx.Context(), nil, "Failed to get user ID from context")
		return errorpkg.ErrInternalServer().WithTraceID(traceID)
	}

	if err := c.svc.DeleteUserAvatar(ctx.Context(), userID); err != nil {
		return err
	}

	return ctx.SendStatus(fiber.StatusNoContent)
}

func (c *userHandler) getLeaderboard(ctx *fiber.Ctx) error {
	resp, err := c.svc.GetLeaderboard(ctx.Context())
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(map[string]interface{}{
		"users": resp,
	})
}
