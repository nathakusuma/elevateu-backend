package handler

import (
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"github.com/nathakusuma/elevateu-backend/domain/contract"
	"github.com/nathakusuma/elevateu-backend/domain/ctxkey"
	"github.com/nathakusuma/elevateu-backend/domain/dto"
	"github.com/nathakusuma/elevateu-backend/domain/errorpkg"
	"github.com/nathakusuma/elevateu-backend/internal/middleware"
	"github.com/nathakusuma/elevateu-backend/pkg/log"
	"github.com/nathakusuma/elevateu-backend/pkg/validator"
)

type authHandler struct {
	val validator.IValidator
	svc contract.IAuthService
}

func InitAuthHandler(
	router fiber.Router,
	midw *middleware.Middleware,
	validator validator.IValidator,
	authSvc contract.IAuthService,
) {
	handler := authHandler{
		svc: authSvc,
		val: validator,
	}

	authGroup := router.Group("/auth")

	otpRateLimiter := midw.RateLimit(middleware.RateLimitConfig{
		MaxRequests:    5,
		PerTimeWindow:  15 * time.Minute,
		ExpirationTime: 30 * time.Minute,
	})

	loginRateLimiter := midw.RateLimit(middleware.RateLimitConfig{
		MaxRequests:    5,
		PerTimeWindow:  5 * time.Minute,
		ExpirationTime: 10 * time.Minute,
	})

	registrationRateLimiter := midw.RateLimit(middleware.RateLimitConfig{
		MaxRequests:    3,
		PerTimeWindow:  15 * time.Minute,
		ExpirationTime: 30 * time.Minute,
	})

	refreshRateLimiter := midw.RateLimit(middleware.RateLimitConfig{
		MaxRequests:    10,
		PerTimeWindow:  5 * time.Minute,
		ExpirationTime: 10 * time.Minute,
	})

	authGroup.Post("/register/otp",
		otpRateLimiter,
		handler.requestOTPRegister)
	authGroup.Post("/register",
		registrationRateLimiter,
		handler.register)
	authGroup.Post("/login",
		loginRateLimiter,
		handler.login)
	authGroup.Post("/refresh",
		refreshRateLimiter,
		handler.refreshToken)
	authGroup.Post("/logout",
		midw.RequireAuthenticated,
		handler.logout)
	authGroup.Post("/reset-password/otp",
		otpRateLimiter,
		handler.requestOTPResetPassword)
	authGroup.Post("/reset-password",
		registrationRateLimiter,
		handler.resetPassword)
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
