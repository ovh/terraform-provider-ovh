package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCloudProjectQuotaDataSource_basic(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")

	config := fmt.Sprintf(`
data "ovh_cloud_project_quota" "quota" {
  service_name = "%s"
}
`, serviceName)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCloudProjectQuota(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.ovh_cloud_project_quota.quota", "service_name", serviceName),
					resource.TestCheckResourceAttrSet(
						"data.ovh_cloud_project_quota.quota", "id"),
					resource.TestCheckResourceAttrSet(
						"data.ovh_cloud_project_quota.quota", "current_state.profile"),
					resource.TestCheckResourceAttrSet(
						"data.ovh_cloud_project_quota.quota", "current_state.available_profiles.#"),
				),
			},
		},
	})
}

func testAccPreCheckCloudProjectQuota(t *testing.T) {
	testAccPreCheckCredentials(t)
	if os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST") == "" {
		t.Skip("OVH_CLOUD_PROJECT_SERVICE_TEST not set")
	}
}
