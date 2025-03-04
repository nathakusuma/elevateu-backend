package entity

import (
	"time"

	"github.com/google/uuid"
)

type CourseVideo struct {
	ID          uuid.UUID `db:"id"`
	CourseID    uuid.UUID `db:"course_id"`
	Title       string    `db:"title"`
	Description string    `db:"description"`
	Duration    int       `db:"duration"`
	IsFree      bool      `db:"is_free"`
	Order       int       `db:"order"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

type CourseMaterial struct {
	ID        uuid.UUID `db:"id"`
	CourseID  uuid.UUID `db:"course_id"`
	Title     string    `db:"title"`
	Subtitle  string    `db:"subtitle"`
	IsFree    bool      `db:"is_free"`
	Order     int       `db:"order"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
