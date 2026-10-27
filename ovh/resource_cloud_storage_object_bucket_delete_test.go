package ovh

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/ovh/go-ovh/ovh"
	"github.com/ovh/terraform-provider-ovh/v2/ovh/ovhwrap"
	"go.uber.org/ratelimit"
)

// newCloudStorageObjectBucketDeleteTest returns the bucket resource talking to a fake API serving
// one bucket with the given objects, its delete state, and the list of the calls received.
// After a DELETE, the bucket is gone, or stays DELETING with a failed task when deleteFails is set.
func newCloudStorageObjectBucketDeleteTest(t *testing.T, objects string, forceDestroy, deleteFails bool) (*cloudStorageObjectBucketResource, tfsdk.State, *[]string) {
	t.Helper()
	ctx := context.Background()
	bucketPath := "/v2/publicCloud/project/xxx/storage/object/bucket/GRA_my-bucket"

	var calls []string
	deleted := false
	mux := http.NewServeMux()
	mux.HandleFunc("/1.0/auth/time", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, time.Now().Unix())
	})
	mux.HandleFunc(bucketPath+"/object", func(w http.ResponseWriter, r *http.Request) {
		calls = append(calls, r.Method+" object?"+r.URL.RawQuery)
		fmt.Fprint(w, objects)
	})
	mux.HandleFunc(bucketPath+"/bulkDeleteObjects", func(w http.ResponseWriter, r *http.Request) {
		calls = append(calls, r.Method+" bulkDeleteObjects")
		objects = "[]"
		fmt.Fprint(w, `{"deleted": [], "errors": []}`)
	})
	mux.HandleFunc(bucketPath, func(w http.ResponseWriter, r *http.Request) {
		calls = append(calls, r.Method+" bucket")
		switch {
		case r.Method == http.MethodDelete:
			deleted = true
		case deleted && deleteFails:
			fmt.Fprint(w, `{"id": "GRA_my-bucket", "resourceStatus": "DELETING", "currentTasks": [{"type": "DELETE_S3_BUCKET", "status": "ERROR", "errors": [{"message": "The bucket is not empty"}]}]}`)
		case deleted:
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, `{"message": "S3 bucket not found"}`)
		}
	})
	server := httptest.NewServer(mux)
	t.Cleanup(server.Close)

	client, err := ovh.NewClient(server.URL+"/1.0", "key", "secret", "consumer")
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	r := &cloudStorageObjectBucketResource{config: &Config{OVHClient: ovhwrap.NewClient(client, ratelimit.NewUnlimited())}}

	var schemaResp resource.SchemaResponse
	r.Schema(ctx, resource.SchemaRequest{}, &schemaResp)
	state := tfsdk.State{Schema: schemaResp.Schema, Raw: tftypes.NewValue(schemaResp.Schema.Type().TerraformType(ctx), nil)}
	state.SetAttribute(ctx, path.Root("service_name"), "xxx")
	state.SetAttribute(ctx, path.Root("id"), "GRA_my-bucket")
	if forceDestroy {
		state.SetAttribute(ctx, path.Root("force_destroy"), true)
	}
	return r, state, &calls
}

func TestCloudStorageObjectBucketDelete(t *testing.T) {
	const twoVersions = `[{"key": "a.txt", "versionId": "2"}, {"key": "a.txt", "versionId": "1"}]`

	for name, tc := range map[string]struct {
		objects      string
		forceDestroy bool
		deleteFails  bool
		wantError    string
		wantCalls    []string
	}{
		"empty bucket": {
			objects:   "[]",
			wantCalls: []string{"GET object?withVersions=true", "DELETE bucket", "GET bucket"},
		},
		// Never send the DELETE: the API would delete the bucket later, once emptied
		"not empty refused": {
			objects:   twoVersions,
			wantError: "Bucket GRA_my-bucket is not empty",
			wantCalls: []string{"GET object?withVersions=true"},
		},
		"not empty with force_destroy": {
			objects:      twoVersions,
			forceDestroy: true,
			wantCalls:    []string{"GET object?withVersions=true", "POST bulkDeleteObjects", "GET object?withVersions=true", "DELETE bucket", "GET bucket"},
		},
		"failed deletion task": {
			objects:     "[]",
			deleteFails: true,
			wantError:   "Error waiting for bucket to be deleted",
			wantCalls:   []string{"GET object?withVersions=true", "DELETE bucket", "GET bucket"},
		},
	} {
		t.Run(name, func(t *testing.T) {
			r, state, calls := newCloudStorageObjectBucketDeleteTest(t, tc.objects, tc.forceDestroy, tc.deleteFails)

			resp := resource.DeleteResponse{State: state}
			r.Delete(context.Background(), resource.DeleteRequest{State: state}, &resp)

			if tc.wantError == "" && resp.Diagnostics.HasError() {
				t.Fatalf("unexpected diagnostics: %v", resp.Diagnostics)
			}
			if tc.wantError != "" && (!resp.Diagnostics.HasError() || resp.Diagnostics.Errors()[0].Summary() != tc.wantError) {
				t.Fatalf("expected error %q, got %v", tc.wantError, resp.Diagnostics)
			}
			if strings.Join(*calls, ", ") != strings.Join(tc.wantCalls, ", ") {
				t.Errorf("unexpected calls:\n got: %v\nwant: %v", *calls, tc.wantCalls)
			}
		})
	}
}
