package user

import (
	"context"
	"testing"

	model "github.com/luongtruong20201/bookmark-management/internal/models"
	"github.com/luongtruong20201/bookmark-management/internal/test/fixture"
	"github.com/luongtruong20201/bookmark-management/pkg/dbutils"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestUser_UpdateUserProfile(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name          string
		setupDB       func(t *testing.T) *gorm.DB
		id            string
		displayName   string
		email         string
		expectedError error
		expectError   bool
		verifyFunc    func(t *testing.T, db *gorm.DB, id, displayName, email string, got *model.User)
	}{
		{
			name: "success - update existing user profile",
			setupDB: func(t *testing.T) *gorm.DB {
				return fixture.NewFixture(t, &fixture.UserCommonTestDB{})
			},
			id:          "9b5c1e3e-7c3b-4f4e-8e7c-6e7a2f5d3a91",
			displayName: "Updated Display Name",
			email:       "updated.email@example.com",
			verifyFunc: func(t *testing.T, db *gorm.DB, id, displayName, email string, got *model.User) {
				assert.NotNil(t, got)
				assert.Equal(t, id, got.ID)
				assert.Equal(t, displayName, got.DisplayName)
				assert.Equal(t, email, got.Email)

				toCheckUser := &model.User{}
				err := db.Where("id = ?", id).First(toCheckUser).Error
				assert.NoError(t, err)
				assert.Equal(t, displayName, toCheckUser.DisplayName)
				assert.Equal(t, email, toCheckUser.Email)
			},
		},
		{
			name: "error - user not found",
			setupDB: func(t *testing.T) *gorm.DB {
				return fixture.NewFixture(t, &fixture.UserCommonTestDB{})
			},
			id:            "00000000-0000-0000-0000-000000000000",
			displayName:   "Any",
			email:         "any@example.com",
			expectedError: dbutils.ErrNotFoundType,
		},
		{
			name: "error - db error during update",
			setupDB: func(t *testing.T) *gorm.DB {
				db := fixture.NewFixture(t, &fixture.UserCommonTestDB{})
				_ = db.Migrator().DropTable(&model.User{})
				return db
			},
			id:          "9b5c1e3e-7c3b-4f4e-8e7c-6e7a2f5d3a91",
			displayName: "Updated Display Name",
			email:       "updated.email@example.com",
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()
			db := tc.setupDB(t)
			repo := NewUser(db)

			res, err := repo.UpdateUserProfile(ctx, tc.id, tc.displayName, tc.email)

			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedError, err)
				assert.Nil(t, res)
				return
			}

			if tc.expectError {
				assert.Error(t, err)
				assert.Nil(t, res)
				return
			}

			assert.NoError(t, err)
			if tc.verifyFunc != nil {
				tc.verifyFunc(t, db, tc.id, tc.displayName, tc.email, res)
			}
		})
	}
}
