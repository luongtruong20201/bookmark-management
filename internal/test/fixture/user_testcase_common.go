// Package fixture provides reusable database fixtures used across tests.
// This file defines a common user dataset shared by repository and service tests.
package fixture

import (
	model "github.com/luongtruong20201/bookmark-management/internal/models"
	"github.com/luongtruong20201/bookmark-management/pkg/utils"
	"gorm.io/gorm"
)

type UserCommonTestDB struct {
	base
}

// Migrate applies the database schema for the User model used in tests.
func (f *UserCommonTestDB) Migrate() error {
	return f.db.AutoMigrate(&model.User{})
}

// GenerateData seeds a common set of demo users into the test database.
// The dataset is reused across multiple tests for consistent expectations.
func (f *UserCommonTestDB) GenerateData() error {
	db := f.db.Session(&gorm.Session{})
	users := []*model.User{
		{
			Base: model.Base{
				ID: "9b5c1e3e-7c3b-4f4e-8e7c-6e7a2f5d3a91",
			},
			DisplayName: "Nguyen Van An",
			Username:    "an.nguyen",
			Password:    "P@ssw0rd1",
			Email:       "an.nguyen@example.com",
		},
		{
			Base: model.Base{
				ID: "2f6c9d14-9e42-4c31-9e6a-0f8c4b2a7c55",
			},
			DisplayName: "Tran Thi Binh",
			Username:    "binh.tran",
			Password:    "P@ssw0rd2",
			Email:       "binh.tran@example.com",
		},
		{
			Base: model.Base{
				ID: "e3c2a8f1-1d3b-4c62-8e54-6b7f9a2d1c90",
			},
			DisplayName: "Le Quang Huy",
			Username:    "huy.le",
			Password:    "P@ssw0rd3",
			Email:       "huy.le@example.com",
		},
		{
			Base: model.Base{
				ID: "7a9f2d41-5b4c-4f3e-9d21-3e8c6a5b4f72",
			},
			DisplayName: "Pham Minh Duc",
			Username:    "duc.pham",
			Password:    "P@ssw0rd4",
			Email:       "duc.pham@example.com",
		},
		{
			Base: model.Base{
				ID: "c4f7b9d2-2a6e-4c8b-b6d1-1f9e7a2d5c44",
			},
			DisplayName: "Vo Thanh Long",
			Username:    "long.vo",
			Password:    "P@ssw0rd5",
			Email:       "long.vo@example.com",
		},
		{
			Base: model.Base{
				ID: "5e8a3b6f-9c1d-4e7b-a2f6-8d1c9b7f2e30",
			},
			DisplayName: "Do Hoang Nam",
			Username:    "nam.do",
			Password:    "P@ssw0rd6",
			Email:       "nam.do@example.com",
		},
		{
			Base: model.Base{
				ID: "1d9e6f4c-7b2a-4f8e-9c6a-5e3b7d2f1a88",
			},
			DisplayName: "Bui Tuan Kiet",
			Username:    "kiet.bui",
			Password:    "P@ssw0rd7",
			Email:       "kiet.bui@example.com",
		},
		{
			Base: model.Base{
				ID: "8f2a6d9c-4e1b-4f73-9a5d-2c7e6b1f8d55",
			},
			DisplayName: "Dang Ngoc Linh",
			Username:    "linh.dang",
			Password:    "P@ssw0rd8",
			Email:       "linh.dang@example.com",
		},
		{
			Base: model.Base{
				ID: "3c6b8d5e-1f9a-4e7d-8c2b-6a9f5d1e7b44",
			},
			DisplayName: "Hoang Gia Bao",
			Username:    "bao.hoang",
			Password:    "P@ssw0rd9",
			Email:       "bao.hoang@example.com",
		},
		{
			Base: model.Base{
				ID: "a7c1d8f9-5e3b-4a6c-9d2f-8b6e1c7a4f20",
			},
			DisplayName: "Nguyen Thi Mai",
			Username:    "mai.nguyen",
			Password:    "P@ssw0rd10",
			Email:       "mai.nguyen@example.com",
		},
		{
			Base: model.Base{
				ID: "550e8400-e29b-41d4-a716-446655440000",
			},
			DisplayName: "John Doe",
			Username:    "johndoe",
			Password:    "P@ssw0rd11",
			Email:       "john.doe@example.com",
		},
	}

	hasher := utils.NewHasher()
	for _, u := range users {
		u.Password = hasher.HashPassword(u.Password)
	}

	return db.CreateInBatches(users, len(users)).Error
}
