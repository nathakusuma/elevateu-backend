package server

import (
	"errors"

	"github.com/gofiber/fiber/v2"

	"github.com/nathakusuma/elevateu-backend/domain/errorpkg"
	"github.com/nathakusuma/elevateu-backend/internal/infra/env"
	"github.com/nathakusuma/elevateu-backend/pkg/validator"
)

func errorHandler(ctx *fiber.Ctx, err error) error {
	var fiberErr *fiber.Error
	if errors.As(err, &fiberErr) {
		return ctx.SendStatus(fiberErr.Code)
	}

	typePrefix := env.GetEnv().AppURL + "/errors"
	ctx.Set("Content-Type", "application/problem+json")

	var apiErr *errorpkg.ResponseError
	if errors.As(err, &apiErr) {
		return ctx.Status(apiErr.Status).JSON(
			apiErr.
				WithTypePrefix(typePrefix).
				WithInstance(env.GetEnv().AppURL + ctx.OriginalURL()),
		)
	}

	var validationErr validator.ValidationErrors
	if errors.As(err, &validationErr) {
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(
			errorpkg.ErrValidation.
				WithValidationErrors(validationErr).
				WithTypePrefix(typePrefix).
				WithInstance(env.GetEnv().AppURL + ctx.OriginalURL()),
		)
	}

	return ctx.Status(fiber.StatusInternalServerError).JSON(
		errorpkg.ErrInternalServer.
			WithTypePrefix(typePrefix),
	)
}
