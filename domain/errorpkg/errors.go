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

	ErrValidation = NewError(http.StatusUnprocessableEntity, "validation-error",
		"There are invalid fields in your request. Please check and try again")
)
