package ovh

import (
	"context"
	"testing"
)

// sampleQuotaResponse returns a representative envelope exercising two regions,
// a populated current_state, and an in-flight currentTask carrying an error.
func sampleQuotaResponse() *CloudQuotaAPIResponse {
	return &CloudQuotaAPIResponse{
		Id:             "proj-1",
		ResourceStatus: "READY",
		Checksum:       "abc123",
		CreatedAt:      "2024-01-01T00:00:00Z",
		UpdatedAt:      "2024-01-02T00:00:00Z",
		TargetSpec: &CloudQuotaTargetSpecAPI{
			PreventAutomaticQuotaUpgrade: true,
			Regions: []CloudQuotaRegionTargetSpecAPI{
				{
					Location: &CloudQuotaLocationAPI{Region: "GRA11"},
					Profile:  "default",
				},
				{
					Location: &CloudQuotaLocationAPI{Region: "SBG5"},
					Profile:  "large",
				},
			},
		},
		CurrentState: &CloudQuotaCurrentStateAPI{
			PreventAutomaticQuotaUpgrade: true,
			AvailableProfiles: []CloudQuotaAvailableProfileAPI{
				{Name: "default", Compute: &CloudQuotaProfileComputeAPI{Cores: 10, Instances: 5, Memory: 20480}},
			},
			Regions: []CloudQuotaRegionCurrentStateAPI{
				{
					Location: &CloudQuotaLocationAPI{Region: "GRA11"},
					Profile:  "default",
					Usage: &CloudQuotaUsageDetailsAPI{
						Compute: &CloudQuotaRegionComputeAPI{
							Cores:     CloudQuotaUsageAPI{Limit: 10, Used: intPtr(2), Unit: "COUNT"},
							Instances: CloudQuotaUsageAPI{Limit: 5, Used: nil, Unit: "COUNT"},
							Memory:    CloudQuotaUsageAPI{Limit: 20480, Used: intPtr(1024), Unit: "MB"},
						},
					},
				},
			},
		},
		CurrentTasks: []CloudQuotaCurrentTaskAPI{
			{
				Id:     "task-1",
				Type:   "quota/edit",
				Status: "ERROR",
				Link:   "/v2/.../task-1",
				Errors: []CloudQuotaTaskErrorAPI{{Message: "region unavailable"}},
			},
		},
	}
}

// TestCloudQuotaResourceModelMergeWith verifies that MergeWith maps the
// flattened targetSpec regions into the model.
func TestCloudQuotaResourceModelMergeWith(t *testing.T) {
	ctx := context.Background()
	var m CloudQuotaResourceModel
	m.MergeWith(ctx, sampleQuotaResponse())

	if got := m.Checksum.ValueString(); got != "abc123" {
		t.Fatalf("checksum = %q, want abc123", got)
	}
	if !m.PreventAutomaticQuotaUpgrade.ValueBool() {
		t.Fatalf("prevent_automatic_quota_upgrade = false, want true")
	}

	var regions []quotaRegionPlan
	if diags := m.Regions.ElementsAs(ctx, &regions, false); diags.HasError() {
		t.Fatalf("regions ElementsAs: %s", diags)
	}
	if len(regions) != 2 {
		t.Fatalf("regions len = %d, want 2", len(regions))
	}
	if regions[0].Region.ValueString() != "GRA11" {
		t.Errorf("regions[0].region = %q, want GRA11", regions[0].Region.ValueString())
	}
	if regions[0].Profile.ValueString() != "default" {
		t.Errorf("regions[0].profile = %q, want default", regions[0].Profile.ValueString())
	}
}

// TestCloudQuotaToTargetSpecAPI verifies the regions are written back into the
// PUT payload.
func TestCloudQuotaToTargetSpecAPI(t *testing.T) {
	ctx := context.Background()
	var m CloudQuotaResourceModel
	m.MergeWith(ctx, sampleQuotaResponse())

	spec, err := m.toTargetSpecAPI(ctx)
	if err != nil {
		t.Fatalf("toTargetSpecAPI: %v", err)
	}
	if !spec.PreventAutomaticQuotaUpgrade {
		t.Errorf("preventAutomaticQuotaUpgrade = false, want true")
	}
	if len(spec.Regions) != 2 {
		t.Fatalf("regions len = %d, want 2", len(spec.Regions))
	}
	if spec.Regions[0].Location == nil || spec.Regions[0].Location.Region != "GRA11" {
		t.Errorf("regions[0] region = %+v, want GRA11", spec.Regions[0].Location)
	}
}

// TestQuotaTaskErrorSuffix verifies failed-reconcile errors are surfaced from
// the envelope's currentTasks.
func TestQuotaTaskErrorSuffix(t *testing.T) {
	if got := quotaTaskErrorSuffix(&CloudQuotaAPIResponse{}); got != "" {
		t.Errorf("empty envelope suffix = %q, want empty", got)
	}
	got := quotaTaskErrorSuffix(sampleQuotaResponse())
	if got != ": region unavailable" {
		t.Errorf("suffix = %q, want ': region unavailable'", got)
	}
}
