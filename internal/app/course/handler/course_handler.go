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

type courseHandler struct {
	val validator.IValidator
	svc contract.ICourseService
}

func InitCourseHandler(
	router fiber.Router,
	midw *middleware.Middleware,
	validator validator.IValidator,
	courseSvc contract.ICourseService,
) {
	handler := courseHandler{
		svc: courseSvc,
		val: validator,
	}

	courseGroup := router.Group("/courses")
	courseGroup.Use(midw.RequireAuthenticated)

	courseGroup.Post("/",
		midw.RequireOneOfRoles(enum.UserRoleAdmin),
		handler.createCourse)
	courseGroup.Get("/", handler.getCourses)
	courseGroup.Get("/:id", handler.getCourseByID)
	courseGroup.Patch("/:id",
		midw.RequireOneOfRoles(enum.UserRoleAdmin),
		handler.updateCourse)
	courseGroup.Delete("/:id",
		midw.RequireOneOfRoles(enum.UserRoleAdmin),
		handler.deleteCourse)
	courseGroup.Get("/:id/preview-video-upload-url",
		midw.RequireOneOfRoles(enum.UserRoleAdmin),
		handler.GetPreviewVideoUploadURL)
}

func (c *courseHandler) createCourse(ctx *fiber.Ctx) error {
	var req dto.CreateCourseRequest
	if err := ctx.BodyParser(&req); err != nil {
		return errorpkg.ErrFailParseRequest
	}

	var err error
	req.TeacherAvatar, err = ctx.FormFile("teacher_avatar")
	if err != nil {
		return errorpkg.ErrFailParseRequest.Build().WithDetail("Fail to parse teacher avatar")
	}
	req.Thumbnail, err = ctx.FormFile("thumbnail")
	if err != nil {
		return errorpkg.ErrFailParseRequest.Build().WithDetail("Fail to parse thumbnail")
	}

	if err = c.val.ValidateStruct(req); err != nil {
		return err
	}

	resp, err := c.svc.CreateCourse(ctx.Context(), &req)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusCreated).JSON(resp)
}

func (c *courseHandler) getCourses(ctx *fiber.Ctx) error {
	var query dto.GetCoursesQuery
	if err := ctx.QueryParser(&query); err != nil {
		return errorpkg.ErrFailParseRequest
	}

	var pageReq dto.PaginationRequest
	if err := ctx.QueryParser(&pageReq); err != nil {
		return errorpkg.ErrFailParseRequest
	}

	if err := c.val.ValidateStruct(query); err != nil {
		return err
	}
	if err := c.val.ValidateStruct(pageReq); err != nil {
		return err
	}

	courses, pageResp, err := c.svc.GetCourses(ctx.Context(), query, pageReq)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(map[string]interface{}{
		"courses":    courses,
		"pagination": pageResp,
	})
}

func (c *courseHandler) getCourseByID(ctx *fiber.Ctx) error {
	id, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return errorpkg.ErrValidation.Build().WithDetail("invalid course ID")
	}

	resp, err := c.svc.GetCourseByID(ctx.Context(), id)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(map[string]interface{}{
		"course": resp,
	})
}

func (c *courseHandler) updateCourse(ctx *fiber.Ctx) error {
	id, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return errorpkg.ErrValidation.Build().WithDetail("invalid course ID")
	}

	var req dto.UpdateCourseRequest
	if err2 := ctx.BodyParser(&req); err2 != nil {
		return errorpkg.ErrFailParseRequest
	}

	// Handle file uploads
	teacherAvatar, err := ctx.FormFile("teacher_avatar")
	if err == nil {
		// Only set if file was uploaded
		req.TeacherAvatar = teacherAvatar
	}

	thumbnail, err := ctx.FormFile("thumbnail")
	if err == nil {
		// Only set if file was uploaded
		req.Thumbnail = thumbnail
	}

	if err2 := c.val.ValidateStruct(req); err2 != nil {
		return err2
	}

	if err2 := c.svc.UpdateCourse(ctx.Context(), id, &req); err2 != nil {
		return err2
	}

	return ctx.SendStatus(fiber.StatusNoContent)
}

func (c *courseHandler) deleteCourse(ctx *fiber.Ctx) error {
	id, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return errorpkg.ErrValidation.Build().WithDetail("invalid course ID")
	}

	if err2 := c.svc.DeleteCourse(ctx.Context(), id); err2 != nil {
		return err2
	}

	return ctx.SendStatus(fiber.StatusNoContent)
}

func (c *courseHandler) GetPreviewVideoUploadURL(ctx *fiber.Ctx) error {
	id, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return errorpkg.ErrValidation.Build().WithDetail("invalid course ID")
	}

	url, err := c.svc.GetPreviewVideoUploadURL(ctx.Context(), id)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(map[string]interface{}{
		"preview_video_upload_url": url,
	})
}
