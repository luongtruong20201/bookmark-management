package repository

import (
	"testing"

	model "github.com/luongtruong20201/bookmark-management/internal/models"
	"github.com/luongtruong20201/bookmark-management/internal/test/fixture"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestUser_CreateUser(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name             string
		setupDB          func(t *testing.T) *gorm.DB
		inputUser        *model.User
		expectedErrorStr string
		expectedOutput   *model.User
		verifyFunc       func(db *gorm.DB, user *model.User)
	}{
		{
			name: "normal case",
			setupDB: func(t *testing.T) *gorm.DB {
				return fixture.NewFixture(t, &fixture.UserCommonTestDB{})
			},
			inputUser: &model.User{
				ID:          "9b5c1e3e-7c3b-4f4e-8e7c-6e7a2f5d3a92",
				DisplayName: "Nguyen Van An",
				Username:    "an.nguyen1",
				Password:    "P@ssw0rd1",
				Email:       "an.nguyen1@example.com",
			},
			expectedErrorStr: "",
			expectedOutput: &model.User{
				ID:          "9b5c1e3e-7c3b-4f4e-8e7c-6e7a2f5d3a92",
				DisplayName: "Nguyen Van An",
				Username:    "an.nguyen1",
				Password:    "P@ssw0rd1",
				Email:       "an.nguyen1@example.com",
			},
			verifyFunc: func(db *gorm.DB, user *model.User) {
				toCheckUser := &model.User{}
				err := db.Where("username = ?", user.Username).First(toCheckUser).Error
				assert.Nil(t, err)
				assert.Equal(t, toCheckUser, user)
			},
		},
		{
			name: "duplicate case",
			setupDB: func(t *testing.T) *gorm.DB {
				return fixture.NewFixture(t, &fixture.UserCommonTestDB{})
			},
			inputUser: &model.User{
				ID:          "9b5c1e3e-7c3b-4f4e-8e7c-6e7a2f5d3a91",
				DisplayName: "Nguyen Van An",
				Username:    "an.nguyen",
				Password:    "P@ssw0rd1",
				Email:       "an.nguyen@example.com",
			},
			expectedErrorStr: "",
			expectedOutput:   nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctx := t.Context()
			db := tc.setupDB(t)
			repo := NewUser(db)

			res, err := repo.CreateUser(ctx, tc.inputUser)
			if err != nil {
				assert.ErrorContains(t, err, tc.expectedErrorStr)
			}

			assert.Equal(t, res, tc.expectedOutput)

			if err == nil {
				tc.verifyFunc(db, tc.expectedOutput)
			}
		})
	}
}
