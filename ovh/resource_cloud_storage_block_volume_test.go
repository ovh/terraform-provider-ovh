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

// TestAccCloudStorageBlockVolume_encryptionCMK verifies that a volume can be created
// with a customer-managed key (CMK) via encryption.kms. It requires a valid OKMS
// domain/service key, so it only runs when the CMK env vars are set.
func TestAccCloudStorageBlockVolume_encryptionCMK(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	region := os.Getenv("OVH_CLOUD_PROJECT_REGION_TEST")
	kmsDomainID := os.Getenv("OVH_CLOUD_PROJECT_KMS_DOMAIN_ID_TEST")
	kmsServiceKeyID := os.Getenv("OVH_CLOUD_PROJECT_KMS_SERVICE_KEY_ID_TEST")

	volumeName := acctest.RandomWithPrefix(testAccResourceCloudStorageBlockVolumeNamePrefix)

	config := fmt.Sprintf(`
resource "ovh_cloud_storage_block_volume" "cmk_volume" {
  service_name = "%s"
  name         = "%s"
  size         = 10
  region       = "%s"
  volume_type  = "HIGH_SPEED"

  encryption = {
    enabled = true
    kms = {
      domain_id      = "%s"
      service_key_id = "%s"
    }
  }
}
`, serviceName, volumeName, region, kmsDomainID, kmsServiceKeyID)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCloud(t)
			testAccCheckCloudProjectExists(t)
			if kmsDomainID == "" || kmsServiceKeyID == "" {
				t.Skip("OVH_CLOUD_PROJECT_KMS_DOMAIN_ID_TEST and OVH_CLOUD_PROJECT_KMS_SERVICE_KEY_ID_TEST must be set for CMK acceptance test")
			}
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_storage_block_volume.cmk_volume", "service_name", serviceName),
					resource.TestCheckResourceAttr("ovh_cloud_storage_block_volume.cmk_volume", "name", volumeName),
					resource.TestCheckResourceAttr("ovh_cloud_storage_block_volume.cmk_volume", "region", region),
					resource.TestCheckResourceAttr("ovh_cloud_storage_block_volume.cmk_volume", "encryption.enabled", "true"),
					resource.TestCheckResourceAttr("ovh_cloud_storage_block_volume.cmk_volume", "encryption.kms.domain_id", kmsDomainID),
					resource.TestCheckResourceAttr("ovh_cloud_storage_block_volume.cmk_volume", "encryption.kms.service_key_id", kmsServiceKeyID),
					resource.TestCheckResourceAttr("ovh_cloud_storage_block_volume.cmk_volume", "current_state.encryption.enabled", "true"),
					resource.TestCheckResourceAttr("ovh_cloud_storage_block_volume.cmk_volume", "current_state.encryption.kms.domain_id", kmsDomainID),
					resource.TestCheckResourceAttr("ovh_cloud_storage_block_volume.cmk_volume", "current_state.encryption.kms.service_key_id", kmsServiceKeyID),
					resource.TestCheckResourceAttrSet("ovh_cloud_storage_block_volume.cmk_volume", "id"),
					resource.TestCheckResourceAttrSet("ovh_cloud_storage_block_volume.cmk_volume", "checksum"),
					resource.TestCheckResourceAttr("ovh_cloud_storage_block_volume.cmk_volume", "resource_status", "READY"),
				),
			},
			// Test import
			{
				ResourceName:      "ovh_cloud_storage_block_volume.cmk_volume",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccCloudStorageBlockVolumeImportStateIdFunc("ovh_cloud_storage_block_volume.cmk_volume"),
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

// TestAccCloudStorageBlockVolume_availabilityZone verifies that a volume can be
// created pinned to a specific availability zone via the root-level
// `availability_zone` attribute, and that the AZ is reflected both at the root
// and inside current_state.location.
func TestAccCloudStorageBlockVolume_availabilityZone(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")

	// Per explicit instruction, these values are hardcoded for this test.
	const (
		region           = "EU-WEST-PAR"
		availabilityZone = "eu-west-par-a"
	)

	volumeName := acctest.RandomWithPrefix(testAccResourceCloudStorageBlockVolumeNamePrefix)

	config := fmt.Sprintf(`
resource "ovh_cloud_storage_block_volume" "volume" {
  service_name      = "%s"
  name              = "%s"
  size              = 10
  region            = "%s"
  availability_zone = "%s"
  volume_type       = "HIGH_SPEED_GEN2"

  encryption = {
    enabled = false
  }
}
`, serviceName, volumeName, region, availabilityZone)

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
					resource.TestCheckResourceAttr("ovh_cloud_storage_block_volume.volume", "region", region),
					resource.TestCheckResourceAttr("ovh_cloud_storage_block_volume.volume", "availability_zone", availabilityZone),
					resource.TestCheckResourceAttr("ovh_cloud_storage_block_volume.volume", "current_state.location.availability_zone", availabilityZone),
					resource.TestCheckResourceAttrSet("ovh_cloud_storage_block_volume.volume", "id"),
					resource.TestCheckResourceAttrSet("ovh_cloud_storage_block_volume.volume", "checksum"),
					resource.TestCheckResourceAttr("ovh_cloud_storage_block_volume.volume", "resource_status", "READY"),
				),
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
