package main

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	model "github.com/luongtruong20201/bookmark-management/internal/models"
	repository "github.com/luongtruong20201/bookmark-management/internal/repositories"
	sqldb "github.com/luongtruong20201/bookmark-management/pkg/sql"
)

func main() {
	db, err := sqldb.NewClient("")
	if err != nil {
		panic(err)
	}

	db.AutoMigrate(&model.User{})

	userRepo := repository.NewUser(db)

	user, err := userRepo.CreateUser(context.Background(), &model.User{
		ID:          uuid.New().String(),
		Username:    "truonglq",
		Password:    "truonglq",
		DisplayName: "truonglq",
		Email:       "truonglq@gmail.com",
	})

	if err != nil {
		panic(err)
	}

	fmt.Println("check user: ", user)
}
