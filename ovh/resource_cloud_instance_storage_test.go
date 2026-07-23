package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
)

// ---------------------------------------------------------------------------
// TIER 2 — group D: ovh_cloud_instance storage compositions (block volumes and
// file shares).
//
// Share tests reuse the shared testAccVrackNetworkSubnetConfig helper (vRack
// net + subnet, cidr 10.0.0.0/24) so the instance can attach the SAME private
// network+subnet that backs the file_share_network, which makes the share
// reachable. access_rules is intentionally left UNSET on shares attached to an
// instance: the instance co-manages the rules and setting them here would drive
// the share OUT_OF_SYNC.
// ---------------------------------------------------------------------------

// TestAccCloudInstance_volumesAttachDetach attaches two CLASSIC block volumes,
// asserts both surface in the observed volumes, then detaches one in place.
func TestAccCloudInstance_volumesAttachDetach(t *testing.T) {
	const rn = "ovh_cloud_instance.test"

	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	region := os.Getenv("OVH_CLOUD_PROJECT_REGION_TEST")
	flavorID := resolveInstanceFlavorID(t, serviceName, region, testAccInstanceFlavorName)
	imageID := resolveInstanceImageID(t, serviceName, region, testAccInstanceImageName)

	volName1 := acctest.RandomWithPrefix("tf-test-inst-vol1")
	volName2 := acctest.RandomWithPrefix("tf-test-inst-vol2")
	name := acctest.RandomWithPrefix("test-inst-vol")

	volumes := fmt.Sprintf(`
resource "ovh_cloud_storage_block_volume" "vol1" {
  service_name = "%s"
  region       = "%s"
  name         = "%s"
  size         = 10
  volume_type  = "CLASSIC"
}

resource "ovh_cloud_storage_block_volume" "vol2" {
  service_name = "%s"
  region       = "%s"
  name         = "%s"
  size         = 10
  volume_type  = "CLASSIC"
}
`, serviceName, region, volName1, serviceName, region, volName2)

	instance := func(volRefs string) string {
		return volumes + fmt.Sprintf(`
resource "ovh_cloud_instance" "test" {
  service_name = "%s"
  region       = "%s"
  name         = "%s"
  flavor_id    = "%s"
  image_id     = "%s"
  volume_ids   = [%s]

  networks = [
    { public = true },
  ]
}
`, serviceName, region, name, flavorID, imageID, volRefs)
	}

	var instanceID string

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckCloudInstanceNet(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: instance("ovh_cloud_storage_block_volume.vol1.id, ovh_cloud_storage_block_volume.vol2.id"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(rn, "resource_status", "READY"),
					resource.TestCheckResourceAttr(rn, "volume_ids.#", "2"),
					resource.TestCheckResourceAttr(rn, "current_state.volumes.#", "2"),
					// Both attached volumes surface (order-independent) with size.
					resource.TestCheckTypeSetElemNestedAttrs(rn, "current_state.volumes.*", map[string]string{
						"name": volName1,
						"size": "10",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(rn, "current_state.volumes.*", map[string]string{
						"name": volName2,
						"size": "10",
					}),
					resource.TestCheckResourceAttrSet(rn, "current_state.volumes.0.id"),
					captureInstanceID(rn, &instanceID),
				),
			},
			{
				// Detach the second volume — in-place update.
				Config: instance("ovh_cloud_storage_block_volume.vol1.id"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(rn, plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(rn, "volume_ids.#", "1"),
					resource.TestCheckResourceAttr(rn, "current_state.volumes.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(rn, "current_state.volumes.*", map[string]string{
						"name": volName1,
						"size": "10",
					}),
					captureInstanceID(rn, &instanceID),
				),
			},
		},
	})
}

// testAccCloudInstanceShareConfig renders a vRack net+subnet, a file share
// network + file share on that subnet (access_rules UNSET), and an instance on
// the same subnet attaching the share. accessLevelLine is the optional
// access_level line inside the shares element (empty to omit it).
func testAccCloudInstanceShareConfig(serviceName, region, vrackNetName, vrackSubnetName, shareNetName, shareName, instName, flavorID, imageID, accessLevelLine string) string {
	return testAccVrackNetworkSubnetConfig(serviceName, region, vrackNetName, vrackSubnetName) + fmt.Sprintf(`
resource "ovh_cloud_storage_file_share_network" "sharenet" {
  service_name = "%s"
  name         = "%s"
  network_id   = ovh_cloud_network_private_vrack.vrack_net.id
  subnet_id    = ovh_cloud_network_private_vrack_subnet.vrack_subnet.id
  region       = "%s"
}

resource "ovh_cloud_storage_file_share" "share" {
  service_name     = "%s"
  name             = "%s"
  size             = 150
  region           = "%s"
  protocol         = "NFS"
  share_type       = "STANDARD_1AZ"
  share_network_id = ovh_cloud_storage_file_share_network.sharenet.id
}

resource "ovh_cloud_instance" "test" {
  service_name = "%s"
  region       = "%s"
  name         = "%s"
  flavor_id    = "%s"
  image_id     = "%s"

  networks = [
    {
      public     = false
      network_id = ovh_cloud_network_private_vrack.vrack_net.id
      subnet_id  = ovh_cloud_network_private_vrack_subnet.vrack_subnet.id
    },
  ]

  shares = [
    {
      id = ovh_cloud_storage_file_share.share.id%s
    },
  ]
}
`, serviceName, shareNetName, region, serviceName, shareName, region, serviceName, region, instName, flavorID, imageID, accessLevelLine)
}

// TestAccCloudInstance_sharesReadWrite attaches a file share to an instance with
// an explicit READ_WRITE access level and asserts both the requested and
// observed share state, including that the share converges (not OUT_OF_SYNC).
func TestAccCloudInstance_sharesReadWrite(t *testing.T) {
	const rn = "ovh_cloud_instance.test"

	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	region := os.Getenv("OVH_CLOUD_PROJECT_REGION_TEST")
	flavorID := resolveInstanceFlavorID(t, serviceName, region, testAccInstanceFlavorName)
	imageID := resolveInstanceImageID(t, serviceName, region, testAccInstanceImageName)

	vrackNetName := acctest.RandomWithPrefix("tf-test-inst-vnet")
	vrackSubnetName := acctest.RandomWithPrefix("tf-test-inst-vsub")
	shareNetName := acctest.RandomWithPrefix("tf-test-inst-sharenet")
	shareName := acctest.RandomWithPrefix("tf-test-inst-share")
	name := acctest.RandomWithPrefix("test-inst-share")

	config := testAccCloudInstanceShareConfig(
		serviceName, region, vrackNetName, vrackSubnetName, shareNetName, shareName, name, flavorID, imageID,
		"\n      access_level = \"READ_WRITE\"",
	)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckCloudInstanceNet(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(rn, "resource_status", "READY"),
					resource.TestCheckResourceAttr(rn, "shares.#", "1"),
					resource.TestCheckResourceAttrPair(rn, "shares.0.id", "ovh_cloud_storage_file_share.share", "id"),
					resource.TestCheckResourceAttr(rn, "shares.0.access_level", "READ_WRITE"),
					// Observed shares are only populated on a single-instance read
					// (which the resource performs).
					resource.TestCheckResourceAttr(rn, "current_state.shares.#", "1"),
					resource.TestCheckResourceAttrSet(rn, "current_state.shares.0.id"),
					resource.TestCheckResourceAttr(rn, "current_state.shares.0.access_level", "READ_WRITE"),
					resource.TestCheckResourceAttrSet(rn, "current_state.shares.0.state"),
				),
			},
		},
	})
}

// TestAccCloudInstance_sharesReadOnlyAndDefault covers the READ_ONLY access
// level and the omitted-access_level case (server-defaulted to READ_WRITE).
func TestAccCloudInstance_sharesReadOnlyAndDefault(t *testing.T) {
	const rn = "ovh_cloud_instance.test"

	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	region := os.Getenv("OVH_CLOUD_PROJECT_REGION_TEST")
	flavorID := resolveInstanceFlavorID(t, serviceName, region, testAccInstanceFlavorName)
	imageID := resolveInstanceImageID(t, serviceName, region, testAccInstanceImageName)

	cases := []struct {
		name            string
		accessLevelLine string
		wantAccessLevel string
	}{
		{
			name:            "read_only",
			accessLevelLine: "\n      access_level = \"READ_ONLY\"",
			wantAccessLevel: "READ_ONLY",
		},
		{
			// access_level omitted: the platform defaults it to READ_WRITE, which
			// must surface in the observed share state.
			name:            "default_read_write",
			accessLevelLine: "",
			wantAccessLevel: "READ_WRITE",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			vrackNetName := acctest.RandomWithPrefix("tf-test-inst-vnet")
			vrackSubnetName := acctest.RandomWithPrefix("tf-test-inst-vsub")
			shareNetName := acctest.RandomWithPrefix("tf-test-inst-sharenet")
			shareName := acctest.RandomWithPrefix("tf-test-inst-share")
			name := acctest.RandomWithPrefix("test-inst-share")

			config := testAccCloudInstanceShareConfig(
				serviceName, region, vrackNetName, vrackSubnetName, shareNetName, shareName, name, flavorID, imageID,
				tc.accessLevelLine,
			)

			resource.Test(t, resource.TestCase{
				PreCheck:                 func() { testAccPreCheckCloudInstanceNet(t) },
				ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
				Steps: []resource.TestStep{
					{
						Config: config,
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr(rn, "resource_status", "READY"),
							resource.TestCheckResourceAttr(rn, "shares.#", "1"),
							resource.TestCheckResourceAttr(rn, "current_state.shares.#", "1"),
							resource.TestCheckResourceAttr(rn, "current_state.shares.0.access_level", tc.wantAccessLevel),
							resource.TestCheckResourceAttrSet(rn, "current_state.shares.0.state"),
						),
					},
				},
			})
		})
	}
}

// TestAccCloudInstance_shareAttachDetach attaches a share, then detaches it in
// place, asserting the share list transitions and the instance stays READY
// (i.e. is not left OUT_OF_SYNC by the co-managed access rules).
func TestAccCloudInstance_shareAttachDetach(t *testing.T) {
	const rn = "ovh_cloud_instance.test"

	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	region := os.Getenv("OVH_CLOUD_PROJECT_REGION_TEST")
	flavorID := resolveInstanceFlavorID(t, serviceName, region, testAccInstanceFlavorName)
	imageID := resolveInstanceImageID(t, serviceName, region, testAccInstanceImageName)

	vrackNetName := acctest.RandomWithPrefix("tf-test-inst-vnet")
	vrackSubnetName := acctest.RandomWithPrefix("tf-test-inst-vsub")
	shareNetName := acctest.RandomWithPrefix("tf-test-inst-sharenet")
	shareName := acctest.RandomWithPrefix("tf-test-inst-share")
	name := acctest.RandomWithPrefix("test-inst-share")

	// Companion resources are identical across both steps; only the instance's
	// shares list changes. The private NIC stays attached in both steps.
	base := testAccVrackNetworkSubnetConfig(serviceName, region, vrackNetName, vrackSubnetName) + fmt.Sprintf(`
resource "ovh_cloud_storage_file_share_network" "sharenet" {
  service_name = "%s"
  name         = "%s"
  network_id   = ovh_cloud_network_private_vrack.vrack_net.id
  subnet_id    = ovh_cloud_network_private_vrack_subnet.vrack_subnet.id
  region       = "%s"
}

resource "ovh_cloud_storage_file_share" "share" {
  service_name     = "%s"
  name             = "%s"
  size             = 150
  region           = "%s"
  protocol         = "NFS"
  share_type       = "STANDARD_1AZ"
  share_network_id = ovh_cloud_storage_file_share_network.sharenet.id
}
`, serviceName, shareNetName, region, serviceName, shareName, region)

	withShare := base + fmt.Sprintf(`
resource "ovh_cloud_instance" "test" {
  service_name = "%s"
  region       = "%s"
  name         = "%s"
  flavor_id    = "%s"
  image_id     = "%s"

  networks = [
    {
      public     = false
      network_id = ovh_cloud_network_private_vrack.vrack_net.id
      subnet_id  = ovh_cloud_network_private_vrack_subnet.vrack_subnet.id
    },
  ]

  shares = [
    {
      id           = ovh_cloud_storage_file_share.share.id
      access_level = "READ_WRITE"
    },
  ]
}
`, serviceName, region, name, flavorID, imageID)

	withoutShare := base + fmt.Sprintf(`
resource "ovh_cloud_instance" "test" {
  service_name = "%s"
  region       = "%s"
  name         = "%s"
  flavor_id    = "%s"
  image_id     = "%s"

  networks = [
    {
      public     = false
      network_id = ovh_cloud_network_private_vrack.vrack_net.id
      subnet_id  = ovh_cloud_network_private_vrack_subnet.vrack_subnet.id
    },
  ]
}
`, serviceName, region, name, flavorID, imageID)

	var instanceID string

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckCloudInstanceNet(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: withShare,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(rn, "resource_status", "READY"),
					resource.TestCheckResourceAttr(rn, "shares.#", "1"),
					resource.TestCheckResourceAttr(rn, "current_state.shares.#", "1"),
					captureInstanceID(rn, &instanceID),
				),
			},
			{
				// Detach the share — in-place update; must stay READY.
				Config: withoutShare,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(rn, plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(rn, "shares.#", "0"),
					resource.TestCheckResourceAttr(rn, "resource_status", "READY"),
					captureInstanceID(rn, &instanceID),
				),
			},
		},
	})
}
