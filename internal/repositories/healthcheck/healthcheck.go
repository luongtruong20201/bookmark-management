package healthcheck

import (
	"context"

	"github.com/redis/go-redis/v9"
)

// HealthCheck defines the interface for health check operations.
// It provides methods to verify the connectivity and health status of external dependencies,
// primarily used for checking Redis connection health.
//
// Implementations of this interface should be lightweight and fast, as they are typically
// called frequently by health check endpoints to monitor service availability.
//
//go:generate mockery --name HealthCheck --filename health_check.go
type HealthCheck interface {
	// Ping checks the connectivity to the underlying data store (e.g., Redis).
	// It sends a PING command to verify that the connection is alive and responsive.
	//
	// Parameters:
	//   - ctx: Context for controlling the request lifecycle. Can be used for cancellation
	//          and timeout control. If the context is cancelled, the operation should
	//          return immediately with an appropriate error.
	//
	// Returns:
	//   - error: Returns nil if the ping operation succeeds, indicating the connection
	//            is healthy. Returns a non-nil error if:
	//            - The connection is closed or unavailable (e.g., redis.ErrClosed)
	//            - A network timeout occurs
	//            - The context is cancelled or times out
	//            - Any other connection-related error occurs
	//
	// Example usage:
	//   ctx := context.Background()
	//   if err := healthCheck.Ping(ctx); err != nil {
	//       log.Printf("Health check failed: %v", err)
	//   }
	Ping(context.Context) error
}

// NewHealthCheck creates a new HealthCheck repository instance that uses Redis
// as the underlying health check mechanism.
//
// This function initializes a redisHealthCheck implementation that wraps a Redis client.
// The returned HealthCheck instance can be used to verify Redis connectivity by sending
// PING commands through the provided Redis client.
//
// Parameters:
//   - redis: A pointer to a configured Redis client. The client should be properly
//     initialized and connected to a Redis server. Must not be nil.
//
// Returns:
//   - HealthCheck: A new instance of HealthCheck that uses the provided Redis client
//     for health check operations. The implementation is thread-safe and
//     can be used concurrently from multiple goroutines.
//
// Example usage:
//
//	redisClient := redis.NewClient(&redis.Options{
//	    Addr: "localhost:6379",
//	})
//	healthCheck := NewHealthCheck(redisClient)
//	err := healthCheck.Ping(context.Background())
func NewHealthCheck(redis *redis.Client) HealthCheck {
	return &redisHealthCheck{
		redis: redis,
	}
}

// redisHealthCheck is the concrete implementation of HealthCheck interface
// that uses Redis for health check operations.
//
// This struct wraps a Redis client and implements the HealthCheck interface
// by delegating ping operations to the underlying Redis connection. It provides
// a simple way to check if the Redis connection is alive and responsive.
//
// Fields:
//   - redis: A pointer to the Redis client used for health check operations.
//     This client is used to send PING commands to verify connectivity.
type redisHealthCheck struct {
	redis *redis.Client
}
