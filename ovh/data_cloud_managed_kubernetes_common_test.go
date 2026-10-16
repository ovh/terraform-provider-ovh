package ovh

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// Decodes the OpenAPI examples into each model and sets them as state, catching tfsdk/schema mismatches offline.
func TestCloudManagedKubernetesReferenceModels_stateSet(t *testing.T) {
	ctx := context.Background()
	cases := []struct {
		ds      datasource.DataSource
		payload string
		target  func() (any, any)
	}{
		{NewCloudManagedKubernetesPlansDataSource(), `[{"name":"FREE","status":"AVAILABLE","limits":{"maxNodes":100},"controlPlane":{"etcd":{"dedicated":false,"storage":{"limit":400,"unit":"MiB"}}}}]`,
			func() (any, any) { m := &cloudManagedKubernetesPlansModel{}; return m, &m.Plans }},
		{NewCloudManagedKubernetesRegionsDataSource(), `[{"name":"BHS5","availabilityZones":[],"availablePlans":["FREE","STANDARD"]}]`,
			func() (any, any) { m := &cloudManagedKubernetesRegionsModel{}; return m, &m.Regions }},
		{NewCloudManagedKubernetesVersionsDataSource(), `{"availableVersions":[{"status":"SUPPORTED","version":"1.34"}],"defaultVersion":"1.34"}`,
			func() (any, any) { m := &cloudManagedKubernetesVersionsModel{}; return m, m }},
	}
	for _, c := range cases {
		var sresp datasource.SchemaResponse
		c.ds.Schema(ctx, datasource.SchemaRequest{}, &sresp)
		model, dst := c.target()
		if err := json.Unmarshal([]byte(c.payload), dst); err != nil {
			t.Fatal(err)
		}
		state := tfsdk.State{Schema: sresp.Schema, Raw: tftypes.NewValue(sresp.Schema.Type().TerraformType(ctx), nil)}
		if diags := state.Set(ctx, model); diags.HasError() {
			t.Fatalf("%T: %v", c.ds, diags)
		}
	}
}
