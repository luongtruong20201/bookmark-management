package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// User represents a user entity in the system.
// It contains user identification and authentication information.
type User struct {
	ID          string `gorm:"type:uuid;primaryKey;column:id" json:"id"`
	Username    string `gorm:"unique;column:username" json:"username"`
	Password    string `gorm:"column:password" json:"-"` 
	DisplayName string `gorm:"column:display_name" json:"display_name"`
	Email       string `gorm:"column:email,unique" json:"email"`
}

// BeforeCreate is a GORM hook that automatically generates a UUID for the user
// if the ID field is empty before creating the record in the database.
func (u *User) BeforeCreate(_ *gorm.DB) error {
	if u.ID == "" {
		u.ID = uuid.New().String()
	}

	return nil
}
