package ovh

import (
	"errors"
	"fmt"
	"testing"

	"github.com/ovh/go-ovh/ovh"
)

func TestIsCloudLoadbalancerBusy(t *testing.T) {
	cases := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name: "parent load balancer not active",
			err: &ovh.APIError{
				Code:    409,
				Class:   "Client::Conflict::ResourceInInvalidState",
				Message: "Failed to modify listener. The parent load balancer is not in ACTIVE state. Wait for ongoing operations to complete and retry. (detail: OpenStack error (HTTP 409))",
			},
			expected: true,
		},
		{
			name: "conflict of another class",
			err: &ovh.APIError{
				Code:    409,
				Class:   "Client::Conflict::Duplicate",
				Message: "Failed to create listener. Protocol port 80 is already used.",
			},
			expected: false,
		},
		{
			name: "same class but permanent conflict",
			err: &ovh.APIError{
				Code:    409,
				Class:   "Client::Conflict::ResourceInInvalidState",
				Message: "Failed to modify listener. The pool is attached to another listener.",
			},
			expected: false,
		},
		{
			name: "bad request",
			err: &ovh.APIError{
				Code:    400,
				Class:   "Client::BadRequest::InvalidParameter",
				Message: "Failed to create listener. TLS container is required for TERMINATED_HTTPS.",
			},
			expected: false,
		},
		{
			name:     "not an API error",
			err:      errors.New("connection reset by peer"),
			expected: false,
		},
		{
			name: "wrapped API error",
			err: fmt.Errorf("calling Post: %w", &ovh.APIError{
				Code:    409,
				Class:   "Client::Conflict::ResourceInInvalidState",
				Message: "The parent load balancer is not in ACTIVE state.",
			}),
			expected: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := isCloudLoadbalancerBusy(tc.err); got != tc.expected {
				t.Errorf("isCloudLoadbalancerBusy() = %v, expected %v", got, tc.expected)
			}
		})
	}
}
