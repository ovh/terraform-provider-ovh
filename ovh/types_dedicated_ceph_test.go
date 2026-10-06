package ovh

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/ovh/go-ovh/ovh"
	"github.com/ovh/terraform-provider-ovh/v2/ovh/ovhwrap"
	ovhtypes "github.com/ovh/terraform-provider-ovh/v2/ovh/types"
	"go.uber.org/ratelimit"
)

func dedicatedCephUserPoolPermissionValue(poolName string, read, write, execute, classRead, classWrite bool) attr.Value {
	return types.ObjectValueMust(DedicatedCephUserPoolPermissionAttrTypes(), map[string]attr.Value{
		"pool_name":   types.StringValue(poolName),
		"read":        types.BoolValue(read),
		"write":       types.BoolValue(write),
		"execute":     types.BoolValue(execute),
		"class_read":  types.BoolValue(classRead),
		"class_write": types.BoolValue(classWrite),
	})
}

func dedicatedCephUserPoolPermissionsList(elements ...attr.Value) types.List {
	return types.ListValueMust(types.ObjectType{AttrTypes: DedicatedCephUserPoolPermissionAttrTypes()}, elements)
}

func dedicatedCephUserPoolPermissionsNull() types.List {
	return types.ListNull(types.ObjectType{AttrTypes: DedicatedCephUserPoolPermissionAttrTypes()})
}

// The attr types of pool_permissions are declared apart from the schema, so a nested attribute
// added to one side only produces a runtime error rather than a compile error.
func TestDedicatedCephUserSchemaMatchesAttrTypes(t *testing.T) {
	var resp resource.SchemaResponse
	(&dedicatedCephUserResource{}).Schema(context.Background(), resource.SchemaRequest{}, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("schema diagnostics: %v", resp.Diagnostics)
	}

	want := types.ListType{ElemType: types.ObjectType{AttrTypes: DedicatedCephUserPoolPermissionAttrTypes()}}
	if got := resp.Schema.Attributes["pool_permissions"].GetType(); !got.Equal(want) {
		t.Errorf("pool_permissions type mismatch:\n got: %v\nwant: %v", got, want)
	}
}

func TestDedicatedCephUserPoolPermissionIsEmpty(t *testing.T) {
	if !(DedicatedCephUserPoolPermission{PoolName: "pool"}).IsEmpty() {
		t.Error("expected a permission without any flag to be empty")
	}

	for name, permission := range map[string]DedicatedCephUserPoolPermission{
		"read":        {Read: true},
		"write":       {Write: true},
		"execute":     {Execute: true},
		"class_read":  {ClassRead: true},
		"class_write": {ClassWrite: true},
	} {
		if permission.IsEmpty() {
			t.Errorf("expected a permission with %s to not be empty", name)
		}
	}
}

func TestDedicatedCephUserPoolPermissionGrantsNothing(t *testing.T) {
	for name, tt := range map[string]struct {
		model DedicatedCephUserPoolPermissionModel
		want  bool
	}{
		"flags omitted": {DedicatedCephUserPoolPermissionModel{
			PoolName: types.StringValue("pool"), Read: types.BoolNull(), Write: types.BoolNull(),
			Execute: types.BoolNull(), ClassRead: types.BoolNull(), ClassWrite: types.BoolNull(),
		}, true},
		"flags false": {DedicatedCephUserPoolPermissionModel{
			PoolName: types.StringValue("pool"), Read: types.BoolValue(false), Write: types.BoolValue(false),
			Execute: types.BoolValue(false), ClassRead: types.BoolValue(false), ClassWrite: types.BoolValue(false),
		}, true},
		"one flag true": {DedicatedCephUserPoolPermissionModel{
			PoolName: types.StringValue("pool"), Read: types.BoolNull(), Write: types.BoolNull(),
			Execute: types.BoolNull(), ClassRead: types.BoolNull(), ClassWrite: types.BoolValue(true),
		}, false},
		"one flag unknown": {DedicatedCephUserPoolPermissionModel{
			PoolName: types.StringValue("pool"), Read: types.BoolUnknown(), Write: types.BoolNull(),
			Execute: types.BoolNull(), ClassRead: types.BoolNull(), ClassWrite: types.BoolNull(),
		}, false},
	} {
		if got := tt.model.GrantsNothing(); got != tt.want {
			t.Errorf("%s: expected GrantsNothing to be %t, got %t", name, tt.want, got)
		}
	}
}

func TestDedicatedCephUserToPoolPermissions(t *testing.T) {
	ctx := context.Background()

	model := DedicatedCephUserModel{PoolPermissions: dedicatedCephUserPoolPermissionsNull()}
	permissions, diags := model.ToPoolPermissions(ctx)
	if diags.HasError() {
		t.Fatalf("unexpected diagnostics: %v", diags)
	}
	if len(permissions) != 0 {
		t.Errorf("expected no permission from a null list, got %v", permissions)
	}

	model.PoolPermissions = dedicatedCephUserPoolPermissionsList(
		dedicatedCephUserPoolPermissionValue("pool", true, false, true, false, true),
	)
	permissions, diags = model.ToPoolPermissions(ctx)
	if diags.HasError() {
		t.Fatalf("unexpected diagnostics: %v", diags)
	}

	want := DedicatedCephUserPoolPermission{PoolName: "pool", Read: true, Execute: true, ClassWrite: true}
	if len(permissions) != 1 || permissions[0] != want {
		t.Errorf("unexpected permissions:\n got: %v\nwant: [%v]", permissions, want)
	}
}

func TestDedicatedCephUserMergeWithPoolPermissions(t *testing.T) {
	ctx := context.Background()

	full := DedicatedCephUserPoolPermission{PoolName: "full", Read: true, Write: true}
	empty := DedicatedCephUserPoolPermission{PoolName: "empty"}

	tests := []struct {
		name        string
		prior       types.List
		permissions []DedicatedCephUserPoolPermission
		want        types.List
	}{
		{
			name:  "unset and no permission stays null",
			prior: dedicatedCephUserPoolPermissionsNull(),
			want:  dedicatedCephUserPoolPermissionsNull(),
		},
		{
			name:        "unset and only empty permissions stays null",
			prior:       dedicatedCephUserPoolPermissionsNull(),
			permissions: []DedicatedCephUserPoolPermission{empty},
			want:        dedicatedCephUserPoolPermissionsNull(),
		},
		{
			name:        "unset and permissions, as after an import",
			prior:       dedicatedCephUserPoolPermissionsNull(),
			permissions: []DedicatedCephUserPoolPermission{full, empty},
			want: dedicatedCephUserPoolPermissionsList(
				dedicatedCephUserPoolPermissionValue("full", true, true, false, false, false),
			),
		},
		{
			name: "configured and permissions removed outside of terraform becomes empty",
			prior: dedicatedCephUserPoolPermissionsList(
				dedicatedCephUserPoolPermissionValue("full", true, true, false, false, false),
			),
			want: dedicatedCephUserPoolPermissionsList(),
		},
		{
			name:  "configured empty list stays empty",
			prior: dedicatedCephUserPoolPermissionsList(),
			want:  dedicatedCephUserPoolPermissionsList(),
		},
		{
			name: "empty permission left out",
			prior: dedicatedCephUserPoolPermissionsList(
				dedicatedCephUserPoolPermissionValue("full", true, true, false, false, false),
			),
			permissions: []DedicatedCephUserPoolPermission{empty, full},
			want: dedicatedCephUserPoolPermissionsList(
				dedicatedCephUserPoolPermissionValue("full", true, true, false, false, false),
			),
		},
		{
			name: "configured order kept, others after in the order of the api",
			prior: dedicatedCephUserPoolPermissionsList(
				dedicatedCephUserPoolPermissionValue("b", true, false, false, false, false),
				dedicatedCephUserPoolPermissionValue("a", true, false, false, false, false),
			),
			permissions: []DedicatedCephUserPoolPermission{
				{PoolName: "a", Read: true},
				{PoolName: "d", Write: true},
				{PoolName: "b", Read: true},
				{PoolName: "c", Execute: true},
			},
			want: dedicatedCephUserPoolPermissionsList(
				dedicatedCephUserPoolPermissionValue("b", true, false, false, false, false),
				dedicatedCephUserPoolPermissionValue("a", true, false, false, false, false),
				dedicatedCephUserPoolPermissionValue("d", false, true, false, false, false),
				dedicatedCephUserPoolPermissionValue("c", false, false, true, false, false),
			),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := DedicatedCephUserModel{PoolPermissions: tt.prior}
			if diags := model.MergeWithPoolPermissions(ctx, tt.permissions); diags.HasError() {
				t.Fatalf("unexpected diagnostics: %v", diags)
			}
			if !model.PoolPermissions.Equal(tt.want) {
				t.Errorf("unexpected pool_permissions:\n got: %v\nwant: %v", model.PoolPermissions, tt.want)
			}
		})
	}
}

func TestDedicatedCephUserMergeWith(t *testing.T) {
	mdsCaps := "allow rw"
	model := DedicatedCephUserModel{ServiceName: ovhtypes.NewTfStringValue("ceph")}

	model.MergeWith(&DedicatedCephUser{Name: "user", Key: "secret", MonCaps: "allow r", OsdCaps: "allow rw pool=p", MdsCaps: &mdsCaps})
	if got := model.ID.ValueString(); got != "ceph/user" {
		t.Errorf("unexpected id %q", got)
	}
	if got := model.MdsCaps.ValueString(); got != mdsCaps {
		t.Errorf("unexpected mds_caps %q", got)
	}

	model.MergeWith(&DedicatedCephUser{Name: "user"})
	if !model.MdsCaps.IsNull() {
		t.Errorf("expected a null mds_caps when the API returns none, got %v", model.MdsCaps)
	}
}

func TestDedicatedCephPoolMergeWith(t *testing.T) {
	model := DedicatedCephPoolModel{ServiceName: ovhtypes.NewTfStringValue("ceph")}
	model.MergeWith(&DedicatedCephPool{Name: "pool", PoolType: "REPLICATED", Backup: true, MinActiveReplicas: 2, ReplicaCount: 3})

	if got := model.ID.ValueString(); got != "ceph/pool" {
		t.Errorf("unexpected id %q", got)
	}
	if got := model.PoolType.ValueString(); got != "REPLICATED" {
		t.Errorf("unexpected pool_type %q", got)
	}
	if !model.Backup.ValueBool() {
		t.Error("expected backup to be true")
	}
	if got := model.MinActiveReplicas.ValueInt64(); got != 2 {
		t.Errorf("unexpected min_active_replicas %d", got)
	}
	if got := model.ReplicaCount.ValueInt64(); got != 3 {
		t.Errorf("unexpected replica_count %d", got)
	}
}

func TestDedicatedCephImportState(t *testing.T) {
	ctx := context.Background()

	for name, r := range map[string]interface {
		resource.Resource
		resource.ResourceWithImportState
	}{
		"pool": &dedicatedCephPoolResource{},
		"user": &dedicatedCephUserResource{},
	} {
		t.Run(name, func(t *testing.T) {
			var schemaResp resource.SchemaResponse
			r.Schema(ctx, resource.SchemaRequest{}, &schemaResp)

			importState := func(id string) resource.ImportStateResponse {
				resp := resource.ImportStateResponse{
					State: tfsdk.State{
						Schema: schemaResp.Schema,
						Raw:    tftypes.NewValue(schemaResp.Schema.Type().TerraformType(ctx), nil),
					},
				}
				r.ImportState(ctx, resource.ImportStateRequest{ID: id}, &resp)
				return resp
			}

			for _, id := range []string{"", "ceph", "ceph/", "/name"} {
				if resp := importState(id); !resp.Diagnostics.HasError() {
					t.Errorf("expected an error importing %q", id)
				}
			}

			resp := importState("ceph/name")
			if resp.Diagnostics.HasError() {
				t.Fatalf("unexpected diagnostics: %v", resp.Diagnostics)
			}
			for attribute, want := range map[string]string{"id": "ceph/name", "service_name": "ceph", "name": "name"} {
				var got ovhtypes.TfStringValue
				resp.Diagnostics.Append(resp.State.GetAttribute(ctx, path.Root(attribute), &got)...)
				if got.ValueString() != want {
					t.Errorf("unexpected %s %q, want %q", attribute, got.ValueString(), want)
				}
			}
			if resp.Diagnostics.HasError() {
				t.Errorf("unexpected diagnostics: %v", resp.Diagnostics)
			}
		})
	}
}

// The API returns the backup of a pool as a string, such as "False", although documented as a JSON
// boolean.
func TestDedicatedCephPoolUnmarshal(t *testing.T) {
	for name, tt := range map[string]struct {
		backup string
		want   bool
	}{
		"returned string":     {`"False"`, false},
		"capitalized true":    {`"True"`, true},
		"lowercase string":    {`"true"`, true},
		"documented boolean":  {`true`, true},
		"documented false":    {`false`, false},
		"missing from answer": {`null`, false},
	} {
		t.Run(name, func(t *testing.T) {
			body := `{"serviceName":"ceph","name":"pool","replicaCount":3,"minActiveReplicas":2,"poolType":"REPLICATED","backup":` + tt.backup + `}`

			var got DedicatedCephPool
			if err := json.Unmarshal([]byte(body), &got); err != nil {
				t.Fatalf("unexpected error: %s", err)
			}

			want := DedicatedCephPool{ServiceName: "ceph", Name: "pool", PoolType: "REPLICATED", Backup: DedicatedCephBool(tt.want), MinActiveReplicas: 2, ReplicaCount: 3}
			if got != want {
				t.Errorf("unexpected pool:\n got: %+v\nwant: %+v", got, want)
			}
		})
	}

	for _, backup := range []string{`"maybe"`, `{}`, `1`} {
		var pool DedicatedCephPool
		err := json.Unmarshal([]byte(`{"backup":`+backup+`}`), &pool)
		if err == nil || !strings.Contains(err.Error(), backup) {
			t.Errorf("expected an error naming %s, got %v", backup, err)
		}
	}
}

func TestRetryDedicatedCephChange(t *testing.T) {
	ctx := context.Background()
	locked := &ovh.APIError{Code: 403, Class: "Client::Forbidden", Message: `Cluster 113198ce-3d23-4c7d-90a2-69b89c9bedbd is locked`}

	t.Run("locked cluster retried until the change is made", func(t *testing.T) {
		calls := 0
		err := retryDedicatedCephChange(ctx, time.Minute, func() error {
			calls++
			if calls == 1 {
				return locked
			}
			if calls == 2 {
				// The pool deletion wraps the client errors.
				return fmt.Errorf("error calling Delete: %w", locked)
			}
			return nil
		})
		if err != nil || calls != 3 {
			t.Errorf("expected a success on the third call, got %v after %d calls", err, calls)
		}
	})

	serverError := &ovh.APIError{Code: 500, Class: "Server::InternalServerError", Message: "Internal server error"}

	t.Run("server error retried", func(t *testing.T) {
		calls := 0
		err := retryDedicatedCephChange(ctx, time.Minute, func() error {
			calls++
			if calls == 1 {
				return serverError
			}
			return nil
		})
		if err != nil || calls != 2 {
			t.Errorf("expected a success on the second call, got %v after %d calls", err, calls)
		}
	})

	t.Run("server error retried a few times only", func(t *testing.T) {
		calls := 0
		err := retryDedicatedCephChange(ctx, time.Minute, func() error {
			calls++
			return fmt.Errorf("error calling Post: %w", serverError)
		})
		if calls != dedicatedCephServerErrorAttempts {
			t.Errorf("expected %d calls, got %d", dedicatedCephServerErrorAttempts, calls)
		}
		if !errors.Is(err, serverError) {
			t.Errorf("expected the server error, got %v", err)
		}
	})

	for name, apiErr := range map[string]error{
		"other forbidden error": &ovh.APIError{Code: 403, Class: "Client::Forbidden", Message: "This call has not been granted"},
		"not found":             &ovh.APIError{Code: 404, Class: "Client::NotFound", Message: "Pool not found"},
		"not an API error":      errors.New("connection refused"),
	} {
		t.Run(name+" not retried", func(t *testing.T) {
			calls := 0
			err := retryDedicatedCephChange(ctx, time.Minute, func() error {
				calls++
				return fmt.Errorf("error calling Post: %w", apiErr)
			})
			if calls != 1 {
				t.Errorf("expected a single call, got %d", calls)
			}
			// The error is returned as is, so that a 404 can still be told apart by the caller.
			if !errors.Is(err, apiErr) {
				t.Errorf("expected the error of the change, got %v", err)
			}
		})
	}
}

func TestLockDedicatedCephCluster(t *testing.T) {
	ctx := context.Background()

	unlock, err := lockDedicatedCephCluster(ctx, "ceph-1")
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	// Another change of the same cluster waits, until given up.
	cancelled, cancel := context.WithCancel(ctx)
	cancel()
	if _, err := lockDedicatedCephCluster(cancelled, "ceph-1"); !errors.Is(err, context.Canceled) {
		t.Errorf("expected the change of a locked cluster to wait, got %v", err)
	}

	// The changes of other clusters do not wait.
	unlockOther, err := lockDedicatedCephCluster(ctx, "ceph-2")
	if err != nil {
		t.Fatalf("expected another cluster to be changed right away, got %s", err)
	}
	unlockOther()

	unlock()
	unlock, err = lockDedicatedCephCluster(ctx, "ceph-1")
	if err != nil {
		t.Fatalf("expected the cluster to be changed once released, got %s", err)
	}
	unlock()
}

// newDedicatedCephTestClient returns a client of a fake API answering the given path with the given
// answers, one per call, the last one being repeated.
func newDedicatedCephTestClient(t *testing.T, endpoint string, calls *int32, answers ...func(w http.ResponseWriter)) *ovhwrap.Client {
	// Only the application key of the client must be configured.
	for _, name := range []string{"OVH_CLIENT_ID", "OVH_CLIENT_SECRET", "OVH_ACCESS_TOKEN"} {
		t.Setenv(name, "")
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/auth/time", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, time.Now().Unix())
	})
	mux.HandleFunc(endpoint, func(w http.ResponseWriter, r *http.Request) {
		call := int(atomic.AddInt32(calls, 1))
		if call > len(answers) {
			call = len(answers)
		}
		answers[call-1](w)
	})

	server := httptest.NewServer(mux)
	t.Cleanup(server.Close)

	client, err := ovh.NewClient(server.URL, "key", "secret", "consumer")
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	return ovhwrap.NewClient(client, ratelimit.NewUnlimited())
}

func TestWaitDedicatedCephIdle(t *testing.T) {
	ctx := context.Background()

	inProgress := func(w http.ResponseWriter) {
		fmt.Fprint(w, `[{"id":"89278f0b-6de8-454d-8f14-36c480454838","name":"createPool"}]`)
	}
	idle := func(w http.ResponseWriter) {
		fmt.Fprint(w, `[]`)
	}
	serverError := func(w http.ResponseWriter) {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, `{"class":"Server::InternalServerError","message":"Internal server error"}`)
	}

	t.Run("waits for the tasks in progress", func(t *testing.T) {
		var calls int32
		client := newDedicatedCephTestClient(t, "/dedicated/ceph/ceph/task", &calls, inProgress, inProgress, idle)
		if err := waitDedicatedCephIdle(ctx, client, "ceph", time.Minute); err != nil {
			t.Fatalf("unexpected error: %s", err)
		}
		if got := atomic.LoadInt32(&calls); got != 3 {
			t.Errorf("expected 3 calls, got %d", got)
		}
	})

	t.Run("tolerates a server error", func(t *testing.T) {
		var calls int32
		client := newDedicatedCephTestClient(t, "/dedicated/ceph/ceph/task", &calls, serverError, idle)
		if err := waitDedicatedCephIdle(ctx, client, "ceph", time.Minute); err != nil {
			t.Fatalf("unexpected error: %s", err)
		}
	})

	t.Run("gives up on server errors in a row", func(t *testing.T) {
		var calls int32
		client := newDedicatedCephTestClient(t, "/dedicated/ceph/ceph/task", &calls, serverError)
		err := waitDedicatedCephIdle(ctx, client, "ceph", time.Minute)

		var errOvh *ovh.APIError
		if !errors.As(err, &errOvh) || errOvh.Code != http.StatusInternalServerError {
			t.Errorf("expected the server error, got %v", err)
		}
		if got := atomic.LoadInt32(&calls); got != dedicatedCephServerErrorAttempts {
			t.Errorf("expected %d calls, got %d", dedicatedCephServerErrorAttempts, got)
		}
	})
}

// The state saved as soon as a creation is accepted must hold no unknown value, which a state
// cannot hold, while the plan it comes from has the attributes read from the API unknown.
func TestDedicatedCephCreatedState(t *testing.T) {
	ctx := context.Background()
	unknownString := ovhtypes.TfStringValue{StringValue: basetypes.NewStringUnknown()}

	storable := func(t *testing.T, r resource.Resource, model interface{}) {
		var schemaResp resource.SchemaResponse
		r.Schema(ctx, resource.SchemaRequest{}, &schemaResp)

		state := tfsdk.State{
			Schema: schemaResp.Schema,
			Raw:    tftypes.NewValue(schemaResp.Schema.Type().TerraformType(ctx), nil),
		}
		if diags := state.Set(ctx, model); diags.HasError() {
			t.Fatalf("unexpected diagnostics: %v", diags)
		}
		if !state.Raw.IsFullyKnown() {
			t.Errorf("expected a fully known state, got %v", state.Raw)
		}
	}

	t.Run("pool", func(t *testing.T) {
		pool := DedicatedCephPoolModel{
			ID:                unknownString,
			ServiceName:       ovhtypes.NewTfStringValue("ceph"),
			Name:              ovhtypes.NewTfStringValue("pool"),
			PoolType:          unknownString,
			Backup:            ovhtypes.TfBoolValue{BoolValue: basetypes.NewBoolUnknown()},
			MinActiveReplicas: ovhtypes.TfInt64Value{Int64Value: basetypes.NewInt64Unknown()},
			ReplicaCount:      ovhtypes.TfInt64Value{Int64Value: basetypes.NewInt64Unknown()},
		}.CreatedState()

		if got := pool.ID.ValueString(); got != "ceph/pool" {
			t.Errorf("unexpected id %q", got)
		}
		storable(t, &dedicatedCephPoolResource{}, &pool)
	})

	t.Run("user", func(t *testing.T) {
		user := DedicatedCephUserModel{
			ID:          unknownString,
			ServiceName: ovhtypes.NewTfStringValue("ceph"),
			Name:        ovhtypes.NewTfStringValue("user"),
			PoolPermissions: dedicatedCephUserPoolPermissionsList(
				dedicatedCephUserPoolPermissionValue("pool", true, false, false, false, false),
			),
			Key:     unknownString,
			MonCaps: unknownString,
			OsdCaps: unknownString,
			MdsCaps: unknownString,
		}.CreatedState()

		if got := user.ID.ValueString(); got != "ceph/user" {
			t.Errorf("unexpected id %q", got)
		}
		// The pool permissions are not set yet when the user is created.
		if !user.PoolPermissions.IsNull() {
			t.Errorf("expected null pool permissions, got %v", user.PoolPermissions)
		}
		storable(t, &dedicatedCephUserResource{}, &user)
	})
}

func TestValidateDedicatedCephUserPoolPermissions(t *testing.T) {
	ctx := context.Background()

	readUnknown := types.ObjectValueMust(DedicatedCephUserPoolPermissionAttrTypes(), map[string]attr.Value{
		"pool_name":   types.StringValue("pending"),
		"read":        types.BoolUnknown(),
		"write":       types.BoolValue(false),
		"execute":     types.BoolValue(false),
		"class_read":  types.BoolValue(false),
		"class_write": types.BoolValue(false),
	})

	for name, permissions := range map[string]types.List{
		"null":  dedicatedCephUserPoolPermissionsNull(),
		"empty": dedicatedCephUserPoolPermissionsList(),
		"granting": dedicatedCephUserPoolPermissionsList(
			dedicatedCephUserPoolPermissionValue("pool", false, false, false, false, true),
		),
		// When validating the configuration, a flag may not be known yet.
		"flag unknown": dedicatedCephUserPoolPermissionsList(readUnknown),
	} {
		if diags := validateDedicatedCephUserPoolPermissions(ctx, permissions); diags.HasError() {
			t.Errorf("%s: unexpected diagnostics: %v", name, diags)
		}
	}

	// When applying, the same flag is known and may turn out to grant nothing.
	diags := validateDedicatedCephUserPoolPermissions(ctx, dedicatedCephUserPoolPermissionsList(
		dedicatedCephUserPoolPermissionValue("pool", true, false, false, false, false),
		dedicatedCephUserPoolPermissionValue("pending", false, false, false, false, false),
	))
	if diags.ErrorsCount() != 1 {
		t.Fatalf("expected a single error, got %v", diags)
	}
	withPath, ok := diags.Errors()[0].(diag.DiagnosticWithPath)
	if !ok || !withPath.Path().Equal(path.Root("pool_permissions").AtListIndex(1)) {
		t.Errorf("expected the error on pool_permissions[1], got %v", diags.Errors()[0])
	}

	// A pool listed twice is refused, even with different permissions.
	diags = validateDedicatedCephUserPoolPermissions(ctx, dedicatedCephUserPoolPermissionsList(
		dedicatedCephUserPoolPermissionValue("pool", true, false, false, false, false),
		dedicatedCephUserPoolPermissionValue("other", true, false, false, false, false),
		dedicatedCephUserPoolPermissionValue("pool", false, true, false, false, false),
	))
	if diags.ErrorsCount() != 1 {
		t.Fatalf("expected a single error, got %v", diags)
	}
	withPath, ok = diags.Errors()[0].(diag.DiagnosticWithPath)
	if !ok || !withPath.Path().Equal(path.Root("pool_permissions").AtListIndex(2).AtName("pool_name")) {
		t.Errorf("expected the error on pool_permissions[2].pool_name, got %v", diags.Errors()[0])
	}

	// When validating the configuration, a pool name may not be known yet.
	unknownName := func(read bool) attr.Value {
		return types.ObjectValueMust(DedicatedCephUserPoolPermissionAttrTypes(), map[string]attr.Value{
			"pool_name":   types.StringUnknown(),
			"read":        types.BoolValue(read),
			"write":       types.BoolValue(!read),
			"execute":     types.BoolValue(false),
			"class_read":  types.BoolValue(false),
			"class_write": types.BoolValue(false),
		})
	}
	diags = validateDedicatedCephUserPoolPermissions(ctx, dedicatedCephUserPoolPermissionsList(unknownName(true), unknownName(false)))
	if diags.HasError() {
		t.Errorf("unexpected diagnostics for unknown pool names: %v", diags)
	}
}

func TestWaitDedicatedCephTask(t *testing.T) {
	const endpoint = "/dedicated/ceph/ceph/task/89278f0b-6de8-454d-8f14-36c480454838"

	task := func(state string) func(w http.ResponseWriter) {
		return func(w http.ResponseWriter) {
			fmt.Fprintf(w, `[{"name":"createPool","state":%q,"type":"pool","createDate":"2026-10-05T08:24:29Z","finishDate":null}]`, state)
		}
	}
	notListed := func(w http.ResponseWriter) {
		fmt.Fprint(w, `[]`)
	}
	notFound := func(w http.ResponseWriter) {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, `{"class":"Client::NotFound","message":"Task not found"}`)
	}
	serverError := func(w http.ResponseWriter) {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, `{"class":"Server::InternalServerError","message":"Internal server error"}`)
	}
	forbidden := func(w http.ResponseWriter) {
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprint(w, `{"class":"Client::Forbidden","message":"This call has not been granted"}`)
	}

	wait := func(t *testing.T, ctx context.Context, timeout time.Duration, calls *int32, answers ...func(w http.ResponseWriter)) error {
		client := newDedicatedCephTestClient(t, endpoint, calls, answers...)
		return waitDedicatedCephTask(ctx, client, "ceph", "89278f0b-6de8-454d-8f14-36c480454838", timeout)
	}

	// The polls happen at about 0s, 0.5s, 1.5s, 3.5s and 7.5s, so the sequences below are as short
	// as they can be: make test runs the unit tests of the package within 30 seconds.

	t.Run("done after being registered", func(t *testing.T) {
		var calls int32
		err := wait(t, context.Background(), time.Minute, &calls, notListed, notFound, task("DONE"))
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}
		if got := atomic.LoadInt32(&calls); got != 3 {
			t.Errorf("expected 3 calls, got %d", got)
		}
	})

	t.Run("failed task stops the wait", func(t *testing.T) {
		var calls int32
		err := wait(t, context.Background(), time.Minute, &calls, task("FAILED"), task("DONE"))
		if err == nil || !strings.Contains(err.Error(), "failed") {
			t.Errorf("expected the task to fail, got %v", err)
		}
		if got := atomic.LoadInt32(&calls); got != 1 {
			t.Errorf("expected a single call, got %d", got)
		}
	})

	t.Run("server errors tolerated when not in a row", func(t *testing.T) {
		// Without the count being reset by a successful poll, the last server error would be the
		// third one and stop the wait.
		var calls int32
		err := wait(t, context.Background(), time.Minute, &calls,
			serverError, task("IN PROGRESS"), serverError, serverError, task("DONE"))
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}
		if got := atomic.LoadInt32(&calls); got != 5 {
			t.Errorf("expected 5 calls, got %d", got)
		}
	})

	t.Run("gives up on server errors in a row", func(t *testing.T) {
		var calls int32
		err := wait(t, context.Background(), time.Minute, &calls, serverError)

		var errOvh *ovh.APIError
		if !errors.As(err, &errOvh) || errOvh.Code != http.StatusInternalServerError {
			t.Errorf("expected the server error, got %v", err)
		}
		if got := atomic.LoadInt32(&calls); got != dedicatedCephServerErrorAttempts {
			t.Errorf("expected %d calls, got %d", dedicatedCephServerErrorAttempts, got)
		}
	})

	t.Run("other client errors stop the wait", func(t *testing.T) {
		var calls int32
		err := wait(t, context.Background(), time.Minute, &calls, forbidden, task("DONE"))

		var errOvh *ovh.APIError
		if !errors.As(err, &errOvh) || errOvh.Code != http.StatusForbidden {
			t.Errorf("expected the forbidden error, got %v", err)
		}
		if got := atomic.LoadInt32(&calls); got != 1 {
			t.Errorf("expected a single call, got %d", got)
		}
	})

	t.Run("times out on a task in progress", func(t *testing.T) {
		var calls int32
		start := time.Now()
		err := wait(t, context.Background(), time.Second, &calls, task("IN PROGRESS"))
		if err == nil {
			t.Fatal("expected the wait to time out")
		}
		if elapsed := time.Since(start); elapsed > 30*time.Second {
			t.Errorf("expected the wait to stop at its timeout, it took %s", elapsed)
		}
	})

	t.Run("stops when cancelled", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		var calls int32
		start := time.Now()
		err := wait(t, ctx, time.Minute, &calls, func(w http.ResponseWriter) {
			// Terraform is cancelled while the task is in progress.
			cancel()
			task("IN PROGRESS")(w)
		})
		if err == nil {
			t.Fatal("expected the wait to stop")
		}
		if elapsed := time.Since(start); elapsed > 30*time.Second {
			t.Errorf("expected the wait to stop when cancelled, it took %s", elapsed)
		}
		if got := atomic.LoadInt32(&calls); got > 2 {
			t.Errorf("expected the polling to stop when cancelled, got %d calls", got)
		}
	})
}
