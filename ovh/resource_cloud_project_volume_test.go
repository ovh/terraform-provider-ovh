package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCloudProjectVolume_basic(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	regionName := os.Getenv("OVH_CLOUD_PROJECT_REGION_TEST")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { testAccPreCheckCloud(t); testAccCheckCloudProjectExists(t) },
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "ovh_cloud_project_volume" "volume"  {
						region_name = "%s"
						service_name = "%s"
						description = "test"
						name = "test"
						size = 15
						type = "classic"
					}
				`,
					regionName,
					serviceName,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_project_volume.volume", "region_name", regionName),
					resource.TestCheckResourceAttr("ovh_cloud_project_volume.volume", "service_name", serviceName),
					resource.TestCheckResourceAttrSet("ovh_cloud_project_volume.volume", "volume_id"),
					resource.TestCheckResourceAttr("ovh_cloud_project_volume.volume", "type", "classic"),
					resource.TestCheckResourceAttr("ovh_cloud_project_volume.volume", "description", "test"),
					resource.TestCheckResourceAttr("ovh_cloud_project_volume.volume", "name", "test"),
					resource.TestCheckResourceAttr("ovh_cloud_project_volume.volume", "size", "15"),
				),
			},
			{
				Config: fmt.Sprintf(`
					resource "ovh_cloud_project_volume" "volume"  {
						region_name = "%s"
						service_name = "%s"
						description = "test_updated"
						name = "test_updated"
						size = 20
						type = "classic"
					}
				`,
					regionName,
					serviceName,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_project_volume.volume", "region_name", regionName),
					resource.TestCheckResourceAttr("ovh_cloud_project_volume.volume", "service_name", serviceName),
					resource.TestCheckResourceAttrSet("ovh_cloud_project_volume.volume", "volume_id"),
					resource.TestCheckResourceAttr("ovh_cloud_project_volume.volume", "type", "classic"),
					resource.TestCheckResourceAttr("ovh_cloud_project_volume.volume", "description", "test_updated"),
					resource.TestCheckResourceAttr("ovh_cloud_project_volume.volume", "name", "test_updated"),
					resource.TestCheckResourceAttr("ovh_cloud_project_volume.volume", "size", "20"),
				),
			},
			{
				ResourceName:            "ovh_cloud_project_volume.volume",
				ImportStateIdPrefix:     serviceName + "/" + regionName + "/",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"type", "description"},
			},
		},
	})
}

// TestAccCloudProjectVolume_encryptionCMK exercises customer managed key (CMK)
// encryption. It requires a region that supports CMK (preprod or 3AZ) and a
// valid OKMS domain/service key, so it only runs when the CMK env vars are set.
func TestAccCloudProjectVolume_encryptionCMK(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	regionName := os.Getenv("OVH_CLOUD_PROJECT_REGION_CMK_TEST")
	kmsDomainID := os.Getenv("OVH_CLOUD_PROJECT_KMS_DOMAIN_ID_TEST")
	kmsServiceKeyID := os.Getenv("OVH_CLOUD_PROJECT_KMS_SERVICE_KEY_ID_TEST")
	// Optional: required when the target region is a 3AZ region.
	availabilityZone := os.Getenv("OVH_CLOUD_PROJECT_AZ_TEST")

	availabilityZoneConfig := ""
	if availabilityZone != "" {
		availabilityZoneConfig = fmt.Sprintf("availability_zone = %q", availabilityZone)
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck: func() {
			testAccPreCheckCloud(t)
			testAccCheckCloudProjectExists(t)
			if regionName == "" || kmsDomainID == "" || kmsServiceKeyID == "" {
				t.Skip("OVH_CLOUD_PROJECT_REGION_CMK_TEST, OVH_CLOUD_PROJECT_KMS_DOMAIN_ID_TEST and OVH_CLOUD_PROJECT_KMS_SERVICE_KEY_ID_TEST must be set for CMK acceptance test")
			}
		},
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "ovh_cloud_project_volume" "cmk_volume" {
						region_name  = "%s"
						service_name = "%s"
						description  = "test_cmk"
						name         = "test_cmk"
						size         = 15
						type         = "high-speed-gen2"
						%s

						encryption = {
							encrypted = true
							kms = {
								domain_id      = "%s"
								service_key_id = "%s"
							}
						}
					}
				`,
					regionName,
					serviceName,
					availabilityZoneConfig,
					kmsDomainID,
					kmsServiceKeyID,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_project_volume.cmk_volume", "region_name", regionName),
					resource.TestCheckResourceAttr("ovh_cloud_project_volume.cmk_volume", "service_name", serviceName),
					resource.TestCheckResourceAttrSet("ovh_cloud_project_volume.cmk_volume", "volume_id"),
					resource.TestCheckResourceAttr("ovh_cloud_project_volume.cmk_volume", "encryption.encrypted", "true"),
					resource.TestCheckResourceAttr("ovh_cloud_project_volume.cmk_volume", "encryption.kms.domain_id", kmsDomainID),
					resource.TestCheckResourceAttr("ovh_cloud_project_volume.cmk_volume", "encryption.kms.service_key_id", kmsServiceKeyID),
				),
			},
		},
	})
}
