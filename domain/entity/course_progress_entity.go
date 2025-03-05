package entity

import "github.com/google/uuid"

type CourseVideoProgress struct {
	StudentID    uuid.UUID `db:"student_id"`
	VideoID      uuid.UUID `db:"video_id"`
	LastPosition int       `db:"last_position"`
	IsCompleted  bool      `db:"is_completed"`
}

type CourseMaterialProgress struct {
	StudentID  uuid.UUID `db:"student_id"`
	MaterialID uuid.UUID `db:"material_id"`
}
