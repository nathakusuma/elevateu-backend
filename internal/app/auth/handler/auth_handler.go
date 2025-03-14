package handler

import (
	"github.com/nathakusuma/elevateu-backend/pkg/log"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"github.com/nathakusuma/elevateu-backend/domain/contract"
	"github.com/nathakusuma/elevateu-backend/domain/ctxkey"
	"github.com/nathakusuma/elevateu-backend/domain/dto"
	"github.com/nathakusuma/elevateu-backend/domain/errorpkg"
	"github.com/nathakusuma/elevateu-backend/internal/middleware"
	"github.com/nathakusuma/elevateu-backend/pkg/validator"
)

type authHandler struct {
	val validator.IValidator
	svc contract.IAuthService
}

func InitAuthHandler(
	router fiber.Router,
	middlewareInstance *middleware.Middleware,
	validator validator.IValidator,
	authSvc contract.IAuthService,
) {
	handler := authHandler{
		svc: authSvc,
		val: validator,
	}

	authGroup := router.Group("/auth")
	authGroup.Post("/register/otp", handler.requestOTPRegister)
	authGroup.Post("/register", handler.register)
	authGroup.Post("/login", handler.login)
	authGroup.Post("/refresh", handler.refreshToken)
	authGroup.Post("/logout", middlewareInstance.RequireAuthenticated, handler.logout)
	authGroup.Post("/reset-password/otp", handler.requestOTPResetPassword)
	authGroup.Post("/reset-password", handler.resetPassword)
}

func (c *authHandler) requestOTPRegister(ctx *fiber.Ctx) error {
	type request struct {
		Email string `json:"email" validate:"required,email,max=320"`
	}

	var req request
	if err := ctx.BodyParser(&req); err != nil {
		return errorpkg.ErrFailParseRequest()
	}

	if err := c.val.ValidateStruct(req); err != nil {
		return err
	}

	err := c.svc.RequestRegisterOTP(ctx.Context(), req.Email)
	if err != nil {
		return err
	}

	return ctx.SendStatus(http.StatusNoContent)
}

func (c *authHandler) register(ctx *fiber.Ctx) error {
	var req dto.RegisterRequest
	if err := ctx.BodyParser(&req); err != nil {
		return errorpkg.ErrFailParseRequest()
	}

	if err := c.val.ValidateStruct(req); err != nil {
		return err
	}

	resp, err := c.svc.Register(ctx.Context(), req)
	if err != nil {
		return err
	}

	return ctx.Status(http.StatusCreated).JSON(resp)
}

func (c *authHandler) login(ctx *fiber.Ctx) error {
	var req dto.LoginRequest
	if err := ctx.BodyParser(&req); err != nil {
		return errorpkg.ErrFailParseRequest()
	}

	if err := c.val.ValidateStruct(req); err != nil {
		return err
	}

	resp, err := c.svc.Login(ctx.Context(), req)
	if err != nil {
		return err
	}

	return ctx.Status(http.StatusOK).JSON(resp)
}

func (c *authHandler) refreshToken(ctx *fiber.Ctx) error {
	type request struct {
		RefreshToken string `json:"refresh_token" validate:"required"`
	}

	var req request
	if err := ctx.BodyParser(&req); err != nil {
		return errorpkg.ErrFailParseRequest()
	}

	if err := c.val.ValidateStruct(req); err != nil {
		return err
	}

	resp, err := c.svc.Refresh(ctx.Context(), req.RefreshToken)
	if err != nil {
		return err
	}

	return ctx.Status(http.StatusOK).JSON(resp)
}

func (c *authHandler) logout(ctx *fiber.Ctx) error {
	userID, ok := ctx.Locals(ctxkey.UserID).(uuid.UUID)
	if !ok {
		traceID := log.ErrorWithTraceID(ctx.Context(), nil, "Failed to get user ID from context")
		return errorpkg.ErrInternalServer().WithTraceID(traceID)
	}

	err := c.svc.Logout(ctx.Context(), userID)
	if err != nil {
		return err
	}

	return ctx.SendStatus(http.StatusNoContent)
}

func (c *authHandler) requestOTPResetPassword(ctx *fiber.Ctx) error {
	type request struct {
		Email string `json:"email" validate:"required,email"`
	}

	var req request
	if err := ctx.BodyParser(&req); err != nil {
		return errorpkg.ErrFailParseRequest()
	}

	if err := c.val.ValidateStruct(req); err != nil {
		return err
	}

	err := c.svc.RequestPasswordResetOTP(ctx.Context(), req.Email)
	if err != nil {
		return err
	}

	return ctx.SendStatus(http.StatusNoContent)
}

func (c *authHandler) resetPassword(ctx *fiber.Ctx) error {
	var req dto.ResetPasswordRequest
	if err := ctx.BodyParser(&req); err != nil {
		return errorpkg.ErrFailParseRequest()
	}

	if err := c.val.ValidateStruct(req); err != nil {
		return err
	}

	resp, err := c.svc.ResetPassword(ctx.Context(), req)
	if err != nil {
		return err
	}

	return ctx.Status(http.StatusOK).JSON(resp)
}
