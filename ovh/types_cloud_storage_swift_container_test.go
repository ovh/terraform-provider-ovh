package ovh

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	ovhtypes "github.com/ovh/terraform-provider-ovh/v2/ovh/types"
)

func storageSwiftContainerTargetSpec(t *testing.T, payload any) map[string]json.RawMessage {
	t.Helper()

	raw, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("marshal payload: %v", err)
	}

	var body struct {
		TargetSpec map[string]json.RawMessage `json:"targetSpec"`
	}
	if err := json.Unmarshal(raw, &body); err != nil {
		t.Fatalf("unmarshal payload: %v", err)
	}

	return body.TargetSpec
}

func TestCloudStorageSwiftContainerToUpdate_CorsAlwaysSent(t *testing.T) {
	cases := map[string]struct {
		cors types.Set
		want string
	}{
		"null set clears origins": {
			cors: types.SetNull(ovhtypes.TfStringType{}),
			want: `[]`,
		},
		"empty set clears origins": {
			cors: types.SetValueMust(ovhtypes.TfStringType{}, []attr.Value{}),
			want: `[]`,
		},
		"values are sent": {
			cors: types.SetValueMust(ovhtypes.TfStringType{}, []attr.Value{
				ovhtypes.NewTfStringValue("https://www.example.com"),
			}),
			want: `["https://www.example.com"]`,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			model := CloudStorageSwiftContainerModel{Type: ovhtypes.NewTfStringValue("PUBLIC"), Cors: tc.cors}
			spec := storageSwiftContainerTargetSpec(t, model.ToUpdate("checksum-123"))

			got, ok := spec["cors"]
			if !ok {
				t.Fatalf("cors key missing from update target spec: %v", spec)
			}
			if string(got) != tc.want {
				t.Fatalf("cors = %s, want %s", got, tc.want)
			}
			if string(spec["type"]) != `"PUBLIC"` {
				t.Fatalf("type = %s, want \"PUBLIC\"", spec["type"])
			}
		})
	}
}

func TestCloudStorageSwiftContainerToCreate_UnsetDefaultsOmitted(t *testing.T) {
	model := CloudStorageSwiftContainerModel{
		Name:    ovhtypes.NewTfStringValue("container"),
		Region:  ovhtypes.NewTfStringValue("GRA"),
		Archive: types.BoolUnknown(),
		Type:    ovhtypes.NewTfStringNull(),
		Cors:    types.SetUnknown(ovhtypes.TfStringType{}),
	}

	spec := storageSwiftContainerTargetSpec(t, model.ToCreate())
	for _, key := range []string{"archive", "type", "cors"} {
		if _, ok := spec[key]; ok {
			t.Fatalf("%s should be omitted from create target spec when unset: %v", key, spec)
		}
	}
	if string(spec["name"]) != `"container"` || string(spec["location"]) != `{"region":"GRA"}` {
		t.Fatalf("unexpected immutable fields in create target spec: %v", spec)
	}
}

func TestCloudStorageSwiftContainerMergeWith_FillsImmutableOnlyWhenUnset(t *testing.T) {
	archive := true
	response := &CloudStorageSwiftContainerAPIResponse{
		Id:             "GRA_container",
		ResourceStatus: "READY",
		TargetSpec: &CloudStorageSwiftContainerAPITargetSpec{
			Name:     "container",
			Location: &CloudStorageSwiftContainerAPILocation{Region: "GRA"},
			Archive:  &archive,
			Type:     "STATIC",
			Cors:     []string{"https://www.example.com"},
		},
		CurrentState: &CloudStorageSwiftContainerAPICurrentState{
			Name:     "container",
			Location: &CloudStorageSwiftContainerAPILocation{Region: "GRA"},
			Archive:  true,
			Type:     "STATIC",
		},
	}

	imported := CloudStorageSwiftContainerModel{}
	imported.MergeWith(context.Background(), response)
	if imported.Name.ValueString() != "container" || imported.Region.ValueString() != "GRA" {
		t.Fatalf("import should fill name and region: %+v", imported)
	}
	if !imported.Archive.ValueBool() || imported.Type.ValueString() != "STATIC" || len(imported.Cors.Elements()) != 1 {
		t.Fatalf("import should fill archive, type and cors: %+v", imported)
	}
	if imported.CurrentState.IsNull() {
		t.Fatalf("current_state should be set")
	}

	configured := CloudStorageSwiftContainerModel{
		Name:   ovhtypes.NewTfStringValue("container"),
		Region: ovhtypes.NewTfStringValue("GRA"),
	}
	response.TargetSpec.Location.Region = "gra"
	configured.MergeWith(context.Background(), response)
	if configured.Region.ValueString() != "GRA" {
		t.Fatalf("a configured region must not be overwritten by the API: %s", configured.Region.ValueString())
	}

	response.TargetSpec = nil
	response.CurrentState = nil
	fresh := CloudStorageSwiftContainerModel{}
	fresh.MergeWith(context.Background(), response)
	if !fresh.CurrentState.IsNull() {
		t.Fatalf("current_state should be a typed null when the API omits it")
	}
}
