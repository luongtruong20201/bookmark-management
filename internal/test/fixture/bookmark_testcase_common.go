package fixture

import (
	model "github.com/luongtruong20201/bookmark-management/internal/models"
	"gorm.io/gorm"
)

// BookmarkCommonTestDB provides a shared bookmark dataset backed by a test database.
// It reuses the common user dataset from UserCommonTestDB and seeds a deterministic
// set of bookmarks used across repository, service, and endpoint tests.
type BookmarkCommonTestDB struct {
	base
}

// Migrate applies the database schema for users and bookmarks used in tests.
func (f *BookmarkCommonTestDB) Migrate() error {
	return f.db.AutoMigrate(&model.User{}, &model.Bookmark{})
}

// GenerateData seeds common users (via UserCommonTestDB) and a fixed set of
// bookmarks for multiple users. The IDs, descriptions, and user relationships
// are chosen to satisfy expectations in tests that assert on specific IDs,
// ordering, and ownership.
func (f *BookmarkCommonTestDB) GenerateData() error {
	userFixture := &UserCommonTestDB{}
	userFixture.SetupDB(f.db)
	if err := userFixture.GenerateData(); err != nil {
		return err
	}

	db := f.db.Session(&gorm.Session{})

	bookmarks := []*model.Bookmark{
		{
			Base: model.Base{
				ID: "a1b2c3d4-e5f6-4a7b-8c9d-0e1f2a3b4c5d",
			},
			Description: "Facebook - Social Media Platform",
			URL:         "https://www.facebook.com",
			Code:        "abc1234",
			UserID:      "9b5c1e3e-7c3b-4f4e-8e7c-6e7a2f5d3a91",
		},
		{
			Base: model.Base{
				ID: "b2c3d4e5-f6a7-4b8c-9d0e-1f2a3b4c5d6e",
			},
			Description: "Google - Search Engine",
			URL:         "https://www.google.com",
			Code:        "def5678",
			UserID:      "9b5c1e3e-7c3b-4f4e-8e7c-6e7a2f5d3a91",
		},
		{
			Base: model.Base{
				ID: "c3d4e5f6-a7b8-4c9d-0e1f-2a3b4c5d6e7f",
			},
			Description: "GitHub - Code Repository",
			URL:         "https://github.com",
			Code:        "ghi9012",
			UserID:      "2f6c9d14-9e42-4c31-9e6a-0f8c4b2a7c55",
		},
		{
			Base: model.Base{
				ID: "d4e5f6a7-b8c9-4d0e-1f2a-3b4c5d6e7f8a",
			},
			Description: "YouTube - Video Platform",
			URL:         "https://youtube.com",
			Code:        "jkl3456",
			UserID:      "2f6c9d14-9e42-4c31-9e6a-0f8c4b2a7c55",
		},
		{
			Base: model.Base{
				ID: "e5f6a7b8-c9d0-4e1f-2a3b-4c5d6e7f8a9b",
			},
			Description: "Stack Overflow - Q&A for Developers",
			URL:         "https://stackoverflow.com",
			Code:        "mno7890",
			UserID:      "e3c2a8f1-1d3b-4c62-8e54-6b7f9a2d1c90",
		},
		{
			Base: model.Base{
				ID: "f6a7b8c9-d0e1-4f2a-3b4c-5d6e7f8a9b0c",
			},
			Description: "Golang - Programming Language",
			URL:         "https://go.dev",
			Code:        "pqr1234",
			UserID:      "e3c2a8f1-1d3b-4c62-8e54-6b7f9a2d1c90",
		},
		{
			Base: model.Base{
				ID: "c1d2e3f4-a5b6-4c7d-8e9f-0a1b2c3d4e5f",
			},
			Description: "Personal Blog",
			URL:         "https://johndoe.example.com",
			Code:        "stu5678",
			UserID:      "550e8400-e29b-41d4-a716-446655440000",
		},
		{
			Base: model.Base{
				ID: "d2e3f4a5-b6c7-4d8e-9f0a-1b2c3d4e5f6a",
			},
			Description: "LinkedIn Profile",
			URL:         "https://linkedin.com/in/johndoe",
			Code:        "vwx9012",
			UserID:      "550e8400-e29b-41d4-a716-446655440000",
		},
	}

	return db.CreateInBatches(bookmarks, len(bookmarks)).Error
}


