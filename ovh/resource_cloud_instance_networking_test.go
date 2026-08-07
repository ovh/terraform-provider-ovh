package ovh

import (
	"fmt"
	"os"
	"strconv"
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
  gateway_ip   = "10.0.0.1"
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
// "gateway_id"). It is order-independent, which matters because the API returns
// current_state.networks sorted by network id. Only direct children of a list
// element match; use testAccCheckInstanceAnyNetworkAddress and its wrappers to
// assert on the nested addresses list.
func testAccCheckInstanceAnyNetworkAttrSet(rn, subKey string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[rn]
		if !ok {
			return fmt.Errorf("not found: %s", rn)
		}
		for k, v := range rs.Primary.Attributes {
			if v == "" {
				continue
			}
			if !strings.HasPrefix(k, "current_state.networks.") {
				continue
			}
			idx, sub, found := strings.Cut(strings.TrimPrefix(k, "current_state.networks."), ".")
			if !found || sub != subKey {
				continue
			}
			if _, err := strconv.Atoi(idx); err != nil {
				continue
			}
			return nil
		}
		return fmt.Errorf("expected some current_state.networks.*.%s to be set", subKey)
	}
}

// testAccCheckInstanceAnyNetworkAddress asserts that some address of some
// observed current_state.networks entry satisfies match. Both index levels are
// walked because the API sorts entries by network id and addresses by
// (type, ip); want describes the expectation for the failure message.
func testAccCheckInstanceAnyNetworkAddress(rn, want string, match func(ip, addrType string) bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[rn]
		if !ok {
			return fmt.Errorf("not found: %s", rn)
		}
		type address struct{ ip, addrType string }
		addresses := map[string]*address{}
		for k, v := range rs.Primary.Attributes {
			rest, ok := strings.CutPrefix(k, "current_state.networks.")
			if !ok {
				continue
			}
			netIdx, rest, ok := strings.Cut(rest, ".addresses.")
			if !ok {
				continue
			}
			addrIdx, sub, ok := strings.Cut(rest, ".")
			if !ok {
				continue
			}
			if _, err := strconv.Atoi(netIdx); err != nil {
				continue
			}
			if _, err := strconv.Atoi(addrIdx); err != nil {
				continue
			}
			key := netIdx + "." + addrIdx
			if addresses[key] == nil {
				addresses[key] = &address{}
			}
			switch sub {
			case "ip":
				addresses[key].ip = v
			case "type":
				addresses[key].addrType = v
			}
		}
		for _, a := range addresses {
			if match(a.ip, a.addrType) {
				return nil
			}
		}
		return fmt.Errorf("expected some current_state.networks.*.addresses.* to be %s", want)
	}
}

func testAccCheckInstanceAnyNetworkAddressIP(rn, wantIP string) resource.TestCheckFunc {
	return testAccCheckInstanceAnyNetworkAddress(rn, "the address "+wantIP, func(ip, _ string) bool {
		return ip == wantIP
	})
}

// testAccCheckInstanceAnyNetworkAddressPair is the addresses[] equivalent of
// resource.TestCheckResourceAttrPair: it asserts some observed address carries
// addrType and the IP held by srcResource's srcAttr.
func testAccCheckInstanceAnyNetworkAddressPair(rn, addrType, srcResource, srcAttr string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		src, ok := s.RootModule().Resources[srcResource]
		if !ok {
			return fmt.Errorf("not found: %s", srcResource)
		}
		wantIP := src.Primary.Attributes[srcAttr]
		if wantIP == "" {
			return fmt.Errorf("expected %s.%s to be set", srcResource, srcAttr)
		}
		want := fmt.Sprintf("the address %s with type %s", wantIP, addrType)
		return testAccCheckInstanceAnyNetworkAddress(rn, want, func(ip, t string) bool {
			return ip == wantIP && t == addrType
		})(s)
	}
}

// TestAccCloudInstance_privateNetwork attaches an instance to a private vRack
// subnet only (no public IP) and asserts the requested + observed network
// state, including that the platform assigned a private address.
func TestAccCloudInstance_privateNetwork(t *testing.T) {
	const rn = "ovh_cloud_instance.test"

	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	region := os.Getenv("OVH_CLOUD_PROJECT_REGION_TEST")
	flavorID := resolveInstanceFlavorID(t, serviceName, region, testAccInstanceFlavorName)
	imageID := resolveInstanceImageID(t, serviceName, region, testAccInstanceImageName)

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
					resource.TestCheckNoResourceAttr(rn, "networks.0.auto_assign_public_ip"),
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
				ResourceName:      rn,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccCloudInstanceImportStateIdFunc(rn),
			},
		},
	})
}

// TestAccCloudInstance_multiNIC attaches a public NIC and a private NIC to the
// same instance and asserts both requested and observed lists carry two entries.
func TestAccCloudInstance_multiNIC(t *testing.T) {
	const rn = "ovh_cloud_instance.test"

	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	region := os.Getenv("OVH_CLOUD_PROJECT_REGION_TEST")
	flavorID := resolveInstanceFlavorID(t, serviceName, region, testAccInstanceFlavorName)
	imageID := resolveInstanceImageID(t, serviceName, region, testAccInstanceImageName)

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
    { auto_assign_public_ip = true },
    {
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
					resource.TestCheckTypeSetElemNestedAttrs(rn, "networks.*", map[string]string{"auto_assign_public_ip": "true"}),
					resource.TestCheckTypeSetElemAttrPair(rn, "networks.*.network_id", "ovh_cloud_network_private_vrack.net", "id"),
					resource.TestCheckResourceAttr(rn, "current_state.networks.#", "2"),
				),
			},
			{
				// networks[].ip and networks[].auto_assign_public_ip are Optional-only:
				// re-planning the same config must produce no diff even though the API
				// echoes back an observed address in current_state.
				Config:   config,
				PlanOnly: true,
			},
		},
	})
}

// TestAccCloudInstance_multiNICPrivateFirst is the ordering regression guard:
// the config lists the private NIC first while apiv2 sorts the derived target
// spec by network UUID and returns the id-less public entry first. `networks`
// is Optional-only, so the applied value must match the config index-for-index
// or the apply fails with "Provider produced inconsistent result after apply".
func TestAccCloudInstance_multiNICPrivateFirst(t *testing.T) {
	const rn = "ovh_cloud_instance.test"

	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	region := os.Getenv("OVH_CLOUD_PROJECT_REGION_TEST")
	flavorID := resolveInstanceFlavorID(t, serviceName, region, testAccInstanceFlavorName)
	imageID := resolveInstanceImageID(t, serviceName, region, testAccInstanceImageName)

	netName := acctest.RandomWithPrefix("tf-test-inst-net")
	subnetName := acctest.RandomWithPrefix("tf-test-inst-subnet")
	name := acctest.RandomWithPrefix("test-inst-privfirst")

	config := testAccCloudInstanceNetSubnetConfig(serviceName, region, netName, subnetName) + fmt.Sprintf(`
resource "ovh_cloud_instance" "test" {
  service_name = "%s"
  region       = "%s"
  name         = "%s"
  flavor_id    = "%s"
  image_id     = "%s"

  networks = [
    {
      network_id = ovh_cloud_network_private_vrack.net.id
      subnet_id  = ovh_cloud_network_private_vrack_subnet.subnet.id
    },
    { auto_assign_public_ip = true },
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
					// Index-for-index: the private entry stays at 0, the public at 1.
					resource.TestCheckResourceAttrPair(rn, "networks.0.network_id", "ovh_cloud_network_private_vrack.net", "id"),
					resource.TestCheckResourceAttrPair(rn, "networks.0.subnet_id", "ovh_cloud_network_private_vrack_subnet.subnet", "id"),
					resource.TestCheckNoResourceAttr(rn, "networks.0.auto_assign_public_ip"),
					resource.TestCheckNoResourceAttr(rn, "networks.1.network_id"),
					resource.TestCheckResourceAttr(rn, "networks.1.auto_assign_public_ip", "true"),
					resource.TestCheckResourceAttr(rn, "current_state.networks.#", "2"),
				),
			},
			{
				// Refresh must re-order the API response back to the config order too.
				Config:   config,
				PlanOnly: true,
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
	region := os.Getenv("OVH_CLOUD_PROJECT_REGION_TEST")
	flavorID := resolveInstanceFlavorID(t, serviceName, region, testAccInstanceFlavorName)
	imageID := resolveInstanceImageID(t, serviceName, region, testAccInstanceImageName)

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
    { auto_assign_public_ip = true },
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
    { auto_assign_public_ip = true },
    {
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
// private NIC through networks[].ip and asserts the association surfaces in the
// observed state. apiv2 only associates a floating IP to a private port, not to
// the Ext-Net/public port. The address lies outside the subnet CIDR, so it
// associates the existing floating IP rather than pinning a fixed IP.
//
// networks[].ip expects the address of an existing floating IP, which the
// ovh_cloud_floating_ip resource exposes as its top-level id.
func TestAccCloudInstance_floatingIP(t *testing.T) {
	const rn = "ovh_cloud_instance.test"

	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	region := os.Getenv("OVH_CLOUD_PROJECT_REGION_TEST")
	flavorID := resolveInstanceFlavorID(t, serviceName, region, testAccInstanceFlavorName)
	imageID := resolveInstanceImageID(t, serviceName, region, testAccInstanceImageName)

	netName := acctest.RandomWithPrefix("tf-test-inst-net")
	subnetName := acctest.RandomWithPrefix("tf-test-inst-subnet")
	name := acctest.RandomWithPrefix("test-inst-fip")

	gwName := acctest.RandomWithPrefix("tf-test-inst-fip-gw")

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
      network_id = ovh_cloud_network_private_vrack.net.id
      subnet_id  = ovh_cloud_network_private_vrack_subnet.subnet.id
      ip         = ovh_cloud_floating_ip.fip.id
    },
  ]

  # A floating IP can only be associated once the subnet is attached to a gateway
  # with an external gateway enabled, so the instance must be created after it.
  depends_on = [ovh_cloud_gateway.gw]
}
`, gwName, region, serviceName, region, acctest.RandomWithPrefix("tf-test-inst-fip"), serviceName, region, name, flavorID, imageID)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckCloudInstanceNet(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(rn, "resource_status", "READY"),
					resource.TestCheckResourceAttr(rn, "networks.#", "1"),
					// The requested ip is the floating IP address.
					resource.TestCheckResourceAttrPair(rn, "networks.0.ip", "ovh_cloud_floating_ip.fip", "id"),
					// The association is reflected in the observed state as a FLOATING
					// address on the private interface (order-independent: the API sorts
					// current_state.networks by network id and addresses by (type, ip)).
					testAccCheckInstanceAnyNetworkAddressPair(rn, "FLOATING", "ovh_cloud_floating_ip.fip", "id"),
				),
			},
			{
				// Regression guard for the Optional-only networks[].ip: the observed
				// address is reported in current_state.networks[].addresses[] and never
				// echoed back into networks[].ip, so re-planning the same config must be
				// a no-op.
				Config:   config,
				PlanOnly: true,
			},
		},
	})
}

// TestAccCloudInstance_ownedPublicIP boots an instance whose only interface is
// {ip = <public IP the project already owns>} — an additional IP or an
// Ext-Net IP of the project in the instance region — and asserts the address
// round-trips with no drift on a second plan. This is the regression guard for
// declaring networks[].ip Optional-only.
func TestAccCloudInstance_ownedPublicIP(t *testing.T) {
	const rn = "ovh_cloud_instance.test"

	ownedIP := os.Getenv("OVH_CLOUD_PROJECT_OWNED_PUBLIC_IP_TEST")
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	region := os.Getenv("OVH_CLOUD_PROJECT_REGION_TEST")
	flavorID := resolveInstanceFlavorID(t, serviceName, region, testAccInstanceFlavorName)
	imageID := resolveInstanceImageID(t, serviceName, region, testAccInstanceImageName)

	name := acctest.RandomWithPrefix("test-inst-ownedip")

	config := fmt.Sprintf(`
resource "ovh_cloud_instance" "test" {
  service_name = "%s"
  region       = "%s"
  name         = "%s"
  flavor_id    = "%s"
  image_id     = "%s"

  networks = [
    { ip = "%s" },
  ]
}
`, serviceName, region, name, flavorID, imageID, ownedIP)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCloudInstanceNet(t)
			checkEnvOrSkip(t, "OVH_CLOUD_PROJECT_OWNED_PUBLIC_IP_TEST")
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(rn, "resource_status", "READY"),
					resource.TestCheckResourceAttr(rn, "networks.#", "1"),
					resource.TestCheckResourceAttr(rn, "networks.0.ip", ownedIP),
					resource.TestCheckNoResourceAttr(rn, "networks.0.network_id"),
					resource.TestCheckNoResourceAttr(rn, "networks.0.auto_assign_public_ip"),
					// The owned address is observed on the public interface, typed FIXED
					// when it is an Ext-Net address and ADDITIONAL when it is routed.
					testAccCheckInstanceAnyNetworkAddressIP(rn, ownedIP),
				),
			},
			{
				Config:   config,
				PlanOnly: true,
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
	region := os.Getenv("OVH_CLOUD_PROJECT_REGION_TEST")
	flavorID := resolveInstanceFlavorID(t, serviceName, region, testAccInstanceFlavorName)
	imageID := resolveInstanceImageID(t, serviceName, region, testAccInstanceImageName)

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
					resource.TestCheckNoResourceAttr(rn, "networks.0.auto_assign_public_ip"),
					// The gateway attached to the subnet surfaces on the observed NIC.
					resource.TestCheckResourceAttrSet(rn, "current_state.networks.0.gateway_id"),
					resource.TestCheckResourceAttrPair(rn, "current_state.networks.0.gateway_id", "ovh_cloud_gateway.gw", "id"),
				),
			},
		},
	})
}

// TestAccCloudInstance_securityGroupUpdate attaches one security group, swaps it
// for a different one, then detaches every group with an explicit empty list,
// asserting each change is applied in place and the requested + observed security
// group lists track it.
func TestAccCloudInstance_securityGroupUpdate(t *testing.T) {
	const rn = "ovh_cloud_instance.test"

	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	region := os.Getenv("OVH_CLOUD_PROJECT_REGION_TEST")
	flavorID := resolveInstanceFlavorID(t, serviceName, region, testAccInstanceFlavorName)
	imageID := resolveInstanceImageID(t, serviceName, region, testAccInstanceImageName)

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
    { auto_assign_public_ip = true },
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
			{
				// An explicit empty list detaches every security group. Omitting the
				// attribute instead would keep the current groups (the API leaves them
				// unchanged when securityGroups is absent from the update target spec).
				Config: instance(""),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(rn, plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(rn, "security_group_ids.#", "0"),
					resource.TestCheckNoResourceAttr(rn, "current_state.security_groups.0.id"),
					captureInstanceID(rn, &instanceID),
				),
			},
			{
				// The empty list must be idempotent, not re-planned on every apply.
				Config: instance(""),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()},
				},
			},
		},
	})
}
