package dto

import (
	"mime/multipart"

	"github.com/google/uuid"

	"github.com/nathakusuma/elevateu-backend/domain/entity"
)

type ChallengeGroupResponse struct {
	ID             uuid.UUID `json:"id"`
	Title          string    `json:"title,omitempty"`
	Description    string    `json:"description,omitempty"`
	ChallengeCount *int      `json:"challenge_count,omitempty"`
	ThumbnailURL   string    `json:"thumbnail_url,omitempty"`
}

func (r *ChallengeGroupResponse) PopulateFromEntity(cg *entity.ChallengeGroup,
	urlSigner func(string) (string, error)) error {
	r.ID = cg.ID
	r.Title = cg.Title
	r.Description = cg.Description
	r.ChallengeCount = &cg.ChallengeCount

	var err error
	r.ThumbnailURL, err = urlSigner("challenge_groups/thumbnail/" + cg.ID.String())
	if err != nil {
		return err
	}

	return nil
}

type ChallengeGroupUpdate struct {
	CategoryID  *uuid.UUID `db:"category_id"`
	Title       *string    `db:"title"`
	Description *string    `db:"description"`
}

type CreateChallengeGroupRequest struct {
	CategoryID  uuid.UUID `form:"category_id" validate:"required"`
	Title       string    `form:"title" validate:"required,min=3,max=50"`
	Description string    `form:"description" validate:"required,min=3,max=1000"`
	Thumbnail   *multipart.FileHeader
}

type GetChallengeGroupQuery struct {
	CategoryID *uuid.UUID `query:"category_id"`
	Title      string     `query:"title"`
}

type UpdateChallengeGroupRequest struct {
	CategoryID  *uuid.UUID `form:"category_id" validate:"omitempty"`
	Title       *string    `form:"title" validate:"omitempty,min=3,max=50"`
	Description *string    `form:"description" validate:"omitempty,min=3,max=1000"`
	Thumbnail   *multipart.FileHeader
}
