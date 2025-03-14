package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"github.com/nathakusuma/elevateu-backend/domain/contract"
	"github.com/nathakusuma/elevateu-backend/domain/enum"
	"github.com/nathakusuma/elevateu-backend/domain/errorpkg"
	"github.com/nathakusuma/elevateu-backend/internal/middleware"
	"github.com/nathakusuma/elevateu-backend/pkg/validator"
)

type categoryHandler struct {
	svc contract.ICategoryService
	val validator.IValidator
}

func InitCategoryHandler(
	router fiber.Router,
	svc contract.ICategoryService,
	midw *middleware.Middleware,
	validator validator.IValidator,
) {
	handler := categoryHandler{
		svc: svc,
		val: validator,
	}

	categoryGroup := router.Group("/categories")
	categoryGroup.Use(midw.RequireAuthenticated)

	categoryGroup.Post("/",
		midw.RequireOneOfRoles(enum.UserRoleAdmin),
		handler.createCategory)
	categoryGroup.Get("/", handler.getAllCategories)
	categoryGroup.Put("/:id",
		midw.RequireOneOfRoles(enum.UserRoleAdmin),
		handler.updateCategory)
	categoryGroup.Delete("/:id",
		midw.RequireOneOfRoles(enum.UserRoleAdmin),
		handler.deleteCategory)
}

func (c *categoryHandler) createCategory(ctx *fiber.Ctx) error {
	type request struct {
		Name string `json:"name" validate:"required,min=3,max=50"`
	}

	var req request
	if err := ctx.BodyParser(&req); err != nil {
		return errorpkg.ErrFailParseRequest()
	}

	if err := c.val.ValidateStruct(req); err != nil {
		return err
	}

	category, err := c.svc.CreateCategory(ctx.Context(), req.Name)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusCreated).JSON(map[string]interface{}{
		"category": category,
	})
}

func (c *categoryHandler) getAllCategories(ctx *fiber.Ctx) error {
	categories, err := c.svc.GetAllCategories(ctx.Context())
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(map[string]interface{}{
		"categories": categories,
	})
}

func (c *categoryHandler) updateCategory(ctx *fiber.Ctx) error {
	type request struct {
		ID   string `param:"id" validate:"required,uuid"`
		Name string `json:"name" validate:"required,min=3,max=50"`
	}

	var req request
	if err := ctx.ParamsParser(&req); err != nil {
		return errorpkg.ErrFailParseRequest()
	}
	if err := ctx.BodyParser(&req); err != nil {
		return errorpkg.ErrFailParseRequest()
	}

	if err := c.val.ValidateStruct(req); err != nil {
		return err
	}

	id, err := uuid.Parse(req.ID)
	if err != nil {
		return errorpkg.ErrValidation().WithDetail("Invalid category ID")
	}

	if err2 := c.svc.UpdateCategory(ctx.Context(), id, req.Name); err2 != nil {
		return err2
	}

	return ctx.SendStatus(fiber.StatusNoContent)
}

func (c *categoryHandler) deleteCategory(ctx *fiber.Ctx) error {
	id, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return errorpkg.ErrValidation().WithDetail("Invalid category ID")
	}

	if err2 := c.svc.DeleteCategory(ctx.Context(), id); err2 != nil {
		return err2
	}

	return ctx.SendStatus(fiber.StatusNoContent)
}
