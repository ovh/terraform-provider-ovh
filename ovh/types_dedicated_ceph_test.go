package ovh

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/ovh/go-ovh/ovh"
	ovhtypes "github.com/ovh/terraform-provider-ovh/v2/ovh/types"
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

func dedicatedCephUserPoolPermissionsSet(elements ...attr.Value) types.Set {
	return types.SetValueMust(types.ObjectType{AttrTypes: DedicatedCephUserPoolPermissionAttrTypes()}, elements)
}

func dedicatedCephUserPoolPermissionsNull() types.Set {
	return types.SetNull(types.ObjectType{AttrTypes: DedicatedCephUserPoolPermissionAttrTypes()})
}

// The attr types of pool_permissions are declared apart from the schema, so a nested attribute
// added to one side only produces a runtime error rather than a compile error.
func TestDedicatedCephUserSchemaMatchesAttrTypes(t *testing.T) {
	var resp resource.SchemaResponse
	(&dedicatedCephUserResource{}).Schema(context.Background(), resource.SchemaRequest{}, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("schema diagnostics: %v", resp.Diagnostics)
	}

	want := types.SetType{ElemType: types.ObjectType{AttrTypes: DedicatedCephUserPoolPermissionAttrTypes()}}
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

func TestDedicatedCephUserToPoolPermissions(t *testing.T) {
	ctx := context.Background()

	model := DedicatedCephUserModel{PoolPermissions: dedicatedCephUserPoolPermissionsNull()}
	permissions, diags := model.ToPoolPermissions(ctx)
	if diags.HasError() {
		t.Fatalf("unexpected diagnostics: %v", diags)
	}
	if len(permissions) != 0 {
		t.Errorf("expected no permission from a null set, got %v", permissions)
	}

	model.PoolPermissions = dedicatedCephUserPoolPermissionsSet(
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
		prior       types.Set
		permissions []DedicatedCephUserPoolPermission
		want        types.Set
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
			want: dedicatedCephUserPoolPermissionsSet(
				dedicatedCephUserPoolPermissionValue("full", true, true, false, false, false),
			),
		},
		{
			name: "set and permissions removed outside of terraform becomes empty",
			prior: dedicatedCephUserPoolPermissionsSet(
				dedicatedCephUserPoolPermissionValue("full", true, true, false, false, false),
			),
			want: dedicatedCephUserPoolPermissionsSet(),
		},
		{
			name:  "configured empty set stays empty",
			prior: dedicatedCephUserPoolPermissionsSet(),
			want:  dedicatedCephUserPoolPermissionsSet(),
		},
		{
			name: "empty permission kept when configured",
			prior: dedicatedCephUserPoolPermissionsSet(
				dedicatedCephUserPoolPermissionValue("empty", false, false, false, false, false),
			),
			permissions: []DedicatedCephUserPoolPermission{full, empty},
			want: dedicatedCephUserPoolPermissionsSet(
				dedicatedCephUserPoolPermissionValue("full", true, true, false, false, false),
				dedicatedCephUserPoolPermissionValue("empty", false, false, false, false, false),
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
