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

type courseFeedbackHandler struct {
	val validator.IValidator
	svc contract.ICourseFeedbackService
}

func InitCourseFeedbackHandler(
	router fiber.Router,
	midw *middleware.Middleware,
	validator validator.IValidator,
	feedbackSvc contract.ICourseFeedbackService,
) {
	handler := courseFeedbackHandler{
		svc: feedbackSvc,
		val: validator,
	}

	coursesGroup := router.Group("/courses")
	coursesGroup.Use(midw.RequireAuthenticated)

	coursesGroup.Patch("feedbacks/:id",
		midw.RequireOneOfRoles(enum.UserRoleStudent),
		handler.updateFeedback)
	coursesGroup.Delete("feedbacks/:id",
		midw.RequireOneOfRoles(enum.UserRoleStudent),
		handler.deleteFeedback)
	coursesGroup.Post("/:courseId/feedbacks",
		midw.RequireOneOfRoles(enum.UserRoleStudent),
		handler.createFeedback)
	coursesGroup.Get("/:courseId/feedbacks", handler.getFeedbacks)
}

func (h *courseFeedbackHandler) createFeedback(ctx *fiber.Ctx) error {
	courseID, err := uuid.Parse(ctx.Params("courseId"))
	if err != nil {
		return errorpkg.ErrValidation().WithDetail("Invalid course ID")
	}

	var req dto.CreateCourseFeedbackRequest
	if err = ctx.BodyParser(&req); err != nil {
		return errorpkg.ErrFailParseRequest()
	}

	if err = h.val.ValidateStruct(req); err != nil {
		return err
	}

	if err = h.svc.CreateFeedback(ctx.Context(), courseID, req); err != nil {
		return err
	}

	return ctx.SendStatus(fiber.StatusCreated)
}

func (h *courseFeedbackHandler) getFeedbacks(c *fiber.Ctx) error {
	courseID, err := uuid.Parse(c.Params("courseId"))
	if err != nil {
		return errorpkg.ErrValidation().WithDetail("Invalid course ID")
	}

	// Parse pagination request
	var pageReq dto.PaginationRequest
	if err = c.QueryParser(&pageReq); err != nil {
		return errorpkg.ErrFailParseRequest()
	}

	if err = h.val.ValidateStruct(pageReq); err != nil {
		return err
	}

	feedbacks, pageResp, err := h.svc.GetFeedbacksByCourseID(c.Context(), courseID, pageReq)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(map[string]interface{}{
		"feedbacks":  feedbacks,
		"pagination": pageResp,
	})
}

func (h *courseFeedbackHandler) updateFeedback(c *fiber.Ctx) error {
	feedbackID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return errorpkg.ErrValidation().WithDetail("Invalid feedback ID")
	}

	var req dto.UpdateCourseFeedbackRequest
	if err = c.BodyParser(&req); err != nil {
		return errorpkg.ErrFailParseRequest()
	}

	if err = h.val.ValidateStruct(req); err != nil {
		return err
	}

	if err = h.svc.UpdateFeedback(c.Context(), feedbackID, req); err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func (h *courseFeedbackHandler) deleteFeedback(c *fiber.Ctx) error {
	feedbackID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return errorpkg.ErrValidation().WithDetail("Invalid feedback ID")
	}

	if err = h.svc.DeleteFeedback(c.Context(), feedbackID); err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusNoContent)
}
