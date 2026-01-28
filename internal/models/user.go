// Package model provides data models and entities used throughout the application.
// It defines the structure of domain objects and their database mappings.
package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// User represents a user entity in the system.
// It contains user identification, authentication, and profile information.
// The struct is mapped to the "users" table in the database using GORM tags.
//
// Fields:
//   - ID: Unique identifier (UUID) for the user, automatically generated if not provided
//   - Username: Unique username for login and identification
//   - Password: Bcrypt-hashed password (excluded from JSON responses for security)
//   - DisplayName: User's display name shown in the application
//   - Email: Unique email address for the user account
type User struct {
	ID          string `gorm:"type:uuid;primaryKey;column:id" json:"id"`
	Username    string `gorm:"unique;column:username" json:"username"`
	Password    string `gorm:"column:password" json:"-"` 
	DisplayName string `gorm:"column:display_name" json:"display_name"`
	Email       string `gorm:"column:email;unique" json:"email"`
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
func (u *User) BeforeCreate(_ *gorm.DB) error {
	if u.ID == "" {
		u.ID = uuid.New().String()
	}

	return nil
}
