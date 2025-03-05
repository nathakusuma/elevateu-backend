package entity

import (
	"time"

	"github.com/google/uuid"
)

type CourseFeedback struct {
	ID        uuid.UUID `db:"id"`
	CourseID  uuid.UUID `db:"course_id"`
	StudentID uuid.UUID `db:"student_id"`
	Rating    float64   `db:"rating"`
	Comment   string    `db:"comment"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`

	User *User `db:"user"`
}
