package ovh

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/ovh/go-ovh/ovh"
	"github.com/ovh/terraform-provider-ovh/v2/ovh/ovhwrap"
	"go.uber.org/ratelimit"
)

// Availabilities in the format returned by GET /cloud/project/{serviceName}/database/availability
const testCloudProjectDatabaseAvailabilities = `[
	{"engine": "postgresql", "version": "14", "plan": "business", "region": "GRA", "lifecycle": {"status": "END_OF_SALE"},
	 "specifications": {"flavor": "db1-4", "network": "public", "nodes": {"minimum": 2, "maximum": 2},
	 "storage": {"minimum": {"unit": "GB", "value": 80}, "maximum": {"unit": "GB", "value": 240}, "step": {"unit": "GB", "value": 10}}}},
	{"engine": "postgresql", "version": "16", "plan": "business", "region": "GRA", "lifecycle": {"status": "STABLE"},
	 "specifications": {"flavor": "db1-4", "network": "public", "nodes": {"minimum": 2, "maximum": 2},
	 "storage": {"minimum": {"unit": "GB", "value": 80}, "maximum": {"unit": "GB", "value": 240}, "step": {"unit": "GB", "value": 10}}}},
	{"engine": "postgresql", "version": "16", "plan": "business", "region": "GRA", "lifecycle": {"status": "STABLE"},
	 "specifications": {"flavor": "db1-7", "network": "public", "nodes": {"minimum": 2, "maximum": 2},
	 "storage": {"minimum": {"unit": "GB", "value": 160}, "maximum": {"unit": "GB", "value": 480}, "step": {"unit": "GB", "value": 10}}}},
	{"engine": "postgresql", "version": "16", "plan": "business", "region": "SBG", "lifecycle": {"status": "STABLE"},
	 "specifications": {"flavor": "db1-4", "network": "private", "nodes": {"minimum": 2, "maximum": 2},
	 "storage": {"minimum": {"unit": "GB", "value": 80}, "maximum": {"unit": "GB", "value": 240}, "step": {"unit": "GB", "value": 10}}}},
	{"engine": "postgresql", "version": "17", "plan": "production", "region": "GRA", "lifecycle": {"status": "STABLE"},
	 "specifications": {"flavor": "db1-4", "network": "public", "nodes": {"minimum": 3, "maximum": 3},
	 "storage": {"minimum": {"unit": "GB", "value": 80}, "maximum": {"unit": "GB", "value": 240}, "step": {"unit": "GB", "value": 10}}}},
	{"engine": "valkey", "version": "8.0", "plan": "essential", "region": "GRA", "lifecycle": {"status": "STABLE"},
	 "specifications": {"flavor": "db1-4", "network": "public", "nodes": {"minimum": 1, "maximum": 1}, "storage": null}},
	{"engine": "mysql", "version": "8", "plan": "essential", "region": "GRA", "lifecycle": {"status": "STABLE"},
	 "specifications": {"flavor": "db1-4", "network": "public", "nodes": {"minimum": 1, "maximum": 1},
	 "storage": {"minimum": {"unit": "MB", "value": 1}, "maximum": {"unit": "MB", "value": 2}, "step": null}}}
]`

func testCloudProjectDatabaseCreateOpts(version, plan, region, flavor string, nodes, diskSize int) *CloudProjectDatabaseCreateOpts {
	return &CloudProjectDatabaseCreateOpts{
		Version: version,
		Plan:    plan,
		Disk:    CloudProjectDatabaseDisk{Size: diskSize},
		NodesPattern: CloudProjectDatabaseNodesPattern{
			Flavor: flavor,
			Number: nodes,
			Region: region,
		},
	}
}

func TestExplainNoMatchingCompute(t *testing.T) {
	var availabilities []CloudProjectDatabaseAvailability
	if err := json.Unmarshal([]byte(testCloudProjectDatabaseAvailabilities), &availabilities); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	privateNetwork := testCloudProjectDatabaseCreateOpts("16", "business", "GRA", "db1-4", 2, 0)
	privateNetwork.NetworkID = "network-id"

	tests := []struct {
		name   string
		engine string
		opts   *CloudProjectDatabaseCreateOpts
		want   string
	}{
		{
			name:   "matching configuration",
			engine: "postgresql",
			opts:   testCloudProjectDatabaseCreateOpts("16", "business", "GRA", "db1-7", 2, 200),
			want:   "",
		},
		{
			name:   "unknown engine",
			engine: "cassandra",
			opts:   testCloudProjectDatabaseCreateOpts("4", "business", "GRA", "db1-4", 2, 0),
			want:   `engine "cassandra" is not available (available: mysql, postgresql, valkey)`,
		},
		{
			name:   "version no longer sold",
			engine: "postgresql",
			opts:   testCloudProjectDatabaseCreateOpts("14", "business", "GRA", "db1-4", 2, 0),
			want:   `version "14" is END_OF_SALE for postgresql (available: 16, 17)`,
		},
		{
			name:   "plan not available for version",
			engine: "postgresql",
			opts:   testCloudProjectDatabaseCreateOpts("17", "business", "GRA", "db1-4", 2, 0),
			want:   `plan "business" is not available for postgresql / 17 (available: production)`,
		},
		{
			name:   "region compared without case",
			engine: "postgresql",
			opts:   testCloudProjectDatabaseCreateOpts("16", "business", "gra", "db1-4", 2, 0),
			want:   "",
		},
		{
			name:   "region not available",
			engine: "postgresql",
			opts:   testCloudProjectDatabaseCreateOpts("16", "business", "BHS", "db1-4", 2, 0),
			want:   `region "BHS" is not available for postgresql / 16 / business (available: GRA, SBG)`,
		},
		{
			name:   "flavor not available",
			engine: "postgresql",
			opts:   testCloudProjectDatabaseCreateOpts("16", "business", "GRA", "b3-8", 2, 0),
			want:   `flavor "b3-8" is not available for postgresql / 16 / business / GRA (available: db1-4, db1-7)`,
		},
		{
			name:   "private network not available",
			engine: "postgresql",
			opts:   privateNetwork,
			want:   `private network is not available for postgresql / 16 / business / GRA / db1-4 (available: public network)`,
		},
		{
			name:   "number of nodes not available",
			engine: "postgresql",
			opts:   testCloudProjectDatabaseCreateOpts("16", "business", "GRA", "db1-4", 3, 0),
			want:   `node count 3 is not available for postgresql / 16 / business / GRA / db1-4 / public network (available: 2 nodes)`,
		},
		{
			name:   "disk size too big",
			engine: "postgresql",
			opts:   testCloudProjectDatabaseCreateOpts("16", "business", "GRA", "db1-4", 2, 300),
			want:   `disk size 300 GB is not available for postgresql / 16 / business / GRA / db1-4 / public network (available: 80 to 240 GB in steps of 10 GB)`,
		},
		{
			name:   "disk size not a step",
			engine: "postgresql",
			opts:   testCloudProjectDatabaseCreateOpts("16", "business", "GRA", "db1-4", 2, 85),
			want:   `disk size 85 GB is not available for postgresql / 16 / business / GRA / db1-4 / public network (available: 80 to 240 GB in steps of 10 GB)`,
		},
		{
			name:   "disk size ignored without storage",
			engine: "valkey",
			opts:   testCloudProjectDatabaseCreateOpts("8.0", "essential", "GRA", "db1-4", 1, 50),
			want:   "",
		},
		{
			name:   "disk size ignored when not in GB",
			engine: "mysql",
			opts:   testCloudProjectDatabaseCreateOpts("8", "essential", "GRA", "db1-4", 1, 50),
			want:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := explainNoMatchingCompute(availabilities, tt.engine, tt.opts); got != tt.want {
				t.Errorf("got:\n\t%s\nwant:\n\t%s", got, tt.want)
			}
		})
	}
}

func TestCloudProjectDatabaseExplainCreateError(t *testing.T) {
	ctx := context.Background()
	opts := testCloudProjectDatabaseCreateOpts("16", "business", "GRA", "b3-8", 2, 0)
	noMatchingCompute := &ovh.APIError{Code: http.StatusNotFound, Class: errClassNoMatchingCompute}

	newConfig := func(t *testing.T, calls *int32, status int, body string) *Config {
		mux := http.NewServeMux()
		mux.HandleFunc("/auth/time", func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, time.Now().Unix())
		})
		mux.HandleFunc("/cloud/project/project-id/database/availability", func(w http.ResponseWriter, r *http.Request) {
			atomic.AddInt32(calls, 1)
			w.WriteHeader(status)
			fmt.Fprint(w, body)
		})

		server := httptest.NewServer(mux)
		t.Cleanup(server.Close)

		client, err := ovh.NewClient(server.URL, "key", "secret", "consumer")
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}
		return &Config{OVHClient: ovhwrap.NewClient(client, ratelimit.NewUnlimited())}
	}

	t.Run("no matching compute", func(t *testing.T) {
		var calls int32
		config := newConfig(t, &calls, http.StatusOK, testCloudProjectDatabaseAvailabilities)

		err := fmt.Errorf("calling Post: %w", error(noMatchingCompute))
		got := cloudProjectDatabaseExplainCreateError(ctx, config, "project-id", "postgresql", opts, err)
		if !strings.HasPrefix(got, "no compute matches the requested configuration, ") || !strings.Contains(got, `flavor "b3-8" is not available`) {
			t.Errorf("unexpected explanation: %s", got)
		}
		if calls != 1 {
			t.Errorf("expected 1 call to the availability endpoint, got %d", calls)
		}
	})

	t.Run("other API error", func(t *testing.T) {
		var calls int32
		config := newConfig(t, &calls, http.StatusOK, testCloudProjectDatabaseAvailabilities)

		err := &ovh.APIError{Code: http.StatusBadRequest, Class: "Client::BadRequest"}
		if got := cloudProjectDatabaseExplainCreateError(ctx, config, "project-id", "postgresql", opts, err); got != "" {
			t.Errorf("expected no explanation, got: %s", got)
		}
		if calls != 0 {
			t.Errorf("expected no call to the availability endpoint, got %d", calls)
		}
	})

	t.Run("availability endpoint fails", func(t *testing.T) {
		var calls int32
		config := newConfig(t, &calls, http.StatusInternalServerError, `{"message": "Internal Server Error"}`)

		if got := cloudProjectDatabaseExplainCreateError(ctx, config, "project-id", "postgresql", opts, noMatchingCompute); got != "" {
			t.Errorf("expected no explanation, got: %s", got)
		}
	})
}
