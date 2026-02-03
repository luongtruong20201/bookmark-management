// Package model provides data models and entities used throughout the application.
// It defines the structure of domain objects and their database mappings.
package model

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
	Base
	Username    string `gorm:"unique;column:username" json:"username"`
	Password    string `gorm:"column:password" json:"-"`
	DisplayName string `gorm:"column:display_name" json:"display_name"`
	Email       string `gorm:"column:email;unique" json:"email"`
}
