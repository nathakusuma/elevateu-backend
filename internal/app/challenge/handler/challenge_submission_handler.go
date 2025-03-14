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

type challengeSubmissionHandler struct {
	val validator.IValidator
	svc contract.IChallengeSubmissionService
}

func InitChallengeSubmissionHandler(
	router fiber.Router,
	midw *middleware.Middleware,
	validator validator.IValidator,
	svc contract.IChallengeSubmissionService,
) {
	handler := challengeSubmissionHandler{
		val: validator,
		svc: svc,
	}

	challengesGroup := router.Group("/challenges")
	challengesGroup.Use(midw.RequireAuthenticated)

	challengesGroup.Post("/submissions/:submission_id/feedbacks",
		midw.RequireOneOfRoles(enum.UserRoleMentor),
		handler.createFeedback)
	challengesGroup.Get("/:challenge_id/submissions/all",
		midw.RequireOneOfRoles(enum.UserRoleMentor, enum.UserRoleAdmin),
		handler.getSubmissionsAsMentor)
	challengesGroup.Post("/:challenge_id/submissions",
		midw.RequireOneOfRoles(enum.UserRoleStudent),
		middleware.RequireSubscription(enum.PaymentTypeChallenge),
		handler.createSubmission)
	challengesGroup.Get("/:challenge_id/submissions",
		midw.RequireOneOfRoles(enum.UserRoleStudent),
		middleware.RequireSubscription(enum.PaymentTypeChallenge),
		handler.getSubmissionAsStudent)
}

func (h *challengeSubmissionHandler) createSubmission(ctx *fiber.Ctx) error {
	var req dto.CreateChallengeSubmissionRequest
	challengeID, err := uuid.Parse(ctx.Params("challenge_id"))
	if err != nil {
		return errorpkg.ErrValidation().WithDetail("Invalid challenge ID")
	}
	req.ChallengeID = challengeID

	if err = ctx.BodyParser(&req); err != nil {
		return errorpkg.ErrFailParseRequest()
	}

	if err = h.val.ValidateStruct(req); err != nil {
		return err
	}

	if err = h.svc.CreateSubmission(ctx.Context(), req); err != nil {
		return err
	}

	return ctx.SendStatus(fiber.StatusCreated)
}

func (h *challengeSubmissionHandler) getSubmissionAsStudent(ctx *fiber.Ctx) error {
	challengeID, err := uuid.Parse(ctx.Params("challenge_id"))
	if err != nil {
		return errorpkg.ErrValidation().WithDetail("Invalid challenge ID")
	}

	submission, err := h.svc.GetSubmissionAsStudent(ctx.Context(), challengeID)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(map[string]interface{}{
		"submission": submission,
	})
}

func (h *challengeSubmissionHandler) getSubmissionsAsMentor(ctx *fiber.Ctx) error {
	challengeID, err := uuid.Parse(ctx.Params("challenge_id"))
	if err != nil {
		return errorpkg.ErrValidation().WithDetail("Invalid challenge ID")
	}

	var pageReq dto.PaginationRequest
	if err = ctx.QueryParser(&pageReq); err != nil {
		return errorpkg.ErrFailParseRequest()
	}

	if err = h.val.ValidateStruct(pageReq); err != nil {
		return err
	}

	submissions, pageResp, err := h.svc.GetSubmissionsAsMentor(ctx.Context(), challengeID, pageReq)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(map[string]interface{}{
		"submissions": submissions,
		"pagination":  pageResp,
	})
}

func (h *challengeSubmissionHandler) createFeedback(ctx *fiber.Ctx) error {
	submissionID, err := uuid.Parse(ctx.Params("submission_id"))
	if err != nil {
		return errorpkg.ErrValidation().WithDetail("Invalid submission ID")
	}

	var req dto.CreateChallengeSubmissionFeedbackRequest
	if err = ctx.BodyParser(&req); err != nil {
		return errorpkg.ErrFailParseRequest()
	}

	if err = h.val.ValidateStruct(req); err != nil {
		return err
	}

	if err = h.svc.CreateFeedback(ctx.Context(), submissionID, req); err != nil {
		return err
	}

	return ctx.SendStatus(fiber.StatusCreated)
}
