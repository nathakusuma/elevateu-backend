package errorpkg

import (
	"strings"

	"github.com/google/uuid"

	"github.com/nathakusuma/elevateu-backend/pkg/validator"
)

type ResponseError struct {
	Type             string                     `json:"type"`
	Title            string                     `json:"title"`
	Status           int                        `json:"status"`
	Detail           string                     `json:"detail,omitempty"`
	Instance         string                     `json:"instance,omitempty"`
	TraceID          *uuid.UUID                 `json:"trace_id,omitempty"`
	ValidationErrors validator.ValidationErrors `json:"validation_errors,omitempty"`
}

func (e *ResponseError) Error() string {
	return e.Title
}

func NewError(status int, errType, title string) *ResponseError {
	return &ResponseError{
		Type:   errType,
		Title:  title,
		Status: status,
	}
}

func (e *ResponseError) WithTypePrefix(prefix string) *ResponseError {
	// Remove trailing slash from prefix if it exists
	prefix = strings.TrimSuffix(prefix, "/")

	// Clean the existing type by removing any existing prefixes
	cleanType := e.Type
	for strings.HasPrefix(cleanType, prefix) {
		cleanType = strings.TrimPrefix(cleanType, prefix+"/")
	}

	e.Type = prefix + "/" + cleanType
	return e
}

func (e *ResponseError) WithDetail(detail string) *ResponseError {
	e.Detail = detail
	return e
}

func (e *ResponseError) WithInstance(instance string) *ResponseError {
	e.Instance = instance
	return e
}

func (e *ResponseError) WithTraceID(traceID uuid.UUID) *ResponseError {
	e.TraceID = &traceID
	return e
}

func (e *ResponseError) WithValidationErrors(validationErrors validator.ValidationErrors) *ResponseError {
	e.ValidationErrors = validationErrors
	return e
}
