package repository

import (
	"context"

	model "github.com/luongtruong20201/bookmark-management/internal/models"
	"gorm.io/gorm"
)

type User interface {
	CreateUser(context.Context, *model.User) (*model.User, error)
}

type user struct {
	db *gorm.DB
}

func NewUser(db *gorm.DB) User {
	return &user{
		db: db,
	}
}

func (u *user) CreateUser(ctx context.Context, user *model.User) (*model.User, error) {
	if err := u.db.WithContext(ctx).Create(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}
