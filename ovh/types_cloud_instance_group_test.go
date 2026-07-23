package ovh

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestUnitCloudInstanceGroupToCreate(t *testing.T) {
	m := &CloudInstanceGroupModel{
		ServiceName: strVal("proj"),
		Region:      strVal("GRA11"),
		Name:        strVal("grp-1"),
		Policy:      strVal("ANTI_AFFINITY"),
	}

	payload := m.ToCreate()

	if payload.TargetSpec == nil {
		t.Fatalf("target spec should not be nil")
	}
	if payload.TargetSpec.Name != "grp-1" {
		t.Fatalf("name = %q, want grp-1", payload.TargetSpec.Name)
	}
	if payload.TargetSpec.Policy != "ANTI_AFFINITY" {
		t.Fatalf("policy = %q, want ANTI_AFFINITY", payload.TargetSpec.Policy)
	}
	if payload.TargetSpec.Location == nil || payload.TargetSpec.Location.Region != "GRA11" {
		t.Fatalf("location.region not set correctly: %+v", payload.TargetSpec.Location)
	}
}

func TestUnitCloudInstanceGroupMergeWith(t *testing.T) {
	ctx := context.Background()
	resp := &CloudInstanceGroupAPIResponse{
		Id:             "grp-1",
		Checksum:       "chk-9",
		CreatedAt:      "2026-07-09T10:00:00Z",
		UpdatedAt:      "2026-07-09T10:05:00Z",
		ResourceStatus: "READY",
		TargetSpec: &CloudInstanceGroupAPITargetSpec{
			Name:     "grp-1",
			Policy:   "ANTI_AFFINITY",
			Location: &CloudInstanceGroupAPILocation{Region: "GRA11"},
		},
		CurrentState: &CloudInstanceGroupAPICurrentState{
			Name:     "grp-1",
			Policy:   "ANTI_AFFINITY",
			Location: &CloudInstanceGroupAPILocation{Region: "GRA11", AvailabilityZone: "GRA11-a"},
			Members:  []CloudInstanceGroupAPIMemberRef{{Id: "inst-1"}},
		},
	}

	m := &CloudInstanceGroupModel{ServiceName: strVal("proj")}
	m.MergeWith(ctx, resp)

	// Envelope
	if m.Id.ValueString() != "grp-1" {
		t.Fatalf("id = %q", m.Id.ValueString())
	}
	if m.Checksum.ValueString() != "chk-9" {
		t.Fatalf("checksum = %q", m.Checksum.ValueString())
	}
	if m.ResourceStatus.ValueString() != "READY" {
		t.Fatalf("resource_status = %q", m.ResourceStatus.ValueString())
	}

	// Flattened fields come from TargetSpec
	if m.Name.ValueString() != "grp-1" {
		t.Fatalf("name = %q, want grp-1", m.Name.ValueString())
	}
	if m.Region.ValueString() != "GRA11" {
		t.Fatalf("region = %q, want GRA11", m.Region.ValueString())
	}
	if m.Policy.ValueString() != "ANTI_AFFINITY" {
		t.Fatalf("policy = %q, want ANTI_AFFINITY", m.Policy.ValueString())
	}

	// current_state must be non-null and carry members
	if m.CurrentState.IsNull() {
		t.Fatalf("current_state should not be null")
	}
	csAttrs := m.CurrentState.Attributes()
	membersList, ok := csAttrs["members"].(types.List)
	if !ok || membersList.IsNull() {
		t.Fatalf("current_state.members missing")
	}
	if len(membersList.Elements()) != 1 {
		t.Fatalf("current_state.members len = %d, want 1", len(membersList.Elements()))
	}
	locObj, ok := csAttrs["location"].(types.Object)
	if !ok || locObj.IsNull() {
		t.Fatalf("current_state.location missing")
	}
}

func TestUnitCloudInstanceGroupMergeWithNilCurrentState(t *testing.T) {
	ctx := context.Background()
	resp := &CloudInstanceGroupAPIResponse{
		Id:             "grp-2",
		ResourceStatus: "CREATING",
		TargetSpec: &CloudInstanceGroupAPITargetSpec{
			Name:     "grp-2",
			Policy:   "AFFINITY",
			Location: &CloudInstanceGroupAPILocation{Region: "GRA11"},
		},
		CurrentState: nil,
	}
	m := &CloudInstanceGroupModel{}
	m.MergeWith(ctx, resp)
	if !m.CurrentState.IsNull() {
		t.Fatalf("current_state should be null when API currentState is nil")
	}
}
