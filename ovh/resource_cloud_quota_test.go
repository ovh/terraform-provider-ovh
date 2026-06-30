package ovh

import (
	"fmt"
	"net/url"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestAccCloudQuotaResource_basic exercises PUT /v2/publicCloud/project/{id}/quota
// through the ovh_cloud_quota resource: apply a profile in a region, flip to a
// second profile to validate the update path, then toggle
// prevent_automatic_quota_upgrade. Finally it imports the singleton envelope.
//
// Required env:
//
//	OVH_CLOUD_PROJECT_SERVICE_TEST — project id
//	OVH_CLOUD_PROJECT_REGION_TEST  — region to mutate (e.g. GRA11)
//
// The two profiles exercised by the apply/update steps are read live from
// current_state.available_profiles. The test skips if fewer than 2 are
// available.
func TestAccCloudQuotaResource_basic(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	region := os.Getenv("OVH_CLOUD_PROJECT_REGION_TEST")

	testAccPreCheckCloudQuotaResource(t)

	endpoint := "/v2/publicCloud/project/" + url.PathEscape(serviceName) + "/quota"
	var resp CloudQuotaAPIResponse
	if err := testAccOVHClient.Get(endpoint, &resp); err != nil {
		t.Fatalf("failed to fetch quota envelope from %s: %s", endpoint, err)
	}

	var names []string
	seen := map[string]bool{}
	if resp.CurrentState != nil {
		for _, p := range resp.CurrentState.AvailableProfiles {
			if p.Name != "" && !seen[p.Name] {
				seen[p.Name] = true
				names = append(names, p.Name)
			}
		}
	}

	// Profile currently applied on the target region. profile1 must differ from
	// it so the first apply performs a real quota change instead of re-applying
	// the same profile (a silent no-op).
	currentProfile := ""
	if resp.CurrentState != nil {
		for _, r := range resp.CurrentState.Regions {
			if r.Location != nil && r.Location.Region == region {
				currentProfile = r.Profile
				break
			}
		}
	}

	profile1 := ""
	for _, n := range names {
		if n != currentProfile {
			profile1 = n
			break
		}
	}
	profile2 := ""
	for _, n := range names {
		if n != profile1 {
			profile2 = n
			break
		}
	}
	if profile1 == "" || profile2 == "" {
		t.Skipf("need 2 distinct quota profiles with profile1 != current (%q) on %s to exercise a real switch; available=%v", currentProfile, region, names)
	}

	configFor := func(profile string, preventUpgrade bool) string {
		return fmt.Sprintf(`
resource "ovh_cloud_quota" "quota" {
  service_name                    = "%s"
  prevent_automatic_quota_upgrade = %t
  regions = [
    {
      region  = "%s"
      profile = "%s"
    },
  ]
}
`, serviceName, preventUpgrade, region, profile)
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckCloudQuotaResource(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: configFor(profile1, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"ovh_cloud_quota.quota", "service_name", serviceName),
					resource.TestCheckResourceAttrSet(
						"ovh_cloud_quota.quota", "id"),
					resource.TestCheckResourceAttrSet(
						"ovh_cloud_quota.quota", "checksum"),
					resource.TestCheckResourceAttr(
						"ovh_cloud_quota.quota", "resource_status", "READY"),
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
				// Update path: switch the applied profile. The switch must fully
				// reconcile — the resource has to land in READY, not OUT_OF_SYNC.
				Config: configFor(profile2, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"ovh_cloud_quota.quota", "regions.0.region", region),
					resource.TestCheckResourceAttr(
						"ovh_cloud_quota.quota", "regions.0.profile", profile2),
					resource.TestCheckResourceAttr(
						"ovh_cloud_quota.quota", "resource_status", "READY"),
					resource.TestCheckResourceAttrSet(
						"ovh_cloud_quota.quota", "checksum"),
				),
			},
			{
				// Toggle the automatic-quota-upgrade switch.
				Config: configFor(profile2, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"ovh_cloud_quota.quota", "prevent_automatic_quota_upgrade", "true"),
				),
			},
			{
				ResourceName:      "ovh_cloud_quota.quota",
				ImportState:       true,
				ImportStateId:     serviceName,
				ImportStateVerify: true,
				// Import reads the singleton envelope, whose targetSpec lists every
				// region of the project — not just the subset declared in config —
				// so the managed `regions` list legitimately differs after import.
				// checksum/updated_at/current_state come from the same GET as the
				// apply and do not move between the two reads, so they are verified.
				ImportStateVerifyIgnore: []string{
					"regions",
				},
			},
		},
	})
}

func testAccPreCheckCloudQuotaResource(t *testing.T) {
	testAccPreCheckCredentials(t)
	for _, v := range []string{
		"OVH_CLOUD_PROJECT_SERVICE_TEST",
		"OVH_CLOUD_PROJECT_REGION_TEST",
	} {
		if os.Getenv(v) == "" {
			t.Skipf("%s not set", v)
		}
	}
}
