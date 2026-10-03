package ovh

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/ovh/go-ovh/ovh"
)

func TestDatabaseNoMatchingComputeError(t *testing.T) {
	params := &CloudProjectDatabaseCreateOpts{
		Plan:    "business",
		Version: "15",
		NodesPattern: CloudProjectDatabaseNodesPattern{
			Flavor: "db1-7",
			Number: 2,
			Region: "GRA",
		},
	}

	noMatchingComputeErr := &ovh.APIError{
		Code:    404,
		Class:   errClassNoMatchingCompute,
		Message: "no compute matches your expectations, please refer to the availability endpoint",
	}

	cases := []struct {
		name        string
		err         error
		expectError bool
		contains    []string
	}{
		{
			name:        "no matching compute API error is rewritten",
			err:         noMatchingComputeErr,
			expectError: true,
			contains: []string{
				`engine "postgresql"`,
				`version "15"`,
				`plan "business"`,
				`flavor "db1-7"`,
				`region "GRA"`,
				"ovh_cloud_project_database_capabilities",
				"no compute matches your expectations",
			},
		},
		{
			name:        "wrapped no matching compute API error is rewritten",
			err:         fmt.Errorf("calling Post: %w", noMatchingComputeErr),
			expectError: true,
			contains:    []string{`flavor "db1-7"`},
		},
		{
			name: "API error of another class is left alone",
			err: &ovh.APIError{
				Code:    404,
				Class:   "Client::NotFound",
				Message: "service not found",
			},
			expectError: false,
		},
		{
			name:        "non API error is left alone",
			err:         errors.New("connection refused"),
			expectError: false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := databaseNoMatchingComputeError(tc.err, "postgresql", params)
			if tc.expectError && got == nil {
				t.Fatalf("expected an error, got nil")
			}
			if !tc.expectError && got != nil {
				t.Fatalf("expected nil, got %v", got)
			}
			if got == nil {
				return
			}
			for _, want := range tc.contains {
				if !strings.Contains(got.Error(), want) {
					t.Errorf("error %q does not contain %q", got.Error(), want)
				}
			}
			if tc.err == noMatchingComputeErr && !errors.Is(got, noMatchingComputeErr) {
				t.Errorf("expected the rewritten error to wrap the API error, got %v", got)
			}
		})
	}
}
