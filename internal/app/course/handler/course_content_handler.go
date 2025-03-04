package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"github.com/nathakusuma/elevateu-backend/domain/contract"
	"github.com/nathakusuma/elevateu-backend/domain/dto"
	"github.com/nathakusuma/elevateu-backend/domain/enum"
	"github.com/nathakusuma/elevateu-backend/domain/errorpkg"
	"github.com/nathakusuma/elevateu-backend/internal/middleware"
	"github.com/nathakusuma/elevateu-backend/pkg/validator"
)

type courseContentHandler struct {
	val validator.IValidator
	svc contract.ICourseContentService
}

func InitCourseContentHandler(
	router fiber.Router,
	midw *middleware.Middleware,
	validator validator.IValidator,
	contentSvc contract.ICourseContentService,
) {
	handler := courseContentHandler{
		svc: contentSvc,
		val: validator,
	}

	coursesGroup := router.Group("/courses")
	coursesGroup.Use(midw.RequireAuthenticated)

	coursesGroup.Patch("/contents/videos/:id",
		midw.RequireOneOfRoles(enum.UserRoleAdmin), handler.updateVideo)
	coursesGroup.Delete("/contents/videos/:id",
		midw.RequireOneOfRoles(enum.UserRoleAdmin), handler.deleteVideo)
	coursesGroup.Get("/contents/videos/:id/upload-url",
		midw.RequireOneOfRoles(enum.UserRoleAdmin), handler.getVideoUploadURLs)

	coursesGroup.Patch("/contents/materials/:id",
		midw.RequireOneOfRoles(enum.UserRoleAdmin), handler.updateMaterial)
	coursesGroup.Delete("/contents/materials/:id",
		midw.RequireOneOfRoles(enum.UserRoleAdmin), handler.deleteMaterial)
	coursesGroup.Get("/contents/materials/:id/upload-url",
		midw.RequireOneOfRoles(enum.UserRoleAdmin), handler.getMaterialUploadURL)

	coursesGroup.Get("/:courseId/contents", handler.getCourseContents)
	coursesGroup.Post("/:courseId/contents/videos",
		midw.RequireOneOfRoles(enum.UserRoleAdmin), handler.createVideo)
	coursesGroup.Post("/:courseId/contents/materials",
		midw.RequireOneOfRoles(enum.UserRoleAdmin), handler.createMaterial)
}

func (h *courseContentHandler) getCourseContents(ctx *fiber.Ctx) error {
	courseID, err := uuid.Parse(ctx.Params("courseId"))
	if err != nil {
		return errorpkg.ErrValidation.Build().WithDetail("invalid course ID")
	}

	contents, err := h.svc.GetCourseContents(ctx.Context(), courseID)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(map[string]interface{}{
		"course_contents": contents,
	})
}

func (h *courseContentHandler) createVideo(ctx *fiber.Ctx) error {
	courseID, err := uuid.Parse(ctx.Params("courseId"))
	if err != nil {
		return errorpkg.ErrValidation.Build().WithDetail("invalid course ID")
	}

	var req dto.CreateCourseVideoRequest
	if err := ctx.BodyParser(&req); err != nil {
		return errorpkg.ErrFailParseRequest
	}

	if err := h.val.ValidateStruct(req); err != nil {
		return err
	}

	resp, err := h.svc.CreateVideo(ctx.Context(), courseID, req)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusCreated).JSON(resp)
}

func (h *courseContentHandler) updateVideo(ctx *fiber.Ctx) error {
	id, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return errorpkg.ErrValidation.Build().WithDetail("invalid video ID")
	}

	var req dto.UpdateCourseVideoRequest
	if err := ctx.BodyParser(&req); err != nil {
		return errorpkg.ErrFailParseRequest
	}

	if err := h.val.ValidateStruct(req); err != nil {
		return err
	}

	if err := h.svc.UpdateVideo(ctx.Context(), id, req); err != nil {
		return err
	}

	return ctx.SendStatus(fiber.StatusNoContent)
}

func (h *courseContentHandler) deleteVideo(ctx *fiber.Ctx) error {
	id, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return errorpkg.ErrValidation.Build().WithDetail("invalid video ID")
	}

	if err := h.svc.DeleteVideo(ctx.Context(), id); err != nil {
		return err
	}

	return ctx.SendStatus(fiber.StatusNoContent)
}

func (h *courseContentHandler) getVideoUploadURLs(ctx *fiber.Ctx) error {
	id, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return errorpkg.ErrValidation.Build().WithDetail("invalid video ID")
	}

	videoURL, thumbnailURL, err := h.svc.GetVideoUploadURLs(ctx.Context(), id)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(map[string]interface{}{
		"video_upload_url":     videoURL,
		"thumbnail_upload_url": thumbnailURL,
	})
}

func (h *courseContentHandler) createMaterial(ctx *fiber.Ctx) error {
	courseID, err := uuid.Parse(ctx.Params("courseId"))
	if err != nil {
		return errorpkg.ErrValidation.Build().WithDetail("invalid course ID")
	}

	var req dto.CreateCourseMaterialRequest
	if err := ctx.BodyParser(&req); err != nil {
		return errorpkg.ErrFailParseRequest
	}

	if err := h.val.ValidateStruct(req); err != nil {
		return err
	}

	resp, err := h.svc.CreateMaterial(ctx.Context(), courseID, req)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusCreated).JSON(resp)
}

func (h *courseContentHandler) updateMaterial(ctx *fiber.Ctx) error {
	id, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return errorpkg.ErrValidation.Build().WithDetail("invalid material ID")
	}

	var req dto.UpdateCourseMaterialRequest
	if err := ctx.BodyParser(&req); err != nil {
		return errorpkg.ErrFailParseRequest
	}

	if err := h.val.ValidateStruct(req); err != nil {
		return err
	}

	if err := h.svc.UpdateCourseMaterial(ctx.Context(), id, req); err != nil {
		return err
	}

	return ctx.SendStatus(fiber.StatusNoContent)
}

func (h *courseContentHandler) deleteMaterial(ctx *fiber.Ctx) error {
	id, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return errorpkg.ErrValidation.Build().WithDetail("invalid material ID")
	}

	if err := h.svc.DeleteCourseMaterial(ctx.Context(), id); err != nil {
		return err
	}

	return ctx.SendStatus(fiber.StatusNoContent)
}

func (h *courseContentHandler) getMaterialUploadURL(ctx *fiber.Ctx) error {
	id, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return errorpkg.ErrValidation.Build().WithDetail("invalid material ID")
	}

	url, err := h.svc.GetMaterialUploadURL(ctx.Context(), id)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(map[string]interface{}{
		"material_upload_url": url,
	})
}
