package healthcheck

import "context"

// Ping sends a PING command to the Redis server to verify connectivity.
//
// This method is a wrapper around the Redis client's Ping method, which sends
// a lightweight PING command to the Redis server. The PING command is designed
// to be fast and non-blocking, making it ideal for health check operations.
//
// The method will return an error if:
//   - The Redis connection is closed or not established
//   - A network error occurs while communicating with Redis
//   - The context is cancelled or times out before the operation completes
//   - The Redis server is unreachable or not responding
//
// Parameters:
//   - ctx: Context for controlling the request lifecycle. The context can be used
//     to set timeouts or cancel the operation. If the context is cancelled,
//     the method will return immediately with a context error.
//
// Returns:
//   - error: Returns nil if the PING command succeeds, indicating that Redis
//     is reachable and responsive. Returns a non-nil error if:
//   - redis.ErrClosed: The Redis client connection is closed
//   - context.DeadlineExceeded: The operation exceeded the context timeout
//   - context.Canceled: The context was cancelled
//   - Any network or Redis-specific error
//
// Thread safety:
//
//	This method is safe to call concurrently from multiple goroutines, as the
//	underlying Redis client is designed to be thread-safe.
//
// Example usage:
//
//	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
//	defer cancel()
//
//	if err := healthCheck.Ping(ctx); err != nil {
//	    // Handle health check failure
//	    return fmt.Errorf("Redis health check failed: %w", err)
//	}
//	// Redis is healthy and reachable
func (r *redisHealthCheck) Ping(ctx context.Context) error {
	return r.redis.Ping(ctx).Err()
}
