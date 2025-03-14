package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"github.com/nathakusuma/elevateu-backend/domain/contract"
	"github.com/nathakusuma/elevateu-backend/domain/ctxkey"
	"github.com/nathakusuma/elevateu-backend/domain/dto"
	"github.com/nathakusuma/elevateu-backend/domain/enum"
	"github.com/nathakusuma/elevateu-backend/domain/errorpkg"
	"github.com/nathakusuma/elevateu-backend/internal/middleware"
	"github.com/nathakusuma/elevateu-backend/pkg/validator"
)

type courseProgressHandler struct {
	val validator.IValidator
	svc contract.ICourseProgressService
}

func InitCourseProgressHandler(
	router fiber.Router,
	midw *middleware.Middleware,
	validator validator.IValidator,
	progressSvc contract.ICourseProgressService,
) {
	handler := courseProgressHandler{
		svc: progressSvc,
		val: validator,
	}

	courseProgressGroup := router.Group("/courses/contents")
	courseProgressGroup.Use(midw.RequireAuthenticated)
	courseProgressGroup.Use(midw.RequireOneOfRoles(enum.UserRoleStudent))

	courseProgressGroup.Post("/videos/:videoId/progresses", handler.updateVideoProgress)
	courseProgressGroup.Post("/materials/:materialId/progresses", handler.completeMaterial)
}

func (h *courseProgressHandler) updateVideoProgress(ctx *fiber.Ctx) error {
	videoID, err := uuid.Parse(ctx.Params("videoId"))
	if err != nil {
		return errorpkg.ErrValidation().WithDetail("Invalid video ID")
	}

	userID, ok := ctx.Locals(ctxkey.UserID).(uuid.UUID)
	if !ok {
		return errorpkg.ErrInvalidBearerToken()
	}

	var req dto.UpdateCourseVideoProgressRequest
	if err = ctx.BodyParser(&req); err != nil {
		return errorpkg.ErrFailParseRequest()
	}

	if err = h.val.ValidateStruct(req); err != nil {
		return err
	}

	if err = h.svc.UpdateVideoProgress(ctx.Context(), userID, videoID, req); err != nil {
		return err
	}

	return ctx.SendStatus(fiber.StatusNoContent)
}

func (h *courseProgressHandler) completeMaterial(ctx *fiber.Ctx) error {
	materialID, err := uuid.Parse(ctx.Params("materialId"))
	if err != nil {
		return errorpkg.ErrValidation().WithDetail("Invalid material ID")
	}

	userID, ok := ctx.Locals(ctxkey.UserID).(uuid.UUID)
	if !ok {
		return errorpkg.ErrInvalidBearerToken()
	}

	if err = h.svc.UpdateMaterialProgress(ctx.Context(), userID, materialID); err != nil {
		return err
	}

	return ctx.SendStatus(fiber.StatusNoContent)
}
