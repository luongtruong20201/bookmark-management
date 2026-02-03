package infrastructure

import (
	"github.com/luongtruong20201/bookmark-management/pkg/common"
	"github.com/luongtruong20201/bookmark-management/pkg/jwt"
)

func CreateJWTProvider() (jwt.JWTGenerator, jwt.JWTValidator) {
	jwtGenerator, err := jwt.NewJWTGenerator("./private.pem")
	common.HandleError(err)

	jwtValidator, err := jwt.NewJWTValidator("./public.pem")
	common.HandleError(err)

	return jwtGenerator, jwtValidator
}
