package dto

import (
	"time"

	"github.com/google/uuid"

	"github.com/nathakusuma/elevateu-backend/domain/entity"
)

type CourseFeedbackResponse struct {
	ID        uuid.UUID `json:"id,omitempty"`
	Rating    float64   `json:"rating,omitempty"`
	Feedback  string    `json:"feedback,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`

	StudentName      string `json:"student_name,omitempty"`
	StudentAvatarURL string `json:"student_avatar_url,omitempty"`
}

func (c *CourseFeedbackResponse) PopulateFromEntity(feedback *entity.CourseFeedback,
	urlSigner func(string) (string, error)) error {
	var err error
	if feedback.User.HasAvatar {
		c.StudentAvatarURL, err = urlSigner("users/avatar/" + feedback.User.ID.String())
	} else {
		c.StudentAvatarURL, err = urlSigner("users/avatar/default")
	}
	if err != nil {
		return err
	}

	c.ID = feedback.ID
	c.Rating = feedback.Rating
	c.Feedback = feedback.Comment
	c.CreatedAt = feedback.CreatedAt
	c.UpdatedAt = feedback.UpdatedAt
	c.StudentName = feedback.User.Name

	return nil
}

type CourseFeedbackUpdate struct {
	Rating   *float64 `db:"rating"`
	Feedback *string  `db:"comment"`
}

type CreateCourseFeedbackRequest struct {
	Rating  int    `json:"rating" validate:"required,gte=1,lte=5"`
	Comment string `json:"comment" validate:"required,min=3,max=500"`
}

type UpdateCourseFeedbackRequest struct {
	Rating  int    `json:"rating" validate:"omitempty,gte=1,lte=5"`
	Comment string `json:"comment" validate:"omitempty,min=3,max=500"`
}
