package repository

import (
	"context"

	model "github.com/luongtruong20201/bookmark-management/internal/models"
	"gorm.io/gorm"
)

//go:generate mockery --name User --filename user.go

// User defines the interface for user repository operations.
// It provides methods to interact with user data in the database.
type User interface {
	CreateUser(context.Context, *model.User) (*model.User, error)
}

type user struct {
	db *gorm.DB
}

// NewUser creates a new user repository with the provided database connection.
func NewUser(db *gorm.DB) User {
	return &user{
		db: db,
	}
}

// CreateUser persists a new user record in the database. The user ID will be
// automatically generated as a UUID if not provided (via BeforeCreate hook).
// Returns the created user with the generated ID.
func (u *user) CreateUser(ctx context.Context, user *model.User) (*model.User, error) {
	if err := u.db.WithContext(ctx).Create(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}
