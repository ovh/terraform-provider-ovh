package helpers

import (
	"errors"
	"fmt"
	"testing"

	"github.com/ovh/go-ovh/ovh"
)

func TestIsNotFound(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{"nil", nil, false},
		{"404", &ovh.APIError{Code: 404}, true},
		{"wrapped 404", fmt.Errorf("calling Get: %w", error(&ovh.APIError{Code: 404})), true},
		{"403", &ovh.APIError{Code: 403}, false},
		{"other error", errors.New("timeout"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsNotFound(tt.err); got != tt.expected {
				t.Fatalf("expected %v, got %v", tt.expected, got)
			}
		})
	}
}
