package ovh

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccCloudStorageBlockVolume_basic(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	region := os.Getenv("OVH_CLOUD_PROJECT_REGION_TEST")

	volumeName := acctest.RandomWithPrefix(testAccResourceCloudStorageBlockVolumeNamePrefix)

	config := fmt.Sprintf(`
resource "ovh_cloud_storage_block_volume" "volume" {
  service_name = "%s"
  name         = "%s"
  size         = 10
  region       = "%s"
  volume_type  = "CLASSIC"

  encryption = {
    enabled = false
  }
}
`, serviceName, volumeName, region)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCloud(t)
			testAccCheckCloudProjectExists(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_storage_block_volume.volume", "service_name", serviceName),
					resource.TestCheckResourceAttr("ovh_cloud_storage_block_volume.volume", "name", volumeName),
					resource.TestCheckResourceAttr("ovh_cloud_storage_block_volume.volume", "size", "10"),
					resource.TestCheckResourceAttr("ovh_cloud_storage_block_volume.volume", "region", region),
					resource.TestCheckResourceAttr("ovh_cloud_storage_block_volume.volume", "volume_type", "CLASSIC"),
					resource.TestCheckResourceAttr("ovh_cloud_storage_block_volume.volume", "encryption.enabled", "false"),
					resource.TestCheckResourceAttr("ovh_cloud_storage_block_volume.volume", "current_state.encryption.enabled", "false"),
					resource.TestCheckResourceAttrSet("ovh_cloud_storage_block_volume.volume", "id"),
					resource.TestCheckResourceAttrSet("ovh_cloud_storage_block_volume.volume", "checksum"),
					resource.TestCheckResourceAttrSet("ovh_cloud_storage_block_volume.volume", "created_at"),
					resource.TestCheckResourceAttrSet("ovh_cloud_storage_block_volume.volume", "resource_status"),
					resource.TestCheckResourceAttrSet("ovh_cloud_storage_block_volume.volume", "current_state.volume_type"),
					resource.TestCheckResourceAttrSet("ovh_cloud_storage_block_volume.volume", "current_state.status"),
				),
			},
			// Test import
			{
				ResourceName:      "ovh_cloud_storage_block_volume.volume",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccCloudStorageBlockVolumeImportStateIdFunc("ovh_cloud_storage_block_volume.volume"),
			},
		},
	})
}

func TestAccCloudStorageBlockVolume_update(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	region := os.Getenv("OVH_CLOUD_PROJECT_REGION_TEST")

	volumeName := acctest.RandomWithPrefix(testAccResourceCloudStorageBlockVolumeNamePrefix)
	updatedName := acctest.RandomWithPrefix(testAccResourceCloudStorageBlockVolumeNamePrefix)

	config := fmt.Sprintf(`
resource "ovh_cloud_storage_block_volume" "volume" {
  service_name = "%s"
  name         = "%s"
  size         = 10
  region       = "%s"
  volume_type  = "CLASSIC"

  encryption = {
    enabled = false
  }
}
`, serviceName, volumeName, region)

	updatedConfig := fmt.Sprintf(`
resource "ovh_cloud_storage_block_volume" "volume" {
  service_name = "%s"
  name         = "%s"
  size         = 20
  region       = "%s"
  volume_type  = "CLASSIC"

  encryption = {
    enabled = false
  }
}
`, serviceName, updatedName, region)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCloud(t)
			testAccCheckCloudProjectExists(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_storage_block_volume.volume", "name", volumeName),
					resource.TestCheckResourceAttr("ovh_cloud_storage_block_volume.volume", "size", "10"),
					resource.TestCheckResourceAttr("ovh_cloud_storage_block_volume.volume", "encryption.enabled", "false"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_storage_block_volume.volume", "name", updatedName),
					resource.TestCheckResourceAttr("ovh_cloud_storage_block_volume.volume", "size", "20"),
					resource.TestCheckResourceAttr("ovh_cloud_storage_block_volume.volume", "encryption.enabled", "false"),
					resource.TestCheckResourceAttrSet("ovh_cloud_storage_block_volume.volume", "id"),
					resource.TestCheckResourceAttrSet("ovh_cloud_storage_block_volume.volume", "checksum"),
				),
			},
		},
	})
}

// TestAccCloudStorageBlockVolume_encryptionToggle verifies that toggling
// `encryption.enabled` recreates the volume.
func TestAccCloudStorageBlockVolume_encryptionToggle(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	region := os.Getenv("OVH_CLOUD_PROJECT_REGION_TEST")

	volumeName := acctest.RandomWithPrefix(testAccResourceCloudStorageBlockVolumeNamePrefix)

	plainConfig := fmt.Sprintf(`
resource "ovh_cloud_storage_block_volume" "volume" {
  service_name = "%s"
  name         = "%s"
  size         = 10
  region       = "%s"
  volume_type  = "CLASSIC"

  encryption = {
    enabled = false
  }
}
`, serviceName, volumeName, region)

	encryptedConfig := fmt.Sprintf(`
resource "ovh_cloud_storage_block_volume" "volume" {
  service_name = "%s"
  name         = "%s"
  size         = 10
  region       = "%s"
  volume_type  = "CLASSIC"

  encryption = {
    enabled = true
  }
}
`, serviceName, volumeName, region)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCloud(t)
			testAccCheckCloudProjectExists(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: plainConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_storage_block_volume.volume", "volume_type", "CLASSIC"),
					resource.TestCheckResourceAttr("ovh_cloud_storage_block_volume.volume", "encryption.enabled", "false"),
					resource.TestCheckResourceAttr("ovh_cloud_storage_block_volume.volume", "current_state.encryption.enabled", "false"),
					resource.TestCheckResourceAttrSet("ovh_cloud_storage_block_volume.volume", "id"),
				),
			},
			{
				// Toggle encryption on — must replace the resource.
				Config: encryptedConfig,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectNonEmptyPlan(),
						plancheck.ExpectResourceAction(
							"ovh_cloud_storage_block_volume.volume",
							plancheck.ResourceActionReplace,
						),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_storage_block_volume.volume", "volume_type", "CLASSIC"),
					resource.TestCheckResourceAttr("ovh_cloud_storage_block_volume.volume", "encryption.enabled", "true"),
					resource.TestCheckResourceAttr("ovh_cloud_storage_block_volume.volume", "current_state.encryption.enabled", "true"),
					resource.TestCheckResourceAttrSet("ovh_cloud_storage_block_volume.volume", "id"),
					resource.TestCheckResourceAttrSet("ovh_cloud_storage_block_volume.volume", "checksum"),
				),
			},
			{
				// Toggle encryption back off — also replaces the resource.
				Config: plainConfig,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectNonEmptyPlan(),
						plancheck.ExpectResourceAction(
							"ovh_cloud_storage_block_volume.volume",
							plancheck.ResourceActionReplace,
						),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_storage_block_volume.volume", "encryption.enabled", "false"),
					resource.TestCheckResourceAttr("ovh_cloud_storage_block_volume.volume", "current_state.encryption.enabled", "false"),
				),
			},
		},
	})
}

// TestAccCloudStorageBlockVolume_createFromImage verifies that a volume can be
// created from a Glance image using the create_from.image_id field.
func TestAccCloudStorageBlockVolume_createFromImage(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	region := os.Getenv("OVH_CLOUD_PROJECT_REGION_TEST")
	imageID := os.Getenv("OVH_CLOUD_PROJECT_GLANCE_IMAGE_ID_TEST")

	if imageID == "" {
		t.Skip("OVH_CLOUD_PROJECT_GLANCE_IMAGE_ID_TEST is not set; skipping create_from.image_id test")
	}

	volumeName := acctest.RandomWithPrefix(testAccResourceCloudStorageBlockVolumeNamePrefix)

	config := fmt.Sprintf(`
resource "ovh_cloud_storage_block_volume" "volume_from_image" {
  service_name = "%s"
  name         = "%s"
  size         = 10
  region       = "%s"
  volume_type  = "CLASSIC"

  create_from = {
    image_id = "%s"
  }
}
`, serviceName, volumeName, region, imageID)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCloud(t)
			testAccCheckCloudProjectExists(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_storage_block_volume.volume_from_image", "service_name", serviceName),
					resource.TestCheckResourceAttr("ovh_cloud_storage_block_volume.volume_from_image", "name", volumeName),
					resource.TestCheckResourceAttr("ovh_cloud_storage_block_volume.volume_from_image", "size", "10"),
					resource.TestCheckResourceAttr("ovh_cloud_storage_block_volume.volume_from_image", "region", region),
					resource.TestCheckResourceAttr("ovh_cloud_storage_block_volume.volume_from_image", "volume_type", "CLASSIC"),
					resource.TestCheckResourceAttr("ovh_cloud_storage_block_volume.volume_from_image", "create_from.image_id", imageID),
					resource.TestCheckResourceAttr("ovh_cloud_storage_block_volume.volume_from_image", "current_state.bootable", "true"),
					resource.TestCheckResourceAttrSet("ovh_cloud_storage_block_volume.volume_from_image", "id"),
					resource.TestCheckResourceAttrSet("ovh_cloud_storage_block_volume.volume_from_image", "checksum"),
					resource.TestCheckResourceAttr("ovh_cloud_storage_block_volume.volume_from_image", "resource_status", "READY"),
				),
			},
			{
				ResourceName:      "ovh_cloud_storage_block_volume.volume_from_image",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccCloudStorageBlockVolumeImportStateIdFunc("ovh_cloud_storage_block_volume.volume_from_image"),
			},
		},
	})
}

// TestAccCloudStorageBlockVolume_recoverDefaults verifies that a volume created
// WITHOUT an `encryption` block does not produce a perpetual diff.
func TestAccCloudStorageBlockVolume_recoverDefaults(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	region := os.Getenv("OVH_CLOUD_PROJECT_REGION_TEST")

	volumeName := acctest.RandomWithPrefix(testAccResourceCloudStorageBlockVolumeNamePrefix)

	// Config omits `encryption` block entirely so that the
	// recover step is responsible for populating both.
	config := fmt.Sprintf(`
resource "ovh_cloud_storage_block_volume" "volume" {
  service_name = "%s"
  name         = "%s"
  size         = 10
  region       = "%s"
  volume_type  = "CLASSIC"
}
`, serviceName, volumeName, region)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCloud(t)
			testAccCheckCloudProjectExists(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_storage_block_volume.volume", "service_name", serviceName),
					resource.TestCheckResourceAttr("ovh_cloud_storage_block_volume.volume", "name", volumeName),
					resource.TestCheckResourceAttr("ovh_cloud_storage_block_volume.volume", "size", "10"),
					resource.TestCheckResourceAttr("ovh_cloud_storage_block_volume.volume", "region", region),
					// volume_type was omitted; recover must have populated it.
					resource.TestCheckResourceAttrSet("ovh_cloud_storage_block_volume.volume", "volume_type"),
					resource.TestCheckResourceAttrSet("ovh_cloud_storage_block_volume.volume", "current_state.volume_type"),
					// encryption block was omitted; recover must have populated enabled.
					resource.TestCheckResourceAttrSet("ovh_cloud_storage_block_volume.volume", "encryption.enabled"),
					resource.TestCheckResourceAttrSet("ovh_cloud_storage_block_volume.volume", "id"),
					resource.TestCheckResourceAttrSet("ovh_cloud_storage_block_volume.volume", "checksum"),
					resource.TestCheckResourceAttrSet("ovh_cloud_storage_block_volume.volume", "resource_status"),
				),
			},
			{
				// Re-applying the SAME config must yield an empty plan. This is the
				// core assertion: the recover-populated volume_type and encryption
				// must not cause an "inconsistent result" or a perpetual diff.
				Config: config,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
			{
				ResourceName:      "ovh_cloud_storage_block_volume.volume",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccCloudStorageBlockVolumeImportStateIdFunc("ovh_cloud_storage_block_volume.volume"),
			},
		},
	})
}

// TestAccCloudStorageBlockVolume_noServiceName validates that when service_name
// is omitted from the resource configuration, the provider falls back to the
// OVH_CLOUD_PROJECT_SERVICE environment variable at plan time (via the
// EnvDefaultString plan modifier) and that this does not produce a perpetual
// diff / phantom replace on subsequent plans.
func TestAccCloudStorageBlockVolume_noServiceName(t *testing.T) {
	region := os.Getenv("OVH_CLOUD_PROJECT_REGION_TEST")

	volumeName := acctest.RandomWithPrefix(testAccResourceCloudStorageBlockVolumeNamePrefix)

	config := fmt.Sprintf(`
resource "ovh_cloud_storage_block_volume" "volume" {
  name        = "%s"
  size        = 10
  region      = "%s"
  volume_type = "CLASSIC"

  encryption = {
    enabled = false
  }
}
`, volumeName, region)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCloud(t)
			testAccCheckCloudProjectExists(t)
			checkEnvOrSkip(t, "OVH_CLOUD_PROJECT_SERVICE")
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_storage_block_volume.volume", "name", volumeName),
					resource.TestCheckResourceAttr("ovh_cloud_storage_block_volume.volume", "service_name", os.Getenv("OVH_CLOUD_PROJECT_SERVICE")),
					resource.TestCheckResourceAttrSet("ovh_cloud_storage_block_volume.volume", "id"),
				),
			},
			{
				// Re-planning with service_name still omitted must be a no-op:
				// the EnvDefaultString modifier injects the env value so the plan
				// matches state and RequiresReplace does not fire (regression guard).
				Config:   config,
				PlanOnly: true,
			},
		},
	})
}

// TestAccCloudStorageBlockVolume_missingServiceName validates that when
// service_name is omitted from the configuration AND the
// OVH_CLOUD_PROJECT_SERVICE environment variable is unset, the EnvDefaultString
// plan modifier (required: true) raises a "Missing" diagnostic at plan time.
func TestAccCloudStorageBlockVolume_missingServiceName(t *testing.T) {
	region := os.Getenv("OVH_CLOUD_PROJECT_REGION_TEST")

	volumeName := acctest.RandomWithPrefix(testAccResourceCloudStorageBlockVolumeNamePrefix)

	config := fmt.Sprintf(`
resource "ovh_cloud_storage_block_volume" "volume" {
  name        = "%s"
  size        = 10
  region      = "%s"
  volume_type = "CLASSIC"

  encryption = {
    enabled = false
  }
}
`, volumeName, region)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCloud(t)
			// Ensure the env var is unset for the duration of this test so the
			// plan modifier cannot resolve service_name from anywhere.
			t.Setenv("OVH_CLOUD_PROJECT_SERVICE", "")
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      config,
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`Missing`),
			},
		},
	})
}

const testAccResourceCloudStorageBlockVolumeNamePrefix = "tf-test-volume-v2-"

func testAccCloudStorageBlockVolumeImportStateIdFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}
		return fmt.Sprintf("%s/%s", rs.Primary.Attributes["service_name"], rs.Primary.Attributes["id"]), nil
	}
}
