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

type challengeHandler struct {
	val validator.IValidator
	svc contract.IChallengeService
}

func InitChallengeHandler(
	router fiber.Router,
	midw *middleware.Middleware,
	validator validator.IValidator,
	challengeSvc contract.IChallengeService,
) {
	handler := challengeHandler{
		svc: challengeSvc,
		val: validator,
	}

	challengesRouter := router.Group("/challenges")
	challengesRouter.Use(midw.RequireAuthenticated)

	challengesRouter.Post("",
		midw.RequireOneOfRoles(enum.UserRoleAdmin),
		handler.createChallenge)
	challengesRouter.Get("", handler.getChallenges)
	challengesRouter.Get("/:id",
		middleware.RequireSubscription(enum.PaymentTypeChallenge),
		handler.getChallengeDetail)
	challengesRouter.Patch("/:id",
		midw.RequireOneOfRoles(enum.UserRoleAdmin),
		handler.updateChallenge)
	challengesRouter.Delete("/:id",
		midw.RequireOneOfRoles(enum.UserRoleAdmin),
		handler.deleteChallenge)
}

func (h *challengeHandler) createChallenge(ctx *fiber.Ctx) error {
	var req dto.CreateChallengeRequest
	if err := ctx.BodyParser(&req); err != nil {
		return errorpkg.ErrFailParseRequest
	}

	if err := h.val.ValidateStruct(req); err != nil {
		return err
	}

	resp, err := h.svc.CreateChallenge(ctx.Context(), &req)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusCreated).JSON(map[string]interface{}{
		"challenge": resp,
	})
}

func (h *challengeHandler) getChallenges(ctx *fiber.Ctx) error {
	var query struct {
		GroupID    uuid.UUID                `query:"group_id" validate:"required"`
		Difficulty enum.ChallengeDifficulty `query:"difficulty" validate:"required,oneof=beginner intermediate advanced"`
	}

	if err := ctx.QueryParser(&query); err != nil {
		return errorpkg.ErrFailParseRequest
	}

	var pageReq dto.PaginationRequest
	if err := ctx.QueryParser(&pageReq); err != nil {
		return errorpkg.ErrFailParseRequest
	}

	if err := h.val.ValidateStruct(query); err != nil {
		return err
	}
	if err := h.val.ValidateStruct(pageReq); err != nil {
		return err
	}

	challenges, pageResp, err := h.svc.GetChallenges(ctx.Context(), query.GroupID, query.Difficulty, pageReq)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(map[string]interface{}{
		"challenges": challenges,
		"pagination": pageResp,
	})
}

func (h *challengeHandler) getChallengeDetail(ctx *fiber.Ctx) error {
	id, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return errorpkg.ErrValidation.Build().WithDetail("Invalid challenge ID")
	}

	challenge, err := h.svc.GetChallengeDetail(ctx.Context(), id)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(map[string]interface{}{
		"challenge": challenge,
	})
}

func (h *challengeHandler) updateChallenge(ctx *fiber.Ctx) error {
	id, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return errorpkg.ErrValidation.Build().WithDetail("Invalid challenge ID")
	}

	var req dto.UpdateChallengeRequest
	if err = ctx.BodyParser(&req); err != nil {
		return errorpkg.ErrFailParseRequest
	}

	if err = h.val.ValidateStruct(req); err != nil {
		return err
	}

	if err = h.svc.UpdateChallenge(ctx.Context(), id, &req); err != nil {
		return err
	}

	return ctx.SendStatus(fiber.StatusNoContent)
}

func (h *challengeHandler) deleteChallenge(ctx *fiber.Ctx) error {
	id, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return errorpkg.ErrValidation.Build().WithDetail("Invalid challenge ID")
	}

	if err = h.svc.DeleteChallenge(ctx.Context(), id); err != nil {
		return err
	}

	return ctx.SendStatus(fiber.StatusNoContent)
}
