package ovh

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestCloudManagedDatabaseCapabilitiesDataSource(t *testing.T) {
	base := "/v2/publicCloud/project/pj/managedDatabase/capabilities/"
	lifecycle := `{"startedAt": "2024-01-01", "endOfSaleAt": null, "endOfLifeAt": null, "status": "STABLE"}`
	state := managedDatabaseDataSourceTestRead(t, NewCloudManagedDatabaseCapabilitiesDataSource(), map[string]string{"service_name": "pj"}, map[string][]managedDatabaseTestAnswer{
		"GET " + base + "disks": {{http.StatusOK, `[{"name": "HIGH_SPEED", "maximumIops": 3000, "iopsPerGigabyte": null, "maximumThroughput": 250.5, "throughputPerGigabyte": null}]`}},
		"GET " + base + "engines": {{http.StatusOK, `[{"name": "postgresql", "description": "PostgreSQL", "category": "DATABASE", "defaultVersion": "17", "sslModes": ["require"],
			"storage": "REPLICATED", "order": 1, "tags": null, "lifecycle": ` + lifecycle + `, "versions": [{"name": "17", "default": true, "tags": null, "lifecycle": ` + lifecycle + `}]}]`}},
		"GET " + base + "flavors": {{http.StatusOK, `[{"name": "b3-8", "generation": "GEN3", "order": 1, "tags": ["new"], "lifecycle": ` + lifecycle + `,
			"specifications": {"vcpus": 2, "memory": {"unit": "GB", "value": 8}}}]`}},
		"GET " + base + "options": {{http.StatusOK, `[{"name": "deletion_protection", "type": "BOOLEAN"}]`}},
		"GET " + base + "plans":   {{http.StatusOK, `[{"name": "production", "description": "Production", "backupRetention": "P14D", "order": 2, "tags": null, "lifecycle": ` + lifecycle + `}]`}},
		"GET " + base + "regions": {{http.StatusOK, `[{"name": "eu-west-par", "type": "REGION-3-AZ", "order": 1, "tags": null, "lifecycle": ` + lifecycle + `}]`}},
	})

	var got cloudManagedDatabaseCapabilitiesModel
	if diags := state.Get(context.Background(), &got); diags.HasError() {
		t.Fatalf("unexpected diagnostics: %v", diags)
	}
	if len(got.Engines) != 1 || got.Engines[0].Versions[0].Name != "17" || got.Flavors[0].Specifications.Memory.Value != 8 ||
		got.Plans[0].Name != "production" || got.Regions[0].Type != "REGION-3-AZ" || got.Disks[0].IopsPerGigabyte != nil || got.Options[0].Type != "BOOLEAN" {
		t.Errorf("unexpected capabilities: %+v", got)
	}
}

func TestAccDataSourceCloudManagedDatabaseCapabilities_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckCloudManagedDatabasePostgresql(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
data "ovh_cloud_managed_database_capabilities" "capabilities" {
  service_name = %q
}
`, os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.ovh_cloud_managed_database_capabilities.capabilities", "engines.0.name"),
					resource.TestCheckResourceAttrSet("data.ovh_cloud_managed_database_capabilities.capabilities", "flavors.0.name"),
					resource.TestCheckResourceAttrSet("data.ovh_cloud_managed_database_capabilities.capabilities", "plans.0.name"),
					resource.TestCheckResourceAttrSet("data.ovh_cloud_managed_database_capabilities.capabilities", "regions.0.name"),
				),
			},
		},
	})
}

func TestAccDataSourceCloudManagedAnalyticsCapabilities_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCloud(t)
			testAccCheckCloudProjectExists(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
data "ovh_cloud_managed_analytics_capabilities" "capabilities" {
  service_name = %q
}
`, os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.ovh_cloud_managed_analytics_capabilities.capabilities", "engines.0.name"),
					resource.TestCheckResourceAttrSet("data.ovh_cloud_managed_analytics_capabilities.capabilities", "flavors.0.name"),
					resource.TestCheckResourceAttrSet("data.ovh_cloud_managed_analytics_capabilities.capabilities", "plans.0.name"),
					resource.TestCheckResourceAttrSet("data.ovh_cloud_managed_analytics_capabilities.capabilities", "regions.0.name"),
				),
			},
		},
	})
}
