package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestAccCloudQuotaResource_basic exercises PUT /v2/publicCloud/project/{id}/quota
// through the ovh_cloud_quota resource: apply a profile in a region, then
// flip to a second profile to validate the update path.
//
// Required env:
//   OVH_CLOUD_PROJECT_SERVICE_TEST        — project id
//   OVH_CLOUD_PROJECT_QUOTA_REGION_TEST   — region to mutate (e.g. GRA11)
//   OVH_CLOUD_PROJECT_QUOTA_PROFILE_TEST  — initial profile name
//   OVH_CLOUD_PROJECT_QUOTA_PROFILE2_TEST — second profile name (update step)
//
// Both profiles must exist in current_state.available_profiles for the
// project. The test skips if any of these are unset.
func TestAccCloudQuotaResource_basic(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	region := os.Getenv("OVH_CLOUD_PROJECT_QUOTA_REGION_TEST")
	profile1 := os.Getenv("OVH_CLOUD_PROJECT_QUOTA_PROFILE_TEST")
	profile2 := os.Getenv("OVH_CLOUD_PROJECT_QUOTA_PROFILE2_TEST")

	configFor := func(profile string) string {
		return fmt.Sprintf(`
resource "ovh_cloud_quota" "quota" {
  service_name                    = "%s"
  prevent_automatic_quota_upgrade = false
  regions = [
    {
      region  = "%s"
      profile = "%s"
    },
  ]
}
`, serviceName, region, profile)
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckCloudQuotaResource(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: configFor(profile1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"ovh_cloud_quota.quota", "service_name", serviceName),
					resource.TestCheckResourceAttrSet(
						"ovh_cloud_quota.quota", "id"),
					resource.TestCheckResourceAttrSet(
						"ovh_cloud_quota.quota", "checksum"),
					resource.TestCheckResourceAttr(
						"ovh_cloud_quota.quota", "prevent_automatic_quota_upgrade", "false"),
					resource.TestCheckResourceAttr(
						"ovh_cloud_quota.quota", "regions.#", "1"),
					resource.TestCheckResourceAttr(
						"ovh_cloud_quota.quota", "regions.0.region", region),
					resource.TestCheckResourceAttr(
						"ovh_cloud_quota.quota", "regions.0.profile", profile1),
					resource.TestCheckResourceAttrSet(
						"ovh_cloud_quota.quota", "current_state.regions.#"),
				),
			},
			{
				Config: configFor(profile2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"ovh_cloud_quota.quota", "regions.0.region", region),
					resource.TestCheckResourceAttr(
						"ovh_cloud_quota.quota", "regions.0.profile", profile2),
					resource.TestCheckResourceAttrSet(
						"ovh_cloud_quota.quota", "checksum"),
				),
			},
			{
				ResourceName:      "ovh_cloud_quota.quota",
				ImportState:       true,
				ImportStateId:     serviceName,
				ImportStateVerify: true,
				// current_state is computed and may drift between PUT-poll and
				// the post-import read; checksum/updated_at can also shift.
				ImportStateVerifyIgnore: []string{
					"checksum",
					"updated_at",
					"current_state",
				},
			},
		},
	})
}

func testAccPreCheckCloudQuotaResource(t *testing.T) {
	testAccPreCheckCredentials(t)
	for _, v := range []string{
		"OVH_CLOUD_PROJECT_SERVICE_TEST",
		"OVH_CLOUD_PROJECT_QUOTA_REGION_TEST",
		"OVH_CLOUD_PROJECT_QUOTA_PROFILE_TEST",
		"OVH_CLOUD_PROJECT_QUOTA_PROFILE2_TEST",
	} {
		if os.Getenv(v) == "" {
			t.Skipf("%s not set", v)
		}
	}
}
