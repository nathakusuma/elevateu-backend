package errorpkg

import (
	"net/http"
)

// General
var (
	ErrInternalServer = NewError(http.StatusInternalServerError,
		"internal-server-error",
		"Something went wrong in our server. Please try again later.")

	ErrFailParseRequest = NewError(http.StatusBadRequest,
		"fail-parse-request",
		"Failed to parse request. Please check your request format.")

	ErrForbiddenRole = NewError(http.StatusForbidden,
		"forbidden-role",
		"You're not allowed to access this resource.")

	ErrForbiddenUser = NewError(http.StatusForbidden,
		"forbidden-user",
		"You're not allowed to access this resource.")

	ErrNotFound = NewError(http.StatusNotFound,
		"not-found",
		"Resource not found.")

	ErrFileTooLarge = NewError(http.StatusRequestEntityTooLarge,
		"file-too-large",
		"File size is too large. Please upload smaller file.")

	ErrInvalidFileFormat = NewError(http.StatusUnprocessableEntity,
		"invalid-file-format",
		"Invalid file format. Please upload a valid file.")

	ErrValidation = NewError(http.StatusUnprocessableEntity, "validation-error",
		"There are invalid fields in your request. Please check and try again")
)

// Auth
var (
	ErrCredentialsNotMatch = NewError(http.StatusUnauthorized,
		"credentials-not-match",
		"Credentials do not match. Please try again.")

	ErrInvalidBearerToken = NewError(http.StatusUnauthorized,
		"invalid-bearer-token",
		"Your auth session is invalid. Please renew your auth session.")

	ErrInvalidOTP = NewError(http.StatusUnauthorized,
		"invalid-otp",
		"Invalid OTP. Please try again or request a new OTP.")

	ErrInvalidRefreshToken = NewError(http.StatusUnauthorized,
		"invalid-refresh-token",
		"Auth session is invalid. Please login again.")

	ErrNoBearerToken = NewError(http.StatusUnauthorized,
		"no-bearer-token",
		"You're not logged in. Please login first.")

	ErrEmailAlreadyRegistered = NewError(http.StatusConflict,
		"email-already-registered",
		"Email already registered. Please login or use another email.")
)
