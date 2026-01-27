package service

import (
	"context"

	model "github.com/luongtruong20201/bookmark-management/internal/models"
	repository "github.com/luongtruong20201/bookmark-management/internal/repositories"
	"github.com/luongtruong20201/bookmark-management/pkg/utils"
)

// User defines the interface for user service operations.
// It provides methods to handle user-related business logic.
//
//go:generate mockery --name User --filename user_service.go
type User interface {
	CreateUser(ctx context.Context, username, password, displayName, email string) (*model.User, error)
}

type user struct {
	repo   repository.User
	hasher utils.Hasher
}

// NewUser creates a new user service with the provided user repository.
func NewUser(repo repository.User, hasher utils.Hasher) User {
	return &user{
		repo:   repo,
		hasher: hasher,
	}
}

// CreateUser creates a new user account. It hashes the password before storing
// the user information in the database. Returns the created user with generated UUID.
func (u *user) CreateUser(ctx context.Context, username, password, displayName, email string) (*model.User, error) {
	hashPassword := u.hasher.HashPassword(password)

	newUser := model.User{
		Username:    username,
		DisplayName: displayName,
		Email:       email,
		Password:    hashPassword,
	}

	res, err := u.repo.CreateUser(ctx, &newUser)
	if err != nil {
		return nil, err
	}

	return res, nil
}
