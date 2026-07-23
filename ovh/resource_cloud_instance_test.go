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
    { public = true },
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
		PreCheck:                 func() { testAccPreCheckCloudInstance(t) },
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
				),
			},
			{
				ResourceName:      "ovh_cloud_instance.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccCloudInstanceImportStateIdFunc("ovh_cloud_instance.test"),
				// current_state / checksum are refreshed on read; ignore volatile fields if import verify complains.
				ImportStateVerifyIgnore: []string{"checksum"},
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
		PreCheck:                 func() { testAccPreCheckCloudInstance(t) },
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
    { public = true },
  ]
}
`, serviceName, region, name, flavorID, imageID, powerState)
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
    { public = true },
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
    { public = true },
  ]
}
`, serviceName, keyName, testAccCloudSshKeyPublicKeyA, serviceName, region, name, flavorID, imageID)
}

// testAccCloudInstanceConfigBootFromVolume renders a bootable block volume
// (created from a Glance image) plus an instance that boots from it: no
// image_id is set and the volume id is placed in volume_ids.
func testAccCloudInstanceConfigBootFromVolume(serviceName, region, flavorID, glanceImageID, volName, name string) string {
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
  volume_ids   = [ovh_cloud_storage_block_volume.boot.id]

  networks = [
    { public = true },
  ]
}
`, serviceName, volName, region, glanceImageID, serviceName, region, name, flavorID)
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
    { public = true },
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
    { public = true },
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
		PreCheck:                 func() { testAccPreCheckCloudInstance(t) },
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
		PreCheck:                 func() { testAccPreCheckCloudInstance(t) },
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

// TestAccCloudInstance_powerState drives the instance through ACTIVE -> SHUTOFF
// -> SHELVED -> ACTIVE and asserts each transition is applied in place.
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
		PreCheck:                 func() { testAccPreCheckCloudInstance(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			powerStep("ACTIVE"),
			powerStep("SHUTOFF"),
			powerStep("SHELVED"),
			powerStep("ACTIVE"),
		},
	})
}

// TestAccCloudInstance_bootFromVolume creates a bootable block volume and boots
// an instance from it (no image_id). It asserts image_id is null and no observed
// image is reported. The bootable source image is resolved from the reference API.
func TestAccCloudInstance_bootFromVolume(t *testing.T) {
	// Boot-from-volume is not yet implemented in the public-cloud-apiv2 backend:
	// instance create always passes imageRef and never builds a block-device
	// mapping (boot_index=0) from the volumes list, so omitting image_id yields
	// a 400 "Missing imageRef attribute". Re-enable once the backend supports it.
	t.Skip("boot-from-volume unsupported by public-cloud-apiv2 backend (create requires imageRef)")

	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	region := os.Getenv("OVH_CLOUD_PROJECT_REGION_TEST")
	flavorID := resolveInstanceFlavorID(t, serviceName, region, testAccInstanceFlavorName)
	glanceImageID := resolveInstanceImageID(t, serviceName, region, testAccInstanceImageName)
	name := acctest.RandomWithPrefix("test-inst-bfv")
	volName := acctest.RandomWithPrefix("test-vol-bfv")

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCloudInstanceE2E(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudInstanceConfigBootFromVolume(serviceName, region, flavorID, glanceImageID, volName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_instance.test", "name", name),
					resource.TestCheckResourceAttr("ovh_cloud_instance.test", "resource_status", "READY"),
					// No image_id was configured: MergeWith stores it as null.
					resource.TestCheckNoResourceAttr("ovh_cloud_instance.test", "image_id"),
					// The boot volume shows up in the observed volumes.
					resource.TestCheckResourceAttrSet("ovh_cloud_instance.test", "current_state.volumes.0.id"),
					// TODO(e2e): confirm the API leaves current_state.image null for
					// boot-from-volume; if it echoes the source image, relax this.
					resource.TestCheckNoResourceAttr("ovh_cloud_instance.test", "current_state.image.id"),
				),
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
		PreCheck:                 func() { testAccPreCheckCloudInstance(t) },
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
		PreCheck:                 func() { testAccPreCheckCloudInstance(t) },
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
		PreCheck:                 func() { testAccPreCheckCloudInstance(t) },
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
		PreCheck:                 func() { testAccPreCheckCloudInstance(t) },
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
		},
	})
}
