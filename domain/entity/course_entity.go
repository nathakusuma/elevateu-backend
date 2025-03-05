package entity

import (
	"time"

	"github.com/google/uuid"
)

type Course struct {
	ID              uuid.UUID `db:"id"`
	CategoryID      uuid.UUID `db:"category_id"`
	Title           string    `db:"title"`
	Description     string    `db:"description"`
	TeacherName     string    `db:"teacher_name"`
	Rating          float64   `db:"rating"`
	RatingCount     int64     `db:"rating_count"`
	TotalRating     float64   `db:"total_rating"`
	EnrollmentCount int64     `db:"enrollment_count"`
	ContentCount    int       `db:"content_count"`
	TotalDuration   int       `db:"total_duration"`
	CreatedAt       time.Time `db:"created_at"`
	UpdatedAt       time.Time `db:"updated_at"`

	Category   *Category         `db:"category"`
	Enrollment *CourseEnrollment `db:"enrollment"`
}
