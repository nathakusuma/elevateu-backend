package handler

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"github.com/nathakusuma/elevateu-backend/domain/contract"
	"github.com/nathakusuma/elevateu-backend/domain/dto"
	"github.com/nathakusuma/elevateu-backend/domain/enum"
	"github.com/nathakusuma/elevateu-backend/domain/errorpkg"
	"github.com/nathakusuma/elevateu-backend/internal/middleware"
	"github.com/nathakusuma/elevateu-backend/pkg/validator"
)

type challengeGroupHandler struct {
	val validator.IValidator
	svc contract.IChallengeGroupService
}

func InitChallengeGroupHandler(
	router fiber.Router,
	midw *middleware.Middleware,
	validator validator.IValidator,
	challengeGroupSvc contract.IChallengeGroupService,
) {
	handler := challengeGroupHandler{
		svc: challengeGroupSvc,
		val: validator,
	}

	groupsRouter := router.Group("/challenge-groups")
	groupsRouter.Use(midw.RequireAuthenticated)

	groupsRouter.Post("",
		midw.RequireOneOfRoles(enum.UserRoleAdmin),
		handler.createGroup)
	groupsRouter.Get("", handler.getGroups)
	groupsRouter.Patch("/:id",
		midw.RequireOneOfRoles(enum.UserRoleAdmin),
		handler.updateGroup)
	groupsRouter.Delete("/:id",
		midw.RequireOneOfRoles(enum.UserRoleAdmin),
		handler.deleteGroup)
}

func (h *challengeGroupHandler) createGroup(ctx *fiber.Ctx) error {
	var req dto.CreateChallengeGroupRequest
	if err := ctx.BodyParser(&req); err != nil {
		return errorpkg.ErrFailParseRequest()
	}

	var err error
	req.Thumbnail, err = ctx.FormFile("thumbnail")
	if err != nil && !errors.Is(err, fiber.ErrUnprocessableEntity) {
		return errorpkg.ErrFailParseRequest().WithDetail("Failed to parse thumbnail")
	}

	if err = h.val.ValidateStruct(req); err != nil {
		return err
	}

	resp, err := h.svc.CreateGroup(ctx.Context(), req)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusCreated).JSON(map[string]interface{}{
		"challenge_group": resp,
	})
}

func (h *challengeGroupHandler) getGroups(ctx *fiber.Ctx) error {
	var query dto.GetChallengeGroupQuery
	if err := ctx.QueryParser(&query); err != nil {
		return errorpkg.ErrFailParseRequest()
	}

	var pageReq dto.PaginationRequest
	if err := ctx.QueryParser(&pageReq); err != nil {
		return errorpkg.ErrFailParseRequest()
	}

	if err := h.val.ValidateStruct(query); err != nil {
		return err
	}
	if err := h.val.ValidateStruct(pageReq); err != nil {
		return err
	}

	groups, pageResp, err := h.svc.GetGroups(ctx.Context(), query, pageReq)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(map[string]interface{}{
		"challenge_groups": groups,
		"pagination":       pageResp,
	})
}

func (h *challengeGroupHandler) updateGroup(ctx *fiber.Ctx) error {
	id, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return errorpkg.ErrValidation().WithDetail("Invalid challenge group ID")
	}

	var req dto.UpdateChallengeGroupRequest
	if err = ctx.BodyParser(&req); err != nil {
		return errorpkg.ErrFailParseRequest()
	}

	thumbnail, err := ctx.FormFile("thumbnail")
	if err == nil {
		req.Thumbnail = thumbnail
	}

	if err = h.val.ValidateStruct(req); err != nil {
		return err
	}

	if err = h.svc.UpdateGroup(ctx.Context(), id, req); err != nil {
		return err
	}

	return ctx.SendStatus(fiber.StatusNoContent)
}

func (h *challengeGroupHandler) deleteGroup(ctx *fiber.Ctx) error {
	id, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return errorpkg.ErrValidation().WithDetail("Invalid challenge group ID")
	}

	if err = h.svc.DeleteGroup(ctx.Context(), id); err != nil {
		return err
	}

	return ctx.SendStatus(fiber.StatusNoContent)
}
