package dto

type UpdateCourseVideoProgressRequest struct {
	LastPosition int  `json:"last_position" validate:"required,gte=0"`
	IsCompleted  bool `json:"is_completed"`
}
