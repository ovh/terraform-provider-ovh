package ovh

import (
	"encoding/base64"
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func testAccCloudInstanceConfig(serviceName, region, flavorID, imageID, name string) string {
	return fmt.Sprintf(`
resource "ovh_cloud_instance" "test" {
  service_name = "%s"
  region       = "%s"
  name         = "%s"
  flavor_id    = "%s"
  image_id     = "%s"

  networks = [
    { auto_assign_public_ip = true },
  ]
}
`, serviceName, region, name, flavorID, imageID)
}

func testAccCloudInstanceImportStateIdFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("resource not found: %s", resourceName)
		}
		return rs.Primary.Attributes["service_name"] + "/" + rs.Primary.Attributes["id"], nil
	}
}

func TestAccCloudInstance_basic(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	region := os.Getenv("OVH_CLOUD_PROJECT_REGION_TEST")
	flavorID := resolveInstanceFlavorID(t, serviceName, region, testAccInstanceFlavorName)
	imageID := resolveInstanceImageID(t, serviceName, region, testAccInstanceImageName)
	name := acctest.RandomWithPrefix("test-inst")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckCloudInstanceV2(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudInstanceConfig(serviceName, region, flavorID, imageID, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_instance.test", "name", name),
					resource.TestCheckResourceAttr("ovh_cloud_instance.test", "region", region),
					resource.TestCheckResourceAttr("ovh_cloud_instance.test", "flavor_id", flavorID),
					resource.TestCheckResourceAttr("ovh_cloud_instance.test", "resource_status", "READY"),
					resource.TestCheckResourceAttrSet("ovh_cloud_instance.test", "id"),
					resource.TestCheckResourceAttrSet("ovh_cloud_instance.test", "checksum"),
					resource.TestCheckResourceAttrSet("ovh_cloud_instance.test", "current_state.power_state"),
					resource.TestCheckResourceAttr("ovh_cloud_instance.test", "billing.period", "HOURLY"),
					resource.TestCheckResourceAttr("ovh_cloud_instance.test", "current_state.billing.period", "HOURLY"),
					resource.TestCheckNoResourceAttr("ovh_cloud_instance.test", "current_state.billing.since"),
				),
			},
			{
				ResourceName:      "ovh_cloud_instance.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccCloudInstanceImportStateIdFunc("ovh_cloud_instance.test"),
			},
		},
	})
}

func TestAccCloudInstance_update(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	region := os.Getenv("OVH_CLOUD_PROJECT_REGION_TEST")
	flavorID := resolveInstanceFlavorID(t, serviceName, region, testAccInstanceFlavorName)
	flavorID2 := os.Getenv("OVH_INSTANCE_FLAVOR_ID_2_TEST") // optional resize target
	imageID := resolveInstanceImageID(t, serviceName, region, testAccInstanceImageName)
	name := acctest.RandomWithPrefix("test-inst")
	nameUpdated := name + "-upd"

	if flavorID2 == "" {
		flavorID2 = flavorID // fall back to a rename-only update
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckCloudInstanceV2(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudInstanceConfig(serviceName, region, flavorID, imageID, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_instance.test", "name", name),
				),
			},
			{
				Config: testAccCloudInstanceConfig(serviceName, region, flavorID2, imageID, nameUpdated),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_instance.test", "name", nameUpdated),
					resource.TestCheckResourceAttr("ovh_cloud_instance.test", "flavor_id", flavorID2),
					resource.TestCheckResourceAttr("ovh_cloud_instance.test", "resource_status", "READY"),
				),
			},
		},
	})
}

// ---------------------------------------------------------------------------
// Additional config builders (kept separate so existing callers of
// testAccCloudInstanceConfig keep the same 5-argument signature).
// ---------------------------------------------------------------------------

// testAccCloudInstanceConfigPowerState renders an instance config with an
// explicit power_state.
func testAccCloudInstanceConfigPowerState(serviceName, region, flavorID, imageID, name, powerState string) string {
	return fmt.Sprintf(`
resource "ovh_cloud_instance" "test" {
  service_name = "%s"
  region       = "%s"
  name         = "%s"
  flavor_id    = "%s"
  image_id     = "%s"
  power_state  = "%s"

  networks = [
    { auto_assign_public_ip = true },
  ]
}
`, serviceName, region, name, flavorID, imageID, powerState)
}

// testAccCloudInstanceConfigBilling renders an instance config with an explicit
// billing period (billing omitted entirely when period == "").
func testAccCloudInstanceConfigBilling(serviceName, region, flavorID, imageID, name, period string) string {
	billing := ""
	if period != "" {
		billing = fmt.Sprintf("\n  billing = {\n    period = %q\n  }\n", period)
	}
	return fmt.Sprintf(`
resource "ovh_cloud_instance" "test" {
  service_name = "%s"
  region       = "%s"
  name         = "%s"
  flavor_id    = "%s"
  image_id     = "%s"
%s
  networks = [
    { auto_assign_public_ip = true },
  ]
}
`, serviceName, region, name, flavorID, imageID, billing)
}

// testAccCloudInstanceConfigAZ renders an instance config, optionally pinning
// the availability_zone (omitted entirely when az == "").
func testAccCloudInstanceConfigAZ(serviceName, region, flavorID, imageID, name, az string) string {
	azLine := ""
	if az != "" {
		azLine = fmt.Sprintf("\n  availability_zone = %q\n", az)
	}
	return fmt.Sprintf(`
resource "ovh_cloud_instance" "test" {
  service_name = "%s"
  region       = "%s"
  name         = "%s"
  flavor_id    = "%s"
  image_id     = "%s"%s
  networks = [
    { auto_assign_public_ip = true },
  ]
}
`, serviceName, region, name, flavorID, imageID, azLine)
}

// testAccCloudInstanceConfigSSHKey renders an ovh_cloud_ssh_key plus an instance
// whose (immutable) ssh_key_name references that key by name.
func testAccCloudInstanceConfigSSHKey(serviceName, region, flavorID, imageID, name, keyName string) string {
	return fmt.Sprintf(`
resource "ovh_cloud_ssh_key" "key" {
  service_name = "%s"
  name         = "%s"
  public_key   = "%s"
}

resource "ovh_cloud_instance" "test" {
  service_name = "%s"
  region       = "%s"
  name         = "%s"
  flavor_id    = "%s"
  image_id     = "%s"
  ssh_key_name = ovh_cloud_ssh_key.key.name

  networks = [
    { auto_assign_public_ip = true },
  ]
}
`, serviceName, keyName, testAccCloudSshKeyPublicKeyA, serviceName, region, name, flavorID, imageID)
}

// testAccCloudInstanceConfigUserData renders an instance config whose user_data
// carries a cloud-init payload, base64-encoded the documented way.
func testAccCloudInstanceConfigUserData(serviceName, region, flavorID, imageID, name, cloudInit string) string {
	return fmt.Sprintf(`
resource "ovh_cloud_instance" "test" {
  service_name = "%s"
  region       = "%s"
  name         = "%s"
  flavor_id    = "%s"
  image_id     = "%s"
  user_data    = base64encode(%q)

  networks = [
    { auto_assign_public_ip = true },
  ]
}
`, serviceName, region, name, flavorID, imageID, cloudInit)
}

// testAccCloudInstanceConfigBootFromVolume renders a bootable block volume
// (created from a Glance image) plus an instance that boots from it, with no
// image_id. volumesBlock is the HCL of the attribute under test: volume_ids
// (deprecated) or volumes.
func testAccCloudInstanceConfigBootFromVolume(serviceName, region, flavorID, glanceImageID, volName, name, volumesBlock string) string {
	return fmt.Sprintf(`
resource "ovh_cloud_storage_block_volume" "boot" {
  service_name = "%s"
  name         = "%s"
  size         = 20
  region       = "%s"
  volume_type  = "CLASSIC"

  create_from = {
    image_id = "%s"
  }
}

resource "ovh_cloud_instance" "test" {
  service_name = "%s"
  region       = "%s"
  name         = "%s"
  flavor_id    = "%s"
  %s

  networks = [
    { auto_assign_public_ip = true },
  ]
}
`, serviceName, volName, region, glanceImageID, serviceName, region, name, flavorID, volumesBlock)
}

// testAccCloudInstanceConfigBootAndDataVolumes renders a bootable block volume
// (created from a Glance image), a data volume, and an instance with no image_id
// that declares them with volumesBlock (volume_ids or volumes).
func testAccCloudInstanceConfigBootAndDataVolumes(serviceName, region, flavorID, glanceImageID, bootName, dataName, name, volumesBlock string) string {
	return fmt.Sprintf(`
resource "ovh_cloud_storage_block_volume" "boot" {
  service_name = "%[1]s"
  name         = "%[2]s"
  size         = 20
  region       = "%[3]s"
  volume_type  = "CLASSIC"

  create_from = {
    image_id = "%[4]s"
  }
}

resource "ovh_cloud_storage_block_volume" "data" {
  service_name = "%[1]s"
  name         = "%[5]s"
  size         = 10
  region       = "%[3]s"
  volume_type  = "CLASSIC"
}

resource "ovh_cloud_instance" "test" {
  service_name = "%[1]s"
  region       = "%[3]s"
  name         = "%[6]s"
  flavor_id    = "%[7]s"
  %[8]s

  networks = [
    { auto_assign_public_ip = true },
  ]
}
`, serviceName, bootName, region, glanceImageID, dataName, name, flavorID, volumesBlock)
}

// testAccCloudInstanceConfigInvalidPowerState renders a config with an invalid
// power_state to exercise the OneOf validator.
func testAccCloudInstanceConfigInvalidPowerState(serviceName, region, flavorID, imageID, name string) string {
	return fmt.Sprintf(`
resource "ovh_cloud_instance" "test" {
  service_name = "%s"
  region       = "%s"
  name         = "%s"
  flavor_id    = "%s"
  image_id     = "%s"
  power_state  = "INVALID"

  networks = [
    { auto_assign_public_ip = true },
  ]
}
`, serviceName, region, name, flavorID, imageID)
}

// testAccCloudInstanceConfigInvalidAccessLevel renders a config with an invalid
// shares.access_level to exercise the nested OneOf validator.
func testAccCloudInstanceConfigInvalidAccessLevel(serviceName, region, flavorID, imageID, name string) string {
	return fmt.Sprintf(`
resource "ovh_cloud_instance" "test" {
  service_name = "%s"
  region       = "%s"
  name         = "%s"
  flavor_id    = "%s"
  image_id     = "%s"

  networks = [
    { auto_assign_public_ip = true },
  ]

  shares = [
    {
      id           = "dummy-share-id"
      access_level = "INVALID"
    },
  ]
}
`, serviceName, region, name, flavorID, imageID)
}

// TestAccCloudInstance_renameAndResize renames and resizes the instance in a
// single update and asserts the change happens in place (id is stable), then
// re-applies the same config and asserts an empty plan.
func TestAccCloudInstance_renameAndResize(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	region := os.Getenv("OVH_CLOUD_PROJECT_REGION_TEST")
	flavorID := resolveInstanceFlavorID(t, serviceName, region, testAccInstanceFlavorName)
	flavorID2 := os.Getenv("OVH_INSTANCE_FLAVOR_ID_2_TEST") // optional resize target
	imageID := resolveInstanceImageID(t, serviceName, region, testAccInstanceImageName)
	name := acctest.RandomWithPrefix("test-inst-resize")
	nameUpdated := name + "-upd"

	if flavorID2 == "" {
		flavorID2 = flavorID // fall back to a rename-only update
	}

	var instanceID string

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckCloudInstanceV2(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudInstanceConfig(serviceName, region, flavorID, imageID, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_instance.test", "name", name),
					resource.TestCheckResourceAttr("ovh_cloud_instance.test", "flavor_id", flavorID),
					resource.TestCheckResourceAttrWith("ovh_cloud_instance.test", "id", func(v string) error {
						if v == "" {
							return fmt.Errorf("expected instance id to be set")
						}
						instanceID = v
						return nil
					}),
				),
			},
			{
				Config: testAccCloudInstanceConfig(serviceName, region, flavorID2, imageID, nameUpdated),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_instance.test", "name", nameUpdated),
					resource.TestCheckResourceAttr("ovh_cloud_instance.test", "flavor_id", flavorID2),
					resource.TestCheckResourceAttr("ovh_cloud_instance.test", "current_state.flavor.id", flavorID2),
					resource.TestCheckResourceAttr("ovh_cloud_instance.test", "resource_status", "READY"),
					resource.TestCheckResourceAttrWith("ovh_cloud_instance.test", "id", func(v string) error {
						if v != instanceID {
							return fmt.Errorf("instance was replaced during rename/resize: id changed from %q to %q", instanceID, v)
						}
						return nil
					}),
				),
			},
			{
				// Re-applying the identical config must be a no-op.
				Config: testAccCloudInstanceConfig(serviceName, region, flavorID2, imageID, nameUpdated),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						ExpectEmptyPlan(),
					},
				},
			},
		},
	})
}

// TestAccCloudInstance_rebuildImage changes image_id and asserts the instance is
// rebuilt in place (id stable) with the observed image reflecting the new image.
func TestAccCloudInstance_rebuildImage(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	region := os.Getenv("OVH_CLOUD_PROJECT_REGION_TEST")
	flavorID := resolveInstanceFlavorID(t, serviceName, region, testAccInstanceFlavorName)
	imageID := resolveInstanceImageID(t, serviceName, region, testAccInstanceImageName)
	imageID2 := os.Getenv("OVH_INSTANCE_IMAGE_ID_2_TEST") // optional rebuild target
	name := acctest.RandomWithPrefix("test-inst-rebuild")

	if imageID2 == "" {
		imageID2 = imageID // fall back to a same-image rebuild
	}

	var instanceID string

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckCloudInstanceV2(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudInstanceConfig(serviceName, region, flavorID, imageID, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_instance.test", "image_id", imageID),
					resource.TestCheckResourceAttrWith("ovh_cloud_instance.test", "id", func(v string) error {
						if v == "" {
							return fmt.Errorf("expected instance id to be set")
						}
						instanceID = v
						return nil
					}),
				),
			},
			{
				Config: testAccCloudInstanceConfig(serviceName, region, flavorID, imageID2, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_instance.test", "image_id", imageID2),
					resource.TestCheckResourceAttr("ovh_cloud_instance.test", "current_state.image.id", imageID2),
					resource.TestCheckResourceAttr("ovh_cloud_instance.test", "resource_status", "READY"),
					resource.TestCheckResourceAttrWith("ovh_cloud_instance.test", "id", func(v string) error {
						if v != instanceID {
							return fmt.Errorf("instance was replaced during image rebuild: id changed from %q to %q", instanceID, v)
						}
						return nil
					}),
				),
			},
		},
	})
}

// TestAccCloudInstance_userData creates an instance carrying cloud-init user data
// and then changes it. The API applies a new userData by rebuilding the instance
// on its current image, so the change must plan as an in-place Update and keep the
// instance id — never a destroy/recreate. The middle step also locks the absence
// of spurious drift: user_data is write-only, so a value the API never echoes back
// must still produce an empty plan.
func TestAccCloudInstance_userData(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	region := os.Getenv("OVH_CLOUD_PROJECT_REGION_TEST")
	flavorID := resolveInstanceFlavorID(t, serviceName, region, testAccInstanceFlavorName)
	imageID := resolveInstanceImageID(t, serviceName, region, testAccInstanceImageName)
	name := acctest.RandomWithPrefix("test-inst-ud")

	const (
		cloudInitA = "#cloud-config\npackages:\n  - nginx\n"
		cloudInitB = "#cloud-config\npackages:\n  - htop\n"
	)
	encodedA := base64.StdEncoding.EncodeToString([]byte(cloudInitA))
	encodedB := base64.StdEncoding.EncodeToString([]byte(cloudInitB))

	var instanceID string

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckCloudInstanceV2(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudInstanceConfigUserData(serviceName, region, flavorID, imageID, name, cloudInitA),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_instance.test", "user_data", encodedA),
					resource.TestCheckResourceAttr("ovh_cloud_instance.test", "resource_status", "READY"),
					resource.TestCheckResourceAttrWith("ovh_cloud_instance.test", "id", func(v string) error {
						if v == "" {
							return fmt.Errorf("expected instance id to be set")
						}
						instanceID = v
						return nil
					}),
				),
			},
			{
				Config:   testAccCloudInstanceConfigUserData(serviceName, region, flavorID, imageID, name, cloudInitA),
				PlanOnly: true,
			},
			{
				Config: testAccCloudInstanceConfigUserData(serviceName, region, flavorID, imageID, name, cloudInitB),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("ovh_cloud_instance.test", plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_instance.test", "user_data", encodedB),
					resource.TestCheckResourceAttr("ovh_cloud_instance.test", "resource_status", "READY"),
					resource.TestCheckResourceAttrWith("ovh_cloud_instance.test", "id", func(v string) error {
						if v != instanceID {
							return fmt.Errorf("instance was replaced by a user_data change: id went from %q to %q", instanceID, v)
						}
						return nil
					}),
				),
			},
		},
	})
}

// TestAccCloudInstance_powerState drives the instance through ACTIVE -> SHUTOFF
// -> SHELVED -> ACTIVE -> PAUSED -> ACTIVE -> SUSPENDED -> ACTIVE and asserts
// each transition is applied in place. Nova only pauses or suspends an ACTIVE
// server, hence the return to ACTIVE in between.
func TestAccCloudInstance_powerState(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	region := os.Getenv("OVH_CLOUD_PROJECT_REGION_TEST")
	flavorID := resolveInstanceFlavorID(t, serviceName, region, testAccInstanceFlavorName)
	imageID := resolveInstanceImageID(t, serviceName, region, testAccInstanceImageName)
	name := acctest.RandomWithPrefix("test-inst-power")

	var instanceID string

	powerStep := func(state string) resource.TestStep {
		return resource.TestStep{
			Config: testAccCloudInstanceConfigPowerState(serviceName, region, flavorID, imageID, name, state),
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttr("ovh_cloud_instance.test", "power_state", state),
				resource.TestCheckResourceAttrSet("ovh_cloud_instance.test", "current_state.power_state"),
				resource.TestCheckResourceAttr("ovh_cloud_instance.test", "resource_status", "READY"),
				resource.TestCheckResourceAttrWith("ovh_cloud_instance.test", "id", func(v string) error {
					if instanceID == "" {
						instanceID = v
						return nil
					}
					if v != instanceID {
						return fmt.Errorf("instance was replaced on power_state change: id changed from %q to %q", instanceID, v)
					}
					return nil
				}),
			),
		}
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckCloudInstanceV2(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			powerStep("ACTIVE"),
			powerStep("SHUTOFF"),
			powerStep("SHELVED"),
			powerStep("ACTIVE"),
			powerStep("PAUSED"),
			powerStep("ACTIVE"),
			powerStep("SUSPENDED"),
			powerStep("ACTIVE"),
		},
	})
}

// testAccCloudInstanceBootFromVolume creates a bootable block volume and boots an
// instance from it (no image_id), declaring the volume with volumesBlock. It
// asserts image_id is null, no observed image is reported, and the role BOOT is
// read back. The bootable source image is resolved from the reference API.
//
// With withImport, a last step imports the instance and verifies the imported
// state against the applied one: an import fills both volumes and the deprecated
// volume_ids, whichever one the configuration declares.
func testAccCloudInstanceBootFromVolume(t *testing.T, volumesBlock string, withImport bool, extra ...resource.TestCheckFunc) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	region := os.Getenv("OVH_CLOUD_PROJECT_REGION_TEST")
	flavorID := resolveInstanceFlavorID(t, serviceName, region, testAccInstanceFlavorName)
	glanceImageID := resolveInstanceImageID(t, serviceName, region, testAccInstanceImageName)
	name := acctest.RandomWithPrefix("test-inst-bfv")
	volName := acctest.RandomWithPrefix("test-vol-bfv")

	checks := append([]resource.TestCheckFunc{
		resource.TestCheckResourceAttr("ovh_cloud_instance.test", "name", name),
		resource.TestCheckResourceAttr("ovh_cloud_instance.test", "resource_status", "READY"),
		// No image_id was configured: MergeWith stores it as null.
		resource.TestCheckNoResourceAttr("ovh_cloud_instance.test", "image_id"),
		// The boot volume shows up in the observed volumes, with its role.
		resource.TestCheckResourceAttrSet("ovh_cloud_instance.test", "current_state.volumes.0.id"),
		resource.TestCheckResourceAttr("ovh_cloud_instance.test", "current_state.volumes.0.role", "BOOT"),
		resource.TestCheckNoResourceAttr("ovh_cloud_instance.test", "current_state.image.id"),
	}, extra...)

	steps := []resource.TestStep{
		{
			Config: testAccCloudInstanceConfigBootFromVolume(serviceName, region, flavorID, glanceImageID, volName, name, volumesBlock),
			Check:  resource.ComposeTestCheckFunc(checks...),
			// What the API echoes back must match the configuration, role included.
			ConfigPlanChecks: resource.ConfigPlanChecks{
				PostApplyPreRefresh:  []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()},
				PostApplyPostRefresh: []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()},
			},
		},
	}
	if withImport {
		steps = append(steps, resource.TestStep{
			ResourceName:      "ovh_cloud_instance.test",
			ImportState:       true,
			ImportStateVerify: true,
			ImportStateIdFunc: testAccCloudInstanceImportStateIdFunc("ovh_cloud_instance.test"),
		})
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCloudInstanceE2E(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps:                    steps,
	})
}

// The deprecated volume_ids keeps booting from its first volume.
func TestAccCloudInstance_bootFromVolume(t *testing.T) {
	testAccCloudInstanceBootFromVolume(t, `volume_ids = [ovh_cloud_storage_block_volume.boot.id]`, true,
		resource.TestCheckResourceAttr("ovh_cloud_instance.test", "volume_ids.#", "1"),
		resource.TestCheckResourceAttr("ovh_cloud_instance.test", "volumes.#", "1"),
	)
}

func TestAccCloudInstance_bootFromVolumes(t *testing.T) {
	testAccCloudInstanceBootFromVolume(t, `volumes = [{ id = ovh_cloud_storage_block_volume.boot.id, role = "BOOT" }]`, true,
		resource.TestCheckResourceAttr("ovh_cloud_instance.test", "volumes.0.role", "BOOT"),
		resource.TestCheckResourceAttr("ovh_cloud_instance.test", "volume_ids.#", "1"),
	)
}

// TestAccCloudInstance_volumeIdsToVolumes moves a boot-from-volume instance from
// the deprecated volume_ids to volumes, with the role BOOT on the same volume. The
// boot volume does not change, so this must not replace the instance, which is the
// destructive path ModifyPlan exists to avoid. The roles
// the API reports must match the configuration, DATA included, so that no plan is
// left once applied. A last step imports the instance, which fills volumes.
func TestAccCloudInstance_volumeIdsToVolumes(t *testing.T) {
	const rn = "ovh_cloud_instance.test"

	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	region := os.Getenv("OVH_CLOUD_PROJECT_REGION_TEST")
	flavorID := resolveInstanceFlavorID(t, serviceName, region, testAccInstanceFlavorName)
	glanceImageID := resolveInstanceImageID(t, serviceName, region, testAccInstanceImageName)
	name := acctest.RandomWithPrefix("test-inst-vol-migr")
	bootName := acctest.RandomWithPrefix("test-vol-boot")
	dataName := acctest.RandomWithPrefix("test-vol-data")

	cfg := func(block string) string {
		return testAccCloudInstanceConfigBootAndDataVolumes(serviceName, region, flavorID, glanceImageID, bootName, dataName, name, block)
	}
	observedRoles := resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttr(rn, "resource_status", "READY"),
		resource.TestCheckNoResourceAttr(rn, "image_id"),
		resource.TestCheckResourceAttr(rn, "current_state.volumes.#", "2"),
		resource.TestCheckTypeSetElemNestedAttrs(rn, "current_state.volumes.*", map[string]string{"name": bootName, "role": "BOOT"}),
		resource.TestCheckTypeSetElemNestedAttrs(rn, "current_state.volumes.*", map[string]string{"name": dataName, "role": "DATA"}),
	)

	var instanceID string

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCloudInstanceE2E(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				// With no image_id, the first volume of volume_ids is the boot volume.
				Config: cfg(`volume_ids = [ovh_cloud_storage_block_volume.boot.id, ovh_cloud_storage_block_volume.data.id]`),
				Check: resource.ComposeAggregateTestCheckFunc(
					observedRoles,
					resource.TestCheckResourceAttr(rn, "volume_ids.#", "2"),
					resource.TestCheckResourceAttr(rn, "volumes.#", "2"),
					captureInstanceID(rn, &instanceID),
				),
			},
			{
				Config: cfg(`volumes = [
    { id = ovh_cloud_storage_block_volume.boot.id, role = "BOOT" },
    { id = ovh_cloud_storage_block_volume.data.id },
  ]`),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					// Never a replacement. volumes already holds the same volumes in
					// the state, computed from the API while volume_ids was used, so
					// the migration changes nothing: a no-op.
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(rn, plancheck.ResourceActionNoop),
					},
					// role is Optional+Computed with a default: once applied, the roles
					// the API reports must leave nothing to change.
					PostApplyPreRefresh:  []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()},
					PostApplyPostRefresh: []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					observedRoles,
					resource.TestCheckResourceAttr(rn, "volumes.#", "2"),
					resource.TestCheckResourceAttr(rn, "volumes.0.role", "BOOT"),
					resource.TestCheckResourceAttr(rn, "volumes.1.role", "DATA"),
					resource.TestCheckResourceAttr(rn, "volume_ids.#", "2"),
					// Fails if the instance was replaced: its id would change.
					captureInstanceID(rn, &instanceID),
				),
			},
			{
				// The role left out on the boot volume stays BOOT: the API keeps the
				// stored BOOT role, so this is the same instance, not a replacement
				// that could never be created (no image_id, no BOOT volume).
				Config: cfg(`volumes = [
    { id = ovh_cloud_storage_block_volume.boot.id },
    { id = ovh_cloud_storage_block_volume.data.id },
  ]`),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(rn, plancheck.ResourceActionNoop),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					observedRoles,
					resource.TestCheckResourceAttr(rn, "volumes.0.role", "BOOT"),
					resource.TestCheckResourceAttr(rn, "volumes.1.role", "DATA"),
					// Fails if the instance was replaced: its id would change.
					captureInstanceID(rn, &instanceID),
				),
			},
			{
				// Turning the boot volume into DATA leaves no boot source: refused at
				// plan, instead of destroying an instance that cannot be recreated.
				Config: cfg(`volumes = [
    { id = ovh_cloud_storage_block_volume.boot.id, role = "DATA" },
    { id = ovh_cloud_storage_block_volume.data.id },
  ]`),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`Boot volume cannot become DATA`),
			},
			{
				// Import reads volumes back, with the BOOT role on the boot volume.
				ResourceName:      rn,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccCloudInstanceImportStateIdFunc(rn),
			},
		},
	})
}

// The refused combination fails at plan, before any call to the API.
func TestAccCloudInstance_volumesConflicts(t *testing.T) {
	cfg := func(block string) string {
		return fmt.Sprintf(`
resource "ovh_cloud_instance" "test" {
  service_name = "xxx"
  region       = "GRA11"
  name         = "conflict"
  flavor_id    = "flavor"
  %s

  networks = [{ auto_assign_public_ip = true }]
}
`, block)
	}
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      cfg("volumes = [{ id = \"a\" }]\n  volume_ids = [\"b\"]"),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`volumes cannot be combined with volume_ids`),
			},
			{
				Config:      cfg("image_id = \"img\"\n  volumes = [{ id = \"a\", role = \"BOOT\" }]"),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`mutually exclusive`),
			},
		},
	})
}

// TestAccCloudInstance_availabilityZone verifies that an omitted availability
// zone does not cause a perpetual diff (empty plan on re-apply, incl. mono-AZ
// regions where the AZ surfaces as null), and, when explicit AZ test values are
// provided, that changing the configured zone forces a replace
// (RequiresReplaceIfConfigured).
func TestAccCloudInstance_availabilityZone(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	region := os.Getenv("OVH_CLOUD_PROJECT_REGION_TEST")
	flavorID := resolveInstanceFlavorID(t, serviceName, region, testAccInstanceFlavorName)
	imageID := resolveInstanceImageID(t, serviceName, region, testAccInstanceImageName)
	azA := os.Getenv("OVH_INSTANCE_AZ_TEST")
	azB := os.Getenv("OVH_INSTANCE_AZ_2_TEST")
	name := acctest.RandomWithPrefix("test-inst-az")

	steps := []resource.TestStep{
		{
			// (a) omit availability_zone -> platform assigns one (may surface as
			// null in mono-AZ regions where OpenStack reports "nova").
			Config: testAccCloudInstanceConfigAZ(serviceName, region, flavorID, imageID, name, ""),
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttr("ovh_cloud_instance.test", "resource_status", "READY"),
			),
		},
		{
			// (a) re-apply with az still omitted -> no diff (UseStateForUnknown).
			Config: testAccCloudInstanceConfigAZ(serviceName, region, flavorID, imageID, name, ""),
			ConfigPlanChecks: resource.ConfigPlanChecks{
				PreApply: []plancheck.PlanCheck{
					ExpectEmptyPlan(),
				},
			},
		},
	}

	// (b) requires two real availability zones to exercise the replace path with
	// applies that actually succeed. Gated on env so the test still runs (a) only.
	if azA != "" && azB != "" {
		steps = append(steps,
			resource.TestStep{
				Config: testAccCloudInstanceConfigAZ(serviceName, region, flavorID, imageID, name, azA),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_instance.test", "availability_zone", azA),
				),
			},
			resource.TestStep{
				Config: testAccCloudInstanceConfigAZ(serviceName, region, flavorID, imageID, name, azB),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("ovh_cloud_instance.test", plancheck.ResourceActionReplace),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_instance.test", "availability_zone", azB),
				),
			},
		)
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckCloudInstanceV2(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps:                    steps,
	})
}

// TestAccCloudInstance_forceNew changes the immutable ssh_key_name and asserts
// the instance is replaced (ssh_key_name has RequiresReplace).
func TestAccCloudInstance_forceNew(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	region := os.Getenv("OVH_CLOUD_PROJECT_REGION_TEST")
	flavorID := resolveInstanceFlavorID(t, serviceName, region, testAccInstanceFlavorName)
	imageID := resolveInstanceImageID(t, serviceName, region, testAccInstanceImageName)
	name := acctest.RandomWithPrefix("test-inst-fn")
	keyA := acctest.RandomWithPrefix("test-key-a")
	keyB := acctest.RandomWithPrefix("test-key-b")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckCloudInstanceV2(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudInstanceConfigSSHKey(serviceName, region, flavorID, imageID, name, keyA),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_instance.test", "ssh_key_name", keyA),
					resource.TestCheckResourceAttr("ovh_cloud_instance.test", "resource_status", "READY"),
				),
			},
			{
				// Changing the (immutable) ssh_key_name must replace the instance.
				Config: testAccCloudInstanceConfigSSHKey(serviceName, region, flavorID, imageID, name, keyB),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("ovh_cloud_instance.test", plancheck.ResourceActionReplace),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_instance.test", "ssh_key_name", keyB),
				),
			},
		},
	})
}

// TestAccCloudInstance_planModifiers asserts that the UnknownDuringUpdate plan
// modifiers mark checksum, updated_at and current_state as unknown in the plan
// when a mutable attribute (name) changes.
func TestAccCloudInstance_planModifiers(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	region := os.Getenv("OVH_CLOUD_PROJECT_REGION_TEST")
	flavorID := resolveInstanceFlavorID(t, serviceName, region, testAccInstanceFlavorName)
	imageID := resolveInstanceImageID(t, serviceName, region, testAccInstanceImageName)
	name := acctest.RandomWithPrefix("test-inst-pm")
	nameUpdated := name + "-upd"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckCloudInstanceV2(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudInstanceConfig(serviceName, region, flavorID, imageID, name),
			},
			{
				Config: testAccCloudInstanceConfig(serviceName, region, flavorID, imageID, nameUpdated),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectUnknownValue("ovh_cloud_instance.test", tfjsonpath.New("checksum")),
						plancheck.ExpectUnknownValue("ovh_cloud_instance.test", tfjsonpath.New("updated_at")),
						plancheck.ExpectUnknownValue("ovh_cloud_instance.test", tfjsonpath.New("current_state")),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_instance.test", "name", nameUpdated),
				),
			},
		},
	})
}

// TestAccCloudInstance_validators asserts schema validators reject invalid
// power_state and shares.access_level values at plan time.
func TestAccCloudInstance_validators(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	region := os.Getenv("OVH_CLOUD_PROJECT_REGION_TEST")
	flavorID := resolveInstanceFlavorID(t, serviceName, region, testAccInstanceFlavorName)
	imageID := resolveInstanceImageID(t, serviceName, region, testAccInstanceImageName)
	name := acctest.RandomWithPrefix("test-inst-val")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckCloudInstanceV2(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccCloudInstanceConfigInvalidPowerState(serviceName, region, flavorID, imageID, name),
				ExpectError: regexp.MustCompile(`must be one of`),
			},
			{
				Config:      testAccCloudInstanceConfigInvalidAccessLevel(serviceName, region, flavorID, imageID, name),
				ExpectError: regexp.MustCompile(`must be one of`),
			},
			{
				Config:      testAccCloudInstanceConfigBilling(serviceName, region, flavorID, imageID, name, "WEEKLY"),
				ExpectError: regexp.MustCompile(`must be one of`),
			},
		},
	})
}

// TestAccCloudInstance_monthlyBilling moves an hourly instance to monthly
// billing, then checks that the API refuses going back to HOURLY. Monthly
// billing cannot be undone and is billed for the month, so the test only runs
// when OVH_CLOUD_INSTANCE_MONTHLY_BILLING_TEST is set.
func TestAccCloudInstance_monthlyBilling(t *testing.T) {
	checkEnvOrSkip(t, "OVH_CLOUD_INSTANCE_MONTHLY_BILLING_TEST")
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	region := os.Getenv("OVH_CLOUD_PROJECT_REGION_TEST")
	flavorID := resolveInstanceFlavorID(t, serviceName, region, testAccInstanceFlavorName)
	imageID := resolveInstanceImageID(t, serviceName, region, testAccInstanceImageName)
	name := acctest.RandomWithPrefix("test-inst-billing")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckCloudInstanceV2(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudInstanceConfigBilling(serviceName, region, flavorID, imageID, name, ""),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_instance.test", "billing.period", "HOURLY"),
					resource.TestCheckResourceAttr("ovh_cloud_instance.test", "current_state.billing.period", "HOURLY"),
				),
			},
			{
				Config: testAccCloudInstanceConfigBilling(serviceName, region, flavorID, imageID, name, "MONTHLY"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_instance.test", "billing.period", "MONTHLY"),
					resource.TestCheckResourceAttr("ovh_cloud_instance.test", "resource_status", "READY"),
					resource.TestCheckResourceAttr("ovh_cloud_instance.test", "current_state.billing.period", "MONTHLY"),
					resource.TestCheckResourceAttrSet("ovh_cloud_instance.test", "current_state.billing.since"),
					resource.TestCheckResourceAttrSet("ovh_cloud_instance.test", "current_state.billing.status"),
				),
			},
			{
				Config:      testAccCloudInstanceConfigBilling(serviceName, region, flavorID, imageID, name, "HOURLY"),
				ExpectError: regexp.MustCompile(`monthly billing cannot be undone`),
			},
			{
				// Leaving billing out keeps the current period: no change.
				Config:   testAccCloudInstanceConfigBilling(serviceName, region, flavorID, imageID, name, ""),
				PlanOnly: true,
			},
			{
				ResourceName:      "ovh_cloud_instance.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccCloudInstanceImportStateIdFunc("ovh_cloud_instance.test"),
			},
		},
	})
}
