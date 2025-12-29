package service

import (
	"testing"

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
	}{
		{
			name:        "success",
			serviceName: "bookmark_service",
			instanceID:  "instance_id",
			expectedMsg: "OK",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			testSvc := NewHealthcheck(tc.serviceName, tc.instanceID)
			message, serviceName, instanceId := testSvc.Check()

			assert.Equal(t, serviceName, tc.serviceName)
			assert.Equal(t, instanceId, instanceId)
			assert.Equal(t, message, "OK")
		})
	}
}
