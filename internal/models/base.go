package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Base struct {
	ID        string     `gorm:"type:uuid;primaryKey;column:id" json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"-"`
}

// BeforeCreate is a GORM hook that automatically generates a UUID for the user
// if the ID field is empty before creating the record in the database.
// This ensures that every user has a unique identifier even if one is not explicitly provided.
//
// Parameters:
//   - _: GORM database instance (unused but required by GORM hook signature)
//
// Returns:
//   - error: Always returns nil, as UUID generation cannot fail
//
// Note: This hook is called automatically by GORM before inserting a new user record.
func (u *Base) BeforeCreate(_ *gorm.DB) error {
	if u.ID == "" {
		u.ID = uuid.New().String()
	}

	return nil
}
