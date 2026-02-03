package healthcheck

import (
	"testing"

	repository "github.com/luongtruong20201/bookmark-management/internal/repositories/healthcheck"
	redisPkg "github.com/luongtruong20201/bookmark-management/pkg/redis"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func TestHealthCheckService_Check(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name            string
		serviceName     string
		instanceID      string
		expectedMsg     string
		expectedSvcName string
		setupRedis      func(t *testing.T) *redis.Client
	}{
		{
			name:        "success",
			serviceName: "bookmark_service",
			instanceID:  "instance_id",
			expectedMsg: "OK",
			setupRedis: func(t *testing.T) *redis.Client {
				redis := redisPkg.InitMockRedis(t)
				return redis
			},
		},
		{
			name:        "fail",
			serviceName: "bookmark_service",
			instanceID:  "instance_id",
			expectedMsg: "NOT_OK",
			setupRedis: func(t *testing.T) *redis.Client {
				redis := redisPkg.InitMockRedis(t)
				_ = redis.Close()
				return redis
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctx := t.Context()
			redis := tc.setupRedis(t)
			repository := repository.NewHealthCheck(redis)
			testSvc := NewHealthcheck(tc.serviceName, tc.instanceID, repository)
			message, serviceName, instanceId := testSvc.Check(ctx)

			assert.Equal(t, tc.instanceID, instanceId)
			assert.Equal(t, message, tc.expectedMsg)
			assert.Equal(t, serviceName, tc.serviceName)
		})
	}
}
