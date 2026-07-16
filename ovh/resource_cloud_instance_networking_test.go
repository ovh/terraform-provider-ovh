package ovh

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// ---------------------------------------------------------------------------
// TIER 2 — group C: ovh_cloud_instance networking compositions.
//
// These tests wire an instance to a private vRack network+subnet, gateway,
// floating IP and security groups by resource reference (no depends_on). They
// all reuse testAccPreCheckCloudInstanceNet (creds + SERVICE/REGION/VRACK/
// FLAVOR_ID/IMAGE_ID) and colocate every companion resource in the instance
// region so associations succeed.
// ---------------------------------------------------------------------------

// testAccCloudInstanceNetSubnetConfig renders a vRack private network and a
// DHCP-enabled /24 subnet, exposed as:
//   - ovh_cloud_network_private_vrack.net
//   - ovh_cloud_network_private_vrack_subnet.subnet
//
// DHCP is enabled so instances attaching to the subnet obtain an address,
// mirroring the canonical TestAccCloudGateway_withSubnets composition.
func testAccCloudInstanceNetSubnetConfig(serviceName, region, netName, subnetName string) string {
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
  dhcp_enabled = true
  region       = "%s"
}
`, serviceName, netName, region, subnetName, region)
}

// captureInstanceID records the instance id on first invocation and, on later
// invocations, fails if it changed — proving updates happen in place rather
// than replacing the instance. Shared by the tier-2 composition tests.
func captureInstanceID(rn string, store *string) resource.TestCheckFunc {
	return resource.TestCheckResourceAttrWith(rn, "id", func(v string) error {
		if v == "" {
			return fmt.Errorf("expected instance id to be set")
		}
		if *store == "" {
			*store = v
			return nil
		}
		if v != *store {
			return fmt.Errorf("instance was replaced: id changed from %q to %q", *store, v)
		}
		return nil
	})
}

// testAccCheckInstanceAnyNetworkAttrSet asserts that at least one element of the
// observed current_state.networks list has a non-empty value for subKey (e.g.
// "gateway_id" or "floating_ip_id"). It is order-independent, which matters
// because the API returns current_state.networks sorted by network id.
func testAccCheckInstanceAnyNetworkAttrSet(rn, subKey string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[rn]
		if !ok {
			return fmt.Errorf("not found: %s", rn)
		}
		for k, v := range rs.Primary.Attributes {
			if strings.HasPrefix(k, "current_state.networks.") && strings.HasSuffix(k, "."+subKey) && v != "" {
				return nil
			}
		}
		return fmt.Errorf("expected some current_state.networks.*.%s to be set", subKey)
	}
}

// TestAccCloudInstance_privateNetwork attaches an instance to a private vRack
// subnet only (public = false) and asserts the requested + observed network
// state, including that the platform assigned a private address.
func TestAccCloudInstance_privateNetwork(t *testing.T) {
	const rn = "ovh_cloud_instance.test"

	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	region := os.Getenv("OVH_INSTANCE_REGION_TEST")
	flavorID := os.Getenv("OVH_INSTANCE_FLAVOR_ID_TEST")
	imageID := os.Getenv("OVH_INSTANCE_IMAGE_ID_TEST")

	netName := acctest.RandomWithPrefix("tf-test-inst-net")
	subnetName := acctest.RandomWithPrefix("tf-test-inst-subnet")
	name := acctest.RandomWithPrefix("test-inst-priv")

	config := testAccCloudInstanceNetSubnetConfig(serviceName, region, netName, subnetName) + fmt.Sprintf(`
resource "ovh_cloud_instance" "test" {
  service_name = "%s"
  region       = "%s"
  name         = "%s"
  flavor_id    = "%s"
  image_id     = "%s"

  networks = [
    {
      public     = false
      network_id = ovh_cloud_network_private_vrack.net.id
      subnet_id  = ovh_cloud_network_private_vrack_subnet.subnet.id
    },
  ]
}
`, serviceName, region, name, flavorID, imageID)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckCloudInstanceNet(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(rn, "resource_status", "READY"),
					resource.TestCheckResourceAttr(rn, "networks.#", "1"),
					resource.TestCheckResourceAttr(rn, "networks.0.public", "false"),
					resource.TestCheckResourceAttrSet(rn, "networks.0.network_id"),
					resource.TestCheckResourceAttrSet(rn, "networks.0.subnet_id"),
					resource.TestCheckResourceAttrPair(rn, "networks.0.network_id", "ovh_cloud_network_private_vrack.net", "id"),
					resource.TestCheckResourceAttrPair(rn, "networks.0.subnet_id", "ovh_cloud_network_private_vrack_subnet.subnet", "id"),
					// Observed private network: subnet id echoed and an address assigned.
					resource.TestCheckResourceAttrSet(rn, "current_state.networks.0.subnet_id"),
					resource.TestCheckResourceAttrSet(rn, "current_state.networks.0.addresses.0.ip"),
				),
			},
			{
				ResourceName:            rn,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateIdFunc:       testAccCloudInstanceImportStateIdFunc(rn),
				ImportStateVerifyIgnore: []string{"checksum"},
			},
		},
	})
}

// TestAccCloudInstance_multiNIC attaches a public NIC and a private NIC to the
// same instance and asserts both requested and observed lists carry two entries.
func TestAccCloudInstance_multiNIC(t *testing.T) {
	const rn = "ovh_cloud_instance.test"

	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	region := os.Getenv("OVH_INSTANCE_REGION_TEST")
	flavorID := os.Getenv("OVH_INSTANCE_FLAVOR_ID_TEST")
	imageID := os.Getenv("OVH_INSTANCE_IMAGE_ID_TEST")

	netName := acctest.RandomWithPrefix("tf-test-inst-net")
	subnetName := acctest.RandomWithPrefix("tf-test-inst-subnet")
	name := acctest.RandomWithPrefix("test-inst-multinic")

	config := testAccCloudInstanceNetSubnetConfig(serviceName, region, netName, subnetName) + fmt.Sprintf(`
resource "ovh_cloud_instance" "test" {
  service_name = "%s"
  region       = "%s"
  name         = "%s"
  flavor_id    = "%s"
  image_id     = "%s"

  networks = [
    { public = true },
    {
      public     = false
      network_id = ovh_cloud_network_private_vrack.net.id
      subnet_id  = ovh_cloud_network_private_vrack_subnet.subnet.id
    },
  ]
}
`, serviceName, region, name, flavorID, imageID)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckCloudInstanceNet(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(rn, "resource_status", "READY"),
					resource.TestCheckResourceAttr(rn, "networks.#", "2"),
					// Requested networks carry both a public and a private entry
					// (order-independent: the API returns them sorted).
					resource.TestCheckTypeSetElemNestedAttrs(rn, "networks.*", map[string]string{"public": "true"}),
					resource.TestCheckTypeSetElemNestedAttrs(rn, "networks.*", map[string]string{"public": "false"}),
					resource.TestCheckResourceAttr(rn, "current_state.networks.#", "2"),
				),
			},
		},
	})
}

// TestAccCloudInstance_attachDetachNIC starts public-only, adds a private NIC in
// place (id stable), then removes it again — proving NIC add/remove are in-place
// updates rather than replacements.
func TestAccCloudInstance_attachDetachNIC(t *testing.T) {
	const rn = "ovh_cloud_instance.test"

	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	region := os.Getenv("OVH_INSTANCE_REGION_TEST")
	flavorID := os.Getenv("OVH_INSTANCE_FLAVOR_ID_TEST")
	imageID := os.Getenv("OVH_INSTANCE_IMAGE_ID_TEST")

	netName := acctest.RandomWithPrefix("tf-test-inst-net")
	subnetName := acctest.RandomWithPrefix("tf-test-inst-subnet")
	name := acctest.RandomWithPrefix("test-inst-nic")

	base := testAccCloudInstanceNetSubnetConfig(serviceName, region, netName, subnetName)

	publicOnly := base + fmt.Sprintf(`
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

	withPrivate := base + fmt.Sprintf(`
resource "ovh_cloud_instance" "test" {
  service_name = "%s"
  region       = "%s"
  name         = "%s"
  flavor_id    = "%s"
  image_id     = "%s"

  networks = [
    { public = true },
    {
      public     = false
      network_id = ovh_cloud_network_private_vrack.net.id
      subnet_id  = ovh_cloud_network_private_vrack_subnet.subnet.id
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
				Config: publicOnly,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(rn, "networks.#", "1"),
					captureInstanceID(rn, &instanceID),
				),
			},
			{
				// Add a private NIC — must be an in-place update.
				Config: withPrivate,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(rn, plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(rn, "networks.#", "2"),
					resource.TestCheckResourceAttr(rn, "current_state.networks.#", "2"),
					captureInstanceID(rn, &instanceID),
				),
			},
			{
				// Remove the private NIC — back to public-only, still in place.
				Config: publicOnly,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(rn, plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(rn, "networks.#", "1"),
					captureInstanceID(rn, &instanceID),
				),
			},
		},
	})
}

// TestAccCloudInstance_floatingIP wires a floating IP (same region) onto a
// public NIC through networks[].floating_ip_id and asserts the association
// surfaces in the observed state.
//
// The instance's floating_ip_id expects the OpenStack UUID of the floating IP,
// which the ovh_cloud_floating_ip resource exposes as current_state.id (its
// top-level id is the IP address).
func TestAccCloudInstance_floatingIP(t *testing.T) {
	const rn = "ovh_cloud_instance.test"

	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	region := os.Getenv("OVH_INSTANCE_REGION_TEST")
	flavorID := os.Getenv("OVH_INSTANCE_FLAVOR_ID_TEST")
	imageID := os.Getenv("OVH_INSTANCE_IMAGE_ID_TEST")

	name := acctest.RandomWithPrefix("test-inst-fip")

	config := fmt.Sprintf(`
resource "ovh_cloud_floating_ip" "fip" {
  service_name = "%s"
  region       = "%s"
  description  = "%s"
}

resource "ovh_cloud_instance" "test" {
  service_name = "%s"
  region       = "%s"
  name         = "%s"
  flavor_id    = "%s"
  image_id     = "%s"

  networks = [
    {
      public         = true
      floating_ip_id = ovh_cloud_floating_ip.fip.current_state.id
    },
  ]
}
`, serviceName, region, acctest.RandomWithPrefix("tf-test-inst-fip"), serviceName, region, name, flavorID, imageID)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckCloudInstanceNet(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(rn, "resource_status", "READY"),
					resource.TestCheckResourceAttr(rn, "networks.#", "1"),
					// The requested floating_ip_id is the floating IP OpenStack UUID.
					resource.TestCheckResourceAttrPair(rn, "networks.0.floating_ip_id", "ovh_cloud_floating_ip.fip", "current_state.id"),
					// The association is reflected in the observed state.
					resource.TestCheckResourceAttrSet(rn, "current_state.networks.0.floating_ip_id"),
				),
			},
		},
	})
}

// TestAccCloudInstance_gatewayEgress creates a network+subnet+gateway (external
// gateway enabled, subnet attached via subnet_ids) and a private instance on
// that subnet, then asserts the observed network reports the gateway id.
func TestAccCloudInstance_gatewayEgress(t *testing.T) {
	const rn = "ovh_cloud_instance.test"

	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	region := os.Getenv("OVH_INSTANCE_REGION_TEST")
	flavorID := os.Getenv("OVH_INSTANCE_FLAVOR_ID_TEST")
	imageID := os.Getenv("OVH_INSTANCE_IMAGE_ID_TEST")

	netName := acctest.RandomWithPrefix("tf-test-inst-net")
	subnetName := acctest.RandomWithPrefix("tf-test-inst-subnet")
	gwName := acctest.RandomWithPrefix("tf-test-inst-gw")
	name := acctest.RandomWithPrefix("test-inst-gw")

	config := testAccCloudInstanceNetSubnetConfig(serviceName, region, netName, subnetName) + fmt.Sprintf(`
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

resource "ovh_cloud_instance" "test" {
  service_name = "%s"
  region       = "%s"
  name         = "%s"
  flavor_id    = "%s"
  image_id     = "%s"

  networks = [
    {
      public     = false
      network_id = ovh_cloud_network_private_vrack.net.id
      subnet_id  = ovh_cloud_network_private_vrack_subnet.subnet.id
    },
  ]
}
`, gwName, region, serviceName, region, name, flavorID, imageID)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckCloudInstanceNet(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(rn, "resource_status", "READY"),
					resource.TestCheckResourceAttr(rn, "networks.0.public", "false"),
					// The gateway attached to the subnet surfaces on the observed NIC.
					resource.TestCheckResourceAttrSet(rn, "current_state.networks.0.gateway_id"),
					resource.TestCheckResourceAttrPair(rn, "current_state.networks.0.gateway_id", "ovh_cloud_gateway.gw", "id"),
				),
			},
		},
	})
}

// TestAccCloudInstance_securityGroupUpdate attaches one security group, then
// swaps it for a different one, asserting the change is applied in place and the
// requested + observed security group lists track it.
func TestAccCloudInstance_securityGroupUpdate(t *testing.T) {
	const rn = "ovh_cloud_instance.test"

	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	region := os.Getenv("OVH_INSTANCE_REGION_TEST")
	flavorID := os.Getenv("OVH_INSTANCE_FLAVOR_ID_TEST")
	imageID := os.Getenv("OVH_INSTANCE_IMAGE_ID_TEST")

	sgName1 := acctest.RandomWithPrefix("tf-test-inst-sg1")
	sgName2 := acctest.RandomWithPrefix("tf-test-inst-sg2")
	name := acctest.RandomWithPrefix("test-inst-sg")

	// Both security groups exist in every step; only the attached one changes.
	sgs := fmt.Sprintf(`
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
`, serviceName, region, sgName1, serviceName, region, sgName2)

	instance := func(sgRef string) string {
		return sgs + fmt.Sprintf(`
resource "ovh_cloud_instance" "test" {
  service_name       = "%s"
  region             = "%s"
  name               = "%s"
  flavor_id          = "%s"
  image_id           = "%s"
  security_group_ids = [%s]

  networks = [
    { public = true },
  ]
}
`, serviceName, region, name, flavorID, imageID, sgRef)
	}

	var instanceID string

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckCloudInstanceNet(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: instance("ovh_cloud_security_group.sg1.id"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(rn, "resource_status", "READY"),
					resource.TestCheckResourceAttr(rn, "security_group_ids.#", "1"),
					resource.TestCheckResourceAttrPair(rn, "security_group_ids.0", "ovh_cloud_security_group.sg1", "id"),
					resource.TestCheckResourceAttr(rn, "current_state.security_groups.#", "1"),
					captureInstanceID(rn, &instanceID),
				),
			},
			{
				// Swap to the second security group — in-place update.
				Config: instance("ovh_cloud_security_group.sg2.id"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(rn, plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(rn, "security_group_ids.#", "1"),
					resource.TestCheckResourceAttrPair(rn, "security_group_ids.0", "ovh_cloud_security_group.sg2", "id"),
					resource.TestCheckResourceAttr(rn, "current_state.security_groups.#", "1"),
					captureInstanceID(rn, &instanceID),
				),
			},
		},
	})
}
