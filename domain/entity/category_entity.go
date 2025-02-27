package entity

import "github.com/google/uuid"

type Category struct {
	ID   uuid.UUID `db:"id" json:"id,omitempty"`
	Name string    `db:"name" json:"name,omitempty"`
}
