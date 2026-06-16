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

// TestAccCloudStorageBlockVolume_createFromBackup verifies volume_type handling
// when restoring from a backup: an explicit value is honored by the API, and an
// omitted value is inferred from the source volume.
func TestAccCloudStorageBlockVolume_createFromBackup(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	region := os.Getenv("OVH_CLOUD_PROJECT_REGION_TEST")

	volumeName := acctest.RandomWithPrefix(testAccResourceCloudStorageBlockVolumeNamePrefix)

	config := fmt.Sprintf(`
resource "ovh_cloud_storage_block_volume" "source" {
  service_name = "%[1]s"
  name         = "%[2]s-source"
  size         = 10
  region       = "%[3]s"
  volume_type  = "CLASSIC"
}

resource "ovh_cloud_storage_block_volume_backup" "backup" {
  service_name = "%[1]s"
  region       = "%[3]s"
  volume_id    = ovh_cloud_storage_block_volume.source.id
  name         = "%[2]s-backup"
}

resource "ovh_cloud_storage_block_volume" "restored_explicit_type" {
  service_name = "%[1]s"
  name         = "%[2]s-restored-explicit"
  size         = 10
  region       = "%[3]s"
  volume_type  = "HIGH_SPEED"

  create_from = {
    backup_id = ovh_cloud_storage_block_volume_backup.backup.id
  }
}

resource "ovh_cloud_storage_block_volume" "restored_inferred_type" {
  service_name = "%[1]s"
  name         = "%[2]s-restored-inferred"
  size         = 10
  region       = "%[3]s"

  create_from = {
    backup_id = ovh_cloud_storage_block_volume_backup.backup.id
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
					resource.TestCheckResourceAttr("ovh_cloud_storage_block_volume.restored_explicit_type", "volume_type", "HIGH_SPEED"),
					resource.TestCheckResourceAttr("ovh_cloud_storage_block_volume.restored_explicit_type", "current_state.volume_type", "HIGH_SPEED"),
					resource.TestCheckResourceAttrPair("ovh_cloud_storage_block_volume.restored_explicit_type", "create_from.backup_id", "ovh_cloud_storage_block_volume_backup.backup", "id"),
					resource.TestCheckResourceAttr("ovh_cloud_storage_block_volume.restored_explicit_type", "resource_status", "READY"),
					resource.TestCheckResourceAttrSet("ovh_cloud_storage_block_volume.restored_explicit_type", "id"),
					resource.TestCheckResourceAttrSet("ovh_cloud_storage_block_volume.restored_explicit_type", "checksum"),
					resource.TestCheckResourceAttrSet("ovh_cloud_storage_block_volume.restored_inferred_type", "volume_type"),
					resource.TestCheckResourceAttrSet("ovh_cloud_storage_block_volume.restored_inferred_type", "current_state.volume_type"),
					resource.TestCheckResourceAttrPair("ovh_cloud_storage_block_volume.restored_inferred_type", "create_from.backup_id", "ovh_cloud_storage_block_volume_backup.backup", "id"),
					resource.TestCheckResourceAttr("ovh_cloud_storage_block_volume.restored_inferred_type", "resource_status", "READY"),
					resource.TestCheckResourceAttrSet("ovh_cloud_storage_block_volume.restored_inferred_type", "id"),
				),
			},
		},
	})
}

// TestAccCloudStorageBlockVolume_createFromSnapshot verifies volume_type
// handling when creating from a snapshot: an explicit value is honored by the
// API, and an omitted value is inferred from the source volume.
func TestAccCloudStorageBlockVolume_createFromSnapshot(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	region := os.Getenv("OVH_CLOUD_PROJECT_REGION_TEST")

	volumeName := acctest.RandomWithPrefix(testAccResourceCloudStorageBlockVolumeNamePrefix)

	config := fmt.Sprintf(`
resource "ovh_cloud_storage_block_volume" "source" {
  service_name = "%[1]s"
  name         = "%[2]s-source"
  size         = 10
  region       = "%[3]s"
  volume_type  = "CLASSIC"
}

resource "ovh_cloud_storage_block_volume_snapshot" "snapshot" {
  service_name = "%[1]s"
  region       = "%[3]s"
  volume_id    = ovh_cloud_storage_block_volume.source.id
  name         = "%[2]s-snapshot"
}

resource "ovh_cloud_storage_block_volume" "from_snapshot_explicit_type" {
  service_name = "%[1]s"
  name         = "%[2]s-from-snapshot-explicit"
  size         = 10
  region       = "%[3]s"
  volume_type  = "HIGH_SPEED"

  create_from = {
    snapshot_id = ovh_cloud_storage_block_volume_snapshot.snapshot.id
  }
}

resource "ovh_cloud_storage_block_volume" "from_snapshot_inferred_type" {
  service_name = "%[1]s"
  name         = "%[2]s-from-snapshot-inferred"
  size         = 10
  region       = "%[3]s"

  create_from = {
    snapshot_id = ovh_cloud_storage_block_volume_snapshot.snapshot.id
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
					resource.TestCheckResourceAttr("ovh_cloud_storage_block_volume.from_snapshot_explicit_type", "volume_type", "HIGH_SPEED"),
					resource.TestCheckResourceAttr("ovh_cloud_storage_block_volume.from_snapshot_explicit_type", "current_state.volume_type", "HIGH_SPEED"),
					resource.TestCheckResourceAttrPair("ovh_cloud_storage_block_volume.from_snapshot_explicit_type", "create_from.snapshot_id", "ovh_cloud_storage_block_volume_snapshot.snapshot", "id"),
					resource.TestCheckResourceAttr("ovh_cloud_storage_block_volume.from_snapshot_explicit_type", "resource_status", "READY"),
					resource.TestCheckResourceAttrSet("ovh_cloud_storage_block_volume.from_snapshot_explicit_type", "id"),
					resource.TestCheckResourceAttrSet("ovh_cloud_storage_block_volume.from_snapshot_explicit_type", "checksum"),
					resource.TestCheckResourceAttrSet("ovh_cloud_storage_block_volume.from_snapshot_inferred_type", "volume_type"),
					resource.TestCheckResourceAttrSet("ovh_cloud_storage_block_volume.from_snapshot_inferred_type", "current_state.volume_type"),
					resource.TestCheckResourceAttrPair("ovh_cloud_storage_block_volume.from_snapshot_inferred_type", "create_from.snapshot_id", "ovh_cloud_storage_block_volume_snapshot.snapshot", "id"),
					resource.TestCheckResourceAttr("ovh_cloud_storage_block_volume.from_snapshot_inferred_type", "resource_status", "READY"),
					resource.TestCheckResourceAttrSet("ovh_cloud_storage_block_volume.from_snapshot_inferred_type", "id"),
				),
			},
		},
	})
}

// TestAccCloudStorageBlockVolume_retype verifies that changing volume_type is an
// online retype applied in place (no replacement).
func TestAccCloudStorageBlockVolume_retype(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	region := os.Getenv("OVH_CLOUD_PROJECT_REGION_TEST")

	volumeName := acctest.RandomWithPrefix(testAccResourceCloudStorageBlockVolumeNamePrefix)

	config := func(volumeType string) string {
		return fmt.Sprintf(`
resource "ovh_cloud_storage_block_volume" "volume" {
  service_name = "%s"
  name         = "%s"
  size         = 10
  region       = "%s"
  volume_type  = "%s"
}
`, serviceName, volumeName, region, volumeType)
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCloud(t)
			testAccCheckCloudProjectExists(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config("CLASSIC"),
				Check: resource.TestCheckResourceAttr(
					"ovh_cloud_storage_block_volume.volume", "volume_type", "CLASSIC"),
			},
			{
				// Retype must be applied in place, not by replacing the volume.
				Config: config("HIGH_SPEED"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(
							"ovh_cloud_storage_block_volume.volume",
							plancheck.ResourceActionUpdate,
						),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_storage_block_volume.volume", "volume_type", "HIGH_SPEED"),
					resource.TestCheckResourceAttr("ovh_cloud_storage_block_volume.volume", "current_state.volume_type", "HIGH_SPEED"),
					resource.TestCheckResourceAttr("ovh_cloud_storage_block_volume.volume", "current_state.status", "AVAILABLE"),
				),
			},
		},
	})
}

// TestAccCloudStorageBlockVolume_encryptedBackupRestore restores a backup of an
// encrypted volume into a new encrypted volume of matching type. Restoring an
// encrypted backup requires the target type and encryption to match the source.
func TestAccCloudStorageBlockVolume_encryptedBackupRestore(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	region := os.Getenv("OVH_CLOUD_PROJECT_REGION_TEST")

	volumeName := acctest.RandomWithPrefix(testAccResourceCloudStorageBlockVolumeNamePrefix)

	config := fmt.Sprintf(`
resource "ovh_cloud_storage_block_volume" "encrypted_source" {
  service_name = "%[1]s"
  name         = "%[2]s-enc-source"
  size         = 10
  region       = "%[3]s"
  volume_type  = "CLASSIC"

  encryption = {
    enabled = true
  }
}

resource "ovh_cloud_storage_block_volume_backup" "encrypted_backup" {
  service_name = "%[1]s"
  region       = "%[3]s"
  volume_id    = ovh_cloud_storage_block_volume.encrypted_source.id
  name         = "%[2]s-enc-backup"
}

resource "ovh_cloud_storage_block_volume" "from_backup_encrypted" {
  service_name = "%[1]s"
  name         = "%[2]s-from-backup-enc"
  size         = 10
  region       = "%[3]s"
  volume_type  = "CLASSIC"

  encryption = {
    enabled = true
  }

  create_from = {
    backup_id = ovh_cloud_storage_block_volume_backup.encrypted_backup.id
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
					resource.TestCheckResourceAttr("ovh_cloud_storage_block_volume.encrypted_source", "current_state.encryption.enabled", "true"),
					resource.TestCheckResourceAttrSet("ovh_cloud_storage_block_volume_backup.encrypted_backup", "id"),
					resource.TestCheckResourceAttr("ovh_cloud_storage_block_volume.from_backup_encrypted", "current_state.status", "AVAILABLE"),
					resource.TestCheckResourceAttr("ovh_cloud_storage_block_volume.from_backup_encrypted", "current_state.encryption.enabled", "true"),
					resource.TestCheckResourceAttrPair("ovh_cloud_storage_block_volume.from_backup_encrypted", "create_from.backup_id", "ovh_cloud_storage_block_volume_backup.encrypted_backup", "id"),
					resource.TestCheckResourceAttrSet("ovh_cloud_storage_block_volume.from_backup_encrypted", "id"),
				),
			},
		},
	})
}

// testAccBlockVolumeBogusID is a syntactically valid UUID guaranteed not to
// exist in the project, used by the negative tests below.
const testAccBlockVolumeBogusID = "00000000-0000-4000-8000-000000000000"

// testAccCloudStorageBlockVolumeExpectError applies a config and asserts the
// apply (or plan) fails with an error matching errRegex.
func testAccCloudStorageBlockVolumeExpectError(t *testing.T, config, errRegex string) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCloud(t)
			testAccCheckCloudProjectExists(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      config,
				ExpectError: regexp.MustCompile(errRegex),
			},
		},
	})
}

// TestAccCloudStorageBlockVolume_invalidSize asserts a size of 0 is rejected.
func TestAccCloudStorageBlockVolume_invalidSize(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	region := os.Getenv("OVH_CLOUD_PROJECT_REGION_TEST")
	volumeName := acctest.RandomWithPrefix(testAccResourceCloudStorageBlockVolumeNamePrefix)

	testAccCloudStorageBlockVolumeExpectError(t, fmt.Sprintf(`
resource "ovh_cloud_storage_block_volume" "neg" {
  service_name = "%[1]s"
  region       = "%[3]s"
  name         = "%[2]s-size0"
  size         = 0
  volume_type  = "CLASSIC"
}
`, serviceName, volumeName, region), "Error calling Post")
}

// TestAccCloudStorageBlockVolume_invalidTwoCreateFromSources asserts a
// create_from with both backup_id and snapshot_id is rejected.
func TestAccCloudStorageBlockVolume_invalidTwoCreateFromSources(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	region := os.Getenv("OVH_CLOUD_PROJECT_REGION_TEST")
	volumeName := acctest.RandomWithPrefix(testAccResourceCloudStorageBlockVolumeNamePrefix)

	testAccCloudStorageBlockVolumeExpectError(t, fmt.Sprintf(`
resource "ovh_cloud_storage_block_volume" "source" {
  service_name = "%[1]s"
  region       = "%[3]s"
  name         = "%[2]s-src"
  size         = 10
  volume_type  = "CLASSIC"
}

resource "ovh_cloud_storage_block_volume_snapshot" "snap" {
  service_name = "%[1]s"
  region       = "%[3]s"
  volume_id    = ovh_cloud_storage_block_volume.source.id
  name         = "%[2]s-snap"
}

resource "ovh_cloud_storage_block_volume_backup" "bkp" {
  service_name = "%[1]s"
  region       = "%[3]s"
  volume_id    = ovh_cloud_storage_block_volume.source.id
  name         = "%[2]s-bkp"
}

resource "ovh_cloud_storage_block_volume" "neg" {
  service_name = "%[1]s"
  region       = "%[3]s"
  name         = "%[2]s-twosrc"
  size         = 10
  create_from = {
    backup_id   = ovh_cloud_storage_block_volume_backup.bkp.id
    snapshot_id = ovh_cloud_storage_block_volume_snapshot.snap.id
  }
}
`, serviceName, volumeName, region), "Error calling Post")
}

// TestAccCloudStorageBlockVolume_invalidEncryptedBackupRestoredPlain asserts a
// backup of an encrypted volume cannot be restored into a non-encrypted volume.
func TestAccCloudStorageBlockVolume_invalidEncryptedBackupRestoredPlain(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	region := os.Getenv("OVH_CLOUD_PROJECT_REGION_TEST")
	volumeName := acctest.RandomWithPrefix(testAccResourceCloudStorageBlockVolumeNamePrefix)

	testAccCloudStorageBlockVolumeExpectError(t, fmt.Sprintf(`
resource "ovh_cloud_storage_block_volume" "enc_source" {
  service_name = "%[1]s"
  region       = "%[3]s"
  name         = "%[2]s-enc-src"
  size         = 10
  volume_type  = "CLASSIC"
  encryption   = { enabled = true }
}

resource "ovh_cloud_storage_block_volume_backup" "enc_bkp" {
  service_name = "%[1]s"
  region       = "%[3]s"
  volume_id    = ovh_cloud_storage_block_volume.enc_source.id
  name         = "%[2]s-enc-bkp"
}

resource "ovh_cloud_storage_block_volume" "neg" {
  service_name = "%[1]s"
  region       = "%[3]s"
  name         = "%[2]s-enc-plain"
  size         = 10
  create_from  = { backup_id = ovh_cloud_storage_block_volume_backup.enc_bkp.id }
  encryption   = { enabled = false }
}
`, serviceName, volumeName, region), "Error calling Post")
}

// TestAccCloudStorageBlockVolume_invalidNonexistentBackup asserts create_from
// with an unknown backup id is rejected.
func TestAccCloudStorageBlockVolume_invalidNonexistentBackup(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	region := os.Getenv("OVH_CLOUD_PROJECT_REGION_TEST")
	volumeName := acctest.RandomWithPrefix(testAccResourceCloudStorageBlockVolumeNamePrefix)

	testAccCloudStorageBlockVolumeExpectError(t, fmt.Sprintf(`
resource "ovh_cloud_storage_block_volume" "neg" {
  service_name = "%[1]s"
  region       = "%[3]s"
  name         = "%[2]s-nobackup"
  size         = 10
  create_from  = { backup_id = "%[4]s" }
}
`, serviceName, volumeName, region, testAccBlockVolumeBogusID), "Error calling Post")
}

// TestAccCloudStorageBlockVolume_invalidNonexistentSnapshot asserts create_from
// with an unknown snapshot id is rejected.
func TestAccCloudStorageBlockVolume_invalidNonexistentSnapshot(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	region := os.Getenv("OVH_CLOUD_PROJECT_REGION_TEST")
	volumeName := acctest.RandomWithPrefix(testAccResourceCloudStorageBlockVolumeNamePrefix)

	testAccCloudStorageBlockVolumeExpectError(t, fmt.Sprintf(`
resource "ovh_cloud_storage_block_volume" "neg" {
  service_name = "%[1]s"
  region       = "%[3]s"
  name         = "%[2]s-nosnap"
  size         = 10
  create_from  = { snapshot_id = "%[4]s" }
}
`, serviceName, volumeName, region, testAccBlockVolumeBogusID), "Error calling Post")
}

// TestAccCloudStorageBlockVolume_invalidShrinkOnRestore asserts restoring a
// snapshot into a volume smaller than the source is rejected.
func TestAccCloudStorageBlockVolume_invalidShrinkOnRestore(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	region := os.Getenv("OVH_CLOUD_PROJECT_REGION_TEST")
	volumeName := acctest.RandomWithPrefix(testAccResourceCloudStorageBlockVolumeNamePrefix)

	testAccCloudStorageBlockVolumeExpectError(t, fmt.Sprintf(`
resource "ovh_cloud_storage_block_volume" "source" {
  service_name = "%[1]s"
  region       = "%[3]s"
  name         = "%[2]s-shrink-src"
  size         = 20
  volume_type  = "CLASSIC"
}

resource "ovh_cloud_storage_block_volume_snapshot" "snap" {
  service_name = "%[1]s"
  region       = "%[3]s"
  volume_id    = ovh_cloud_storage_block_volume.source.id
  name         = "%[2]s-shrink-snap"
}

resource "ovh_cloud_storage_block_volume" "neg" {
  service_name = "%[1]s"
  region       = "%[3]s"
  name         = "%[2]s-shrink"
  size         = 10
  create_from  = { snapshot_id = ovh_cloud_storage_block_volume_snapshot.snap.id }
}
`, serviceName, volumeName, region), "Error calling Post")
}

// TestAccCloudStorageBlockVolume_invalidSnapshotOfNonexistentVolume asserts a
// snapshot of an unknown volume is rejected.
func TestAccCloudStorageBlockVolume_invalidSnapshotOfNonexistentVolume(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	region := os.Getenv("OVH_CLOUD_PROJECT_REGION_TEST")
	volumeName := acctest.RandomWithPrefix(testAccResourceCloudStorageBlockVolumeNamePrefix)

	testAccCloudStorageBlockVolumeExpectError(t, fmt.Sprintf(`
resource "ovh_cloud_storage_block_volume_snapshot" "neg" {
  service_name = "%[1]s"
  region       = "%[3]s"
  volume_id    = "%[4]s"
  name         = "%[2]s-snap-novol"
}
`, serviceName, volumeName, region, testAccBlockVolumeBogusID), "Error calling Post")
}

// TestAccCloudStorageBlockVolume_invalidBackupOfNonexistentVolume asserts a
// backup of an unknown volume is rejected.
func TestAccCloudStorageBlockVolume_invalidBackupOfNonexistentVolume(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	region := os.Getenv("OVH_CLOUD_PROJECT_REGION_TEST")
	volumeName := acctest.RandomWithPrefix(testAccResourceCloudStorageBlockVolumeNamePrefix)

	testAccCloudStorageBlockVolumeExpectError(t, fmt.Sprintf(`
resource "ovh_cloud_storage_block_volume_backup" "neg" {
  service_name = "%[1]s"
  region       = "%[3]s"
  volume_id    = "%[4]s"
  name         = "%[2]s-bkp-novol"
}
`, serviceName, volumeName, region, testAccBlockVolumeBogusID), "Error calling Post")
}

// TestAccCloudStorageBlockVolume_invalidVolumeType asserts a volume_type outside
// the schema's allowed set is rejected at plan by the OneOf validator.
func TestAccCloudStorageBlockVolume_invalidVolumeType(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	region := os.Getenv("OVH_CLOUD_PROJECT_REGION_TEST")
	volumeName := acctest.RandomWithPrefix(testAccResourceCloudStorageBlockVolumeNamePrefix)

	testAccCloudStorageBlockVolumeExpectError(t, fmt.Sprintf(`
resource "ovh_cloud_storage_block_volume" "neg" {
  service_name = "%[1]s"
  region       = "%[3]s"
  name         = "%[2]s-luks-type"
  size         = 10
  volume_type  = "CLASSIC_LUKS"
}
`, serviceName, volumeName, region), "value must be one of")
}
