package dto

import (
	"time"

	"github.com/google/uuid"

	"github.com/nathakusuma/elevateu-backend/domain/entity"
)

type ChallengeSubmissionFeedbackResponse struct {
	Score           *int      `json:"score"`
	Feedback        string    `json:"feedback"`
	MentorName      string    `json:"mentor_name"`
	MentorAvatarURL string    `json:"mentor_avatar_url"`
	CreatedAt       time.Time `json:"created_at"`
}

type ChallengeSubmissionResponse struct {
	ID               uuid.UUID                            `json:"id,omitempty"`
	URL              string                               `json:"url,omitempty"`
	StudentName      string                               `json:"student_name,omitempty"`
	StudentAvatarURL string                               `json:"student_avatar_url,omitempty"`
	CreatedAt        time.Time                            `json:"created_at,omitempty"`
	Feedback         *ChallengeSubmissionFeedbackResponse `json:"feedback,omitempty"`
}

func (r *ChallengeSubmissionFeedbackResponse) populateFromEntity(feedback *entity.ChallengeSubmissionFeedback,
	urlSigner func(string) (string, error)) error {
	r.Score = &feedback.Score
	r.Feedback = feedback.Feedback

	if !feedback.CreatedAt.IsZero() {
		r.CreatedAt = feedback.CreatedAt
	}

	if feedback.Mentor != nil {
		r.MentorName = feedback.Mentor.Name
		var err error
		if feedback.Mentor.HasAvatar {
			r.MentorAvatarURL, err = urlSigner("users/avatar/" + feedback.Mentor.ID.String())
		} else {
			r.MentorAvatarURL, err = urlSigner("users/avatar/default")
		}
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *ChallengeSubmissionResponse) PopulateFromEntity(submission *entity.ChallengeSubmission,
	urlSigner func(string) (string, error)) error {
	r.ID = submission.ID
	r.URL = submission.URL

	if !submission.CreatedAt.IsZero() {
		r.CreatedAt = submission.CreatedAt
	}

	if submission.Student != nil {
		r.StudentName = submission.Student.Name
		var err error
		if submission.Student.HasAvatar {
			r.StudentAvatarURL, err = urlSigner("users/avatar/" + submission.Student.ID.String())
		} else {
			r.StudentAvatarURL, err = urlSigner("users/avatar/default")
		}
		if err != nil {
			return err
		}
	}

	if submission.Feedback != nil {
		r.Feedback = &ChallengeSubmissionFeedbackResponse{}
		if err := r.Feedback.populateFromEntity(submission.Feedback, urlSigner); err != nil {
			return err
		}
	}

	return nil
}

type CreateChallengeSubmissionRequest struct {
	ChallengeID uuid.UUID `param:"challenge_id" validate:"required"`
	URL         string    `json:"url" validate:"required,url"`
}

type CreateChallengeSubmissionFeedbackRequest struct {
	Score    *int   `json:"score" validate:"required,min=0,max=100"`
	Feedback string `json:"feedback" validate:"required,min=3,max=1000"`
}
