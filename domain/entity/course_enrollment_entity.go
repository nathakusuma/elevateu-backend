package entity

import (
	"time"

	"github.com/google/uuid"
)

type CourseEnrollment struct {
	CourseID         uuid.UUID `db:"course_id"`
	StudentID        uuid.UUID `db:"student_id"`
	ContentCompleted int       `db:"content_completed"`
	IsCompleted      bool      `db:"is_completed"`
	CreatedAt        time.Time `db:"created_at"`
	LastAccessedAt   time.Time `db:"last_accessed_at"`
}
