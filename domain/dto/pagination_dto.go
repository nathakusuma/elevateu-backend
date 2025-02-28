package dto

import "github.com/google/uuid"

type PaginationRequest struct {
	Cursor    uuid.UUID `query:"cursor" validate:"omitempty,required_with=Direction"`
	Limit     int       `query:"limit" validate:"required,min=1,max=10"`
	Direction string    `query:"direction" validate:"omitempty,required_with=Cursor,oneof=next prev"`
}

type PaginationResponse struct {
	HasMore bool `json:"has_more"`
}
