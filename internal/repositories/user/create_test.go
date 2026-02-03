package user

import (
	"context"
	"strings"
	"testing"

	model "github.com/luongtruong20201/bookmark-management/internal/models"
	"github.com/luongtruong20201/bookmark-management/internal/test/fixture"
	"github.com/luongtruong20201/bookmark-management/pkg/dbutils"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestUser_CreateUser(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name           string
		setupDB        func(t *testing.T) *gorm.DB
		inputUser      *model.User
		expectedError  error
		expectedOutput *model.User
		verifyFunc     func(t *testing.T, db *gorm.DB, user *model.User)
	}{
		{
			name: "success - create new user",
			setupDB: func(t *testing.T) *gorm.DB {
				return fixture.NewFixture(t, &fixture.UserCommonTestDB{})
			},
			inputUser: &model.User{
				Base: model.Base{
					ID: "9b5c1e3e-7c3b-4f4e-8e7c-6e7a2f5d3a92",
				},
				DisplayName: "Nguyen Van An",
				Username:    "an.nguyen1",
				Password:    "P@ssw0rd1",
				Email:       "an.nguyen1@example.com",
			},
			expectedError: nil,
			expectedOutput: &model.User{
				Base: model.Base{
					ID: "9b5c1e3e-7c3b-4f4e-8e7c-6e7a2f5d3a92",
				},
				DisplayName: "Nguyen Van An",
				Username:    "an.nguyen1",
				Password:    "P@ssw0rd1",
				Email:       "an.nguyen1@example.com",
			},
			verifyFunc: func(t *testing.T, db *gorm.DB, user *model.User) {
				toCheckUser := &model.User{}
				err := db.Where("username = ?", user.Username).First(toCheckUser).Error
				assert.Nil(t, err)
				assert.Equal(t, user.ID, toCheckUser.ID)
				assert.Equal(t, user.Username, toCheckUser.Username)
				assert.Equal(t, user.Email, toCheckUser.Email)
				assert.Equal(t, user.DisplayName, toCheckUser.DisplayName)
			},
		},
		{
			name: "success - create user without ID (auto-generate UUID)",
			setupDB: func(t *testing.T) *gorm.DB {
				return fixture.NewFixture(t, &fixture.UserCommonTestDB{})
			},
			inputUser: &model.User{
				DisplayName: "Test User",
				Username:    "test.user",
				Password:    "P@ssw0rd1",
				Email:       "test.user@example.com",
			},
			expectedError:  nil,
			expectedOutput: nil,
			verifyFunc: func(t *testing.T, db *gorm.DB, user *model.User) {
				toCheckUser := &model.User{}
				err := db.Where("username = ?", user.Username).First(toCheckUser).Error
				assert.Nil(t, err)
				assert.NotEmpty(t, toCheckUser.ID)
				assert.Equal(t, user.Username, toCheckUser.Username)
				assert.Equal(t, user.Email, toCheckUser.Email)
				assert.Equal(t, user.DisplayName, toCheckUser.DisplayName)
			},
		},
		{
			name: "error - duplicate username",
			setupDB: func(t *testing.T) *gorm.DB {
				return fixture.NewFixture(t, &fixture.UserCommonTestDB{})
			},
			inputUser: &model.User{
				Base: model.Base{
					ID: "9b5c1e3e-7c3b-4f4e-8e7c-6e7a2f5d3a99",
				},
				DisplayName: "Duplicate User",
				Username:    "an.nguyen",
				Password:    "P@ssw0rd1",
				Email:       "duplicate@example.com",
			},
			expectedError:  dbutils.ErrDuplicationType,
			expectedOutput: nil,
			verifyFunc:     nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()
			db := tc.setupDB(t)
			repo := NewUser(db)

			res, err := repo.CreateUser(ctx, tc.inputUser)

			if tc.expectedError != nil {
				assert.Error(t, err)

				if err != tc.expectedError {

					errStr := strings.ToLower(err.Error())
					assert.True(t, strings.Contains(errStr, "unique constraint") || err == dbutils.ErrDuplicationType,
						"expected duplicate constraint error, got: %v", err)
				}
				assert.Nil(t, res)
			} else {
				assert.NoError(t, err)
				if tc.expectedOutput != nil {
					assert.Equal(t, tc.expectedOutput.ID, res.ID)
					assert.Equal(t, tc.expectedOutput.Username, res.Username)
					assert.Equal(t, tc.expectedOutput.Email, res.Email)
					assert.Equal(t, tc.expectedOutput.DisplayName, res.DisplayName)
				}
				if tc.verifyFunc != nil {
					tc.verifyFunc(t, db, res)
				}
			}
		})
	}
}
