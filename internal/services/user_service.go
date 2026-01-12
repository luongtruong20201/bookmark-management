package service

import (
	"context"

	model "github.com/luongtruong20201/bookmark-management/internal/models"
	repository "github.com/luongtruong20201/bookmark-management/internal/repositories"
	"github.com/luongtruong20201/bookmark-management/pkg/utils"
)

type User interface {
	CreateUser(ctx context.Context, username, password, displayName, email string) (*model.User, error)
}

type user struct {
	repo repository.User
}

func NewUser(repo repository.User) User {
	return &user{
		repo: repo,
	}
}

func (u *user) CreateUser(ctx context.Context, username, password, displayName, email string) (*model.User, error) {
	hashPassword := utils.HashPassword(password)

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
