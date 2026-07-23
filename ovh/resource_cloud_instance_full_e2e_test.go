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
// TIER 2 — group G: the flagship full-stack composition.
//
// TestAccCloudInstance_fullStack wires an ovh_cloud_instance to EVERY companion
// resource it can consume — a vRack network + subnet, an egress gateway, a
// floating IP, security groups, an ssh key, two block volumes, a file share
// (via its share network), and a placement group — plus the singular and plural
// data sources reading it back. It then drives an in-place update touching the
// mutable surface (rename, resize, power off, second volume, security-group
// swap) and proves the config is idempotent.
// ---------------------------------------------------------------------------

// fullStackNames bundles the randomized names of the companion resources so
// they stay identical across the create and update steps (otherwise the
// companions would be recreated and the instance replaced).
type fullStackNames struct {
	net, subnet, gw, sg1, sg2, key, vol1, vol2, sharenet, share, grp string
}

// testAccCloudInstanceFullStackConfig renders the whole graph. The companion
// resources are constant; only the instance's mutable fields (name, flavor,
// power state, attached volumes, attached security group) vary per step.
func testAccCloudInstanceFullStackConfig(serviceName, region, imageID, sshPublicKey string, n fullStackNames, instName, instFlavor, powerState, sgRef, volRefs string) string {
	powerLine := ""
	if powerState != "" {
		powerLine = fmt.Sprintf("\n  power_state  = %q", powerState)
	}

	return fmt.Sprintf(`
resource "ovh_cloud_network_private_vrack" "net" {
  service_name = "%s"
  name         = "%s"
  region       = "%s"
}

resource "ovh_cloud_network_private_vrack_subnet" "subnet" {
  service_name = ovh_cloud_network_private_vrack.net.service_name
  network_id   = ovh_cloud_network_private_vrack.net.id
  name         = "%s"
  cidr         = "10.0.0.0/24"
  gateway_ip   = "10.0.0.1"
  dhcp_enabled = true
  region       = "%s"
}

resource "ovh_cloud_gateway" "gw" {
  service_name = ovh_cloud_network_private_vrack.net.service_name
  name         = "%s"
  region       = "%s"
  subnet_ids   = [ovh_cloud_network_private_vrack_subnet.subnet.id]

  external_gateway = {
    enabled = true
    model   = "S"
  }
}

resource "ovh_cloud_floating_ip" "fip" {
  service_name = "%s"
  region       = "%s"
}

resource "ovh_cloud_security_group" "sg1" {
  service_name = "%s"
  region       = "%s"
  name         = "%s"
}

resource "ovh_cloud_security_group" "sg2" {
  service_name = "%s"
  region       = "%s"
  name         = "%s"
}

resource "ovh_cloud_ssh_key" "key" {
  service_name = "%s"
  name         = "%s"
  public_key   = "%s"
}

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

resource "ovh_cloud_storage_file_share_network" "sharenet" {
  service_name = "%s"
  name         = "%s"
  network_id   = ovh_cloud_network_private_vrack.net.id
  subnet_id    = ovh_cloud_network_private_vrack_subnet.subnet.id
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

resource "ovh_cloud_instance_group" "grp" {
  service_name = "%s"
  region       = "%s"
  name         = "%s"
  policy       = "ANTI_AFFINITY"
}

resource "ovh_cloud_instance" "test" {
  service_name       = "%s"
  region             = "%s"
  name               = "%s"
  flavor_id          = "%s"
  image_id           = "%s"%s
  group_id           = ovh_cloud_instance_group.grp.id
  ssh_key_name       = ovh_cloud_ssh_key.key.name
  security_group_ids = [%s]
  volume_ids         = [%s]

  networks = [
    {
      public     = false
      network_id = ovh_cloud_network_private_vrack.net.id
      subnet_id  = ovh_cloud_network_private_vrack_subnet.subnet.id
    },
    {
      public         = true
      floating_ip_id = ovh_cloud_floating_ip.fip.current_state.id
    },
  ]

  shares = [
    {
      id           = ovh_cloud_storage_file_share.share.id
      access_level = "READ_WRITE"
    },
  ]
}

data "ovh_cloud_instance" "by_id" {
  service_name = ovh_cloud_instance.test.service_name
  id           = ovh_cloud_instance.test.id
}

data "ovh_cloud_instances" "all" {
  service_name = ovh_cloud_instance.test.service_name
}
`,
		serviceName, n.net, region, // net
		n.subnet, region, // subnet
		n.gw, region, // gateway
		serviceName, region, // floating ip
		serviceName, region, n.sg1, // sg1
		serviceName, region, n.sg2, // sg2
		serviceName, n.key, sshPublicKey, // ssh key
		serviceName, region, n.vol1, // vol1
		serviceName, region, n.vol2, // vol2
		serviceName, n.sharenet, region, // share network
		serviceName, n.share, region, // share
		serviceName, region, n.grp, // instance group
		serviceName, region, instName, instFlavor, imageID, powerLine, sgRef, volRefs, // instance
	)
}

func TestAccCloudInstance_fullStack(t *testing.T) {
	const rn = "ovh_cloud_instance.test"

	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	region := os.Getenv("OVH_CLOUD_PROJECT_REGION_TEST")
	flavorID := resolveInstanceFlavorID(t, serviceName, region, testAccInstanceFlavorName)
	flavorID2 := os.Getenv("OVH_INSTANCE_FLAVOR_ID_2_TEST") // optional resize target
	imageID := resolveInstanceImageID(t, serviceName, region, testAccInstanceImageName)

	if flavorID2 == "" {
		flavorID2 = flavorID // fall back to a rename-only resize
	}

	n := fullStackNames{
		net:      acctest.RandomWithPrefix("tf-test-fs-net"),
		subnet:   acctest.RandomWithPrefix("tf-test-fs-subnet"),
		gw:       acctest.RandomWithPrefix("tf-test-fs-gw"),
		sg1:      acctest.RandomWithPrefix("tf-test-fs-sg1"),
		sg2:      acctest.RandomWithPrefix("tf-test-fs-sg2"),
		key:      acctest.RandomWithPrefix("tf-test-fs-key"),
		vol1:     acctest.RandomWithPrefix("tf-test-fs-vol1"),
		vol2:     acctest.RandomWithPrefix("tf-test-fs-vol2"),
		sharenet: acctest.RandomWithPrefix("tf-test-fs-sharenet"),
		share:    acctest.RandomWithPrefix("tf-test-fs-share"),
		grp:      acctest.RandomWithPrefix("tf-test-fs-grp"),
	}

	name := acctest.RandomWithPrefix("test-inst-fullstack")
	nameUpdated := name + "-upd"

	// Step 1: create the full graph; only vol1 and sg1 are attached to start.
	createConfig := testAccCloudInstanceFullStackConfig(
		serviceName, region, imageID, testAccCloudSshKeyPublicKeyA, n,
		name, flavorID, "",
		"ovh_cloud_security_group.sg1.id",
		"ovh_cloud_storage_block_volume.vol1.id",
	)

	// Step 2: in-place update — rename, resize, power off, attach vol2, swap SG.
	updateConfig := testAccCloudInstanceFullStackConfig(
		serviceName, region, imageID, testAccCloudSshKeyPublicKeyA, n,
		nameUpdated, flavorID2, "SHUTOFF",
		"ovh_cloud_security_group.sg2.id",
		"ovh_cloud_storage_block_volume.vol1.id, ovh_cloud_storage_block_volume.vol2.id",
	)

	var instanceID string

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckCloudInstanceE2E(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: createConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(rn, "resource_status", "READY"),
					resource.TestCheckResourceAttr(rn, "name", name),
					// Flavor + image observed as requested.
					resource.TestCheckResourceAttr(rn, "flavor_id", flavorID),
					resource.TestCheckResourceAttr(rn, "current_state.flavor.id", flavorID),
					resource.TestCheckResourceAttr(rn, "image_id", imageID),
					resource.TestCheckResourceAttr(rn, "current_state.image.id", imageID),
					// SSH key (immutable) injected at boot.
					resource.TestCheckResourceAttrPair(rn, "ssh_key_name", "ovh_cloud_ssh_key.key", "name"),
					resource.TestCheckResourceAttr(rn, "current_state.ssh_key_name", n.key),
					// Placement group membership.
					resource.TestCheckResourceAttrPair(rn, "group_id", "ovh_cloud_instance_group.grp", "id"),
					resource.TestCheckResourceAttrPair(rn, "current_state.group.id", "ovh_cloud_instance_group.grp", "id"),
					// Two NICs: gateway egress on the private one, floating IP on the public one.
					resource.TestCheckResourceAttr(rn, "networks.#", "2"),
					resource.TestCheckResourceAttr(rn, "current_state.networks.#", "2"),
					resource.TestCheckResourceAttrSet(rn, "current_state.networks.0.id"),
					resource.TestCheckResourceAttrSet(rn, "current_state.networks.1.id"),
					testAccCheckInstanceAnyNetworkAttrSet(rn, "gateway_id"),
					testAccCheckInstanceAnyNetworkAttrSet(rn, "floating_ip_id"),
					// One volume attached to start.
					resource.TestCheckResourceAttr(rn, "volume_ids.#", "1"),
					resource.TestCheckResourceAttr(rn, "current_state.volumes.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(rn, "current_state.volumes.*", map[string]string{
						"name": n.vol1,
						"size": "10",
					}),
					// One security group attached to start.
					resource.TestCheckResourceAttr(rn, "security_group_ids.#", "1"),
					resource.TestCheckResourceAttrPair(rn, "security_group_ids.0", "ovh_cloud_security_group.sg1", "id"),
					resource.TestCheckResourceAttr(rn, "current_state.security_groups.#", "1"),
					// Share attached and converged.
					resource.TestCheckResourceAttr(rn, "shares.#", "1"),
					resource.TestCheckResourceAttr(rn, "shares.0.access_level", "READ_WRITE"),
					resource.TestCheckResourceAttr(rn, "current_state.shares.#", "1"),
					resource.TestCheckResourceAttr(rn, "current_state.shares.0.access_level", "READ_WRITE"),
					resource.TestCheckResourceAttrSet(rn, "current_state.shares.0.state"),
					// Data-source read-backs.
					resource.TestCheckResourceAttr("data.ovh_cloud_instance.by_id", "name", name),
					resource.TestCheckResourceAttr("data.ovh_cloud_instance.by_id", "flavor_id", flavorID),
					resource.TestCheckResourceAttrPair("data.ovh_cloud_instance.by_id", "id", rn, "id"),
					resource.TestCheckResourceAttrPair("data.ovh_cloud_instance.by_id", "current_state.group.id", "ovh_cloud_instance_group.grp", "id"),
					resource.TestCheckResourceAttr("data.ovh_cloud_instance.by_id", "current_state.shares.#", "1"),
					resource.TestCheckResourceAttrWith("data.ovh_cloud_instances.all", "instances.#", testAccCheckCloudPublicIPListNotEmpty),
					captureInstanceID(rn, &instanceID),
				),
			},
			{
				// In-place update across the mutable surface; id must be stable.
				Config: updateConfig,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(rn, plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(rn, "resource_status", "READY"),
					resource.TestCheckResourceAttr(rn, "name", nameUpdated),
					resource.TestCheckResourceAttr(rn, "flavor_id", flavorID2),
					resource.TestCheckResourceAttr(rn, "current_state.flavor.id", flavorID2),
					resource.TestCheckResourceAttr(rn, "power_state", "SHUTOFF"),
					resource.TestCheckResourceAttr(rn, "current_state.power_state", "SHUTOFF"),
					// Second volume now attached.
					resource.TestCheckResourceAttr(rn, "volume_ids.#", "2"),
					resource.TestCheckResourceAttr(rn, "current_state.volumes.#", "2"),
					// Security group swapped.
					resource.TestCheckResourceAttr(rn, "security_group_ids.#", "1"),
					resource.TestCheckResourceAttrPair(rn, "security_group_ids.0", "ovh_cloud_security_group.sg2", "id"),
					resource.TestCheckResourceAttr(rn, "current_state.security_groups.#", "1"),
					captureInstanceID(rn, &instanceID),
				),
			},
			{
				// Re-applying the update config must be a no-op.
				Config: updateConfig,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
		},
	})
}
