package user

import (
	"github.com/gin-gonic/gin"
	service "github.com/luongtruong20201/bookmark-management/internal/services/user"
)

// User defines the interface for user handlers.
// It provides methods to handle HTTP requests for user-related operations including
// registration, authentication, profile retrieval, and profile updates.
//
//go:generate mockery --name=Service --filename=user.go --output=./mocks
type User interface {
	// RegisterUser handles user registration requests.
	// It validates input, creates a new user account, and returns the created user information.
	RegisterUser(c *gin.Context)
	// Login handles user authentication requests.
	// It validates credentials and returns a JWT token upon successful authentication.
	Login(c *gin.Context)
	// GetProfile retrieves the profile information of the currently authenticated user.
	// The user ID is extracted from the JWT token claims in the request context.
	GetProfile(c *gin.Context)
	// UpdateProfile updates the profile information (display name and email) of the currently authenticated user.
	// The user ID is extracted from the JWT token claims in the request context.
	UpdateProfile(c *gin.Context)
}

// user implements the User interface and provides HTTP handlers for user operations.
// It encapsulates the user service dependency for business logic execution.
type user struct {
	svc service.User
}

// NewUser creates a new user handler with the provided user service.
func NewUser(svc service.User) User {
	return &user{
		svc: svc,
	}
}
