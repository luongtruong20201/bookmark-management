package password

import "github.com/luongtruong20201/bookmark-management/pkg/stringutils"

const (
	passLength = 10
)

// passwordService implements the Password interface and provides business logic
// for password generation. It uses a key generator to create random alphanumeric passwords.
type passwordService struct{}

// Password interface represents password service
//
//go:generate mockery --name Password --filename password_service.go
type Password interface {
	GeneratePassword() (string, error)
}

// NewPassword return a new instance of the password service
func NewPassword() Password {
	return &passwordService{}
}

// GeneratePassword generates a random password of the specified length using
// alphanumeric characters. Returns an error if random number generation fails.
func (s *passwordService) GeneratePassword() (string, error) {
	return stringutils.GenerateCode(passLength)
}
