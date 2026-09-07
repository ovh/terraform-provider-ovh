package ovh

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/ovh/go-ovh/ovh"
	"github.com/ovh/terraform-provider-ovh/v2/ovh/helpers"
)

const (
	testAccResourceCloudPrivateEndpointNamePrefix = "tf-test-pe-"
	testAccPrivateEndpointResource                = "ovh_cloud_private_endpoint.endpoint"
)

func init() {
	resource.AddTestSweepers("ovh_cloud_private_endpoint", &resource.Sweeper{
		Name: "ovh_cloud_private_endpoint",
		F:    testSweepCloudPrivateEndpoint,
	})
}

// testSweepCloudPrivateEndpoint removes the endpoints an aborted acceptance run
// left behind, releasing the ports they allocated in the subnet.
func testSweepCloudPrivateEndpoint(region string) error {
	client, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting client: %s", err)
	}

	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	if serviceName == "" {
		return nil
	}

	endpoint := "/v2/publicCloud/project/" + url.PathEscape(serviceName) + "/privateEndpoint"
	endpoints, err := helpers.GetAllPagesV2[CloudPrivateEndpointAPIResponse](context.Background(), client, endpoint)
	if err != nil {
		return fmt.Errorf("error listing private endpoints: %s", err)
	}

	for _, pe := range endpoints {
		if pe.TargetSpec == nil || !strings.HasPrefix(pe.TargetSpec.Name, testAccResourceCloudPrivateEndpointNamePrefix) {
			continue
		}
		if err := client.Delete(endpoint+"/"+url.PathEscape(pe.Id), nil); err != nil {
			return fmt.Errorf("error deleting private endpoint %s: %s", pe.Id, err)
		}
	}

	return nil
}

// testAccCheckCloudPrivateEndpointDestroy asserts every private endpoint left
// in the state is gone from the API once the test tears down.
func testAccCheckCloudPrivateEndpointDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ovh_cloud_private_endpoint" {
			continue
		}

		endpoint := "/v2/publicCloud/project/" + url.PathEscape(rs.Primary.Attributes["service_name"]) + "/privateEndpoint/" + url.PathEscape(rs.Primary.Attributes["id"])
		err := testAccOVHClient.Get(endpoint, &CloudPrivateEndpointAPIResponse{})
		if err == nil {
			return fmt.Errorf("private endpoint %s still exists", rs.Primary.Attributes["id"])
		}
		if errOvh, ok := err.(*ovh.APIError); ok && errOvh.Code == 404 {
			continue
		}
		return fmt.Errorf("error checking private endpoint %s was destroyed: %w", rs.Primary.Attributes["id"], err)
	}
	return nil
}

// privateEndpointEnv holds the inputs the private endpoint acceptance tests
// read from the environment. Fields a test does not need stay empty.
type privateEndpointEnv struct {
	serviceName string
	region      string
	service     string
	tier        string
	networkID   string
	subnetID    string
	vipAddress  string
}

func privateEndpointTestEnv() privateEndpointEnv {
	return privateEndpointEnv{
		serviceName: os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST"),
		region:      os.Getenv("OVH_CLOUD_PROJECT_REGION_TEST"),
		service:     os.Getenv("OVH_CLOUD_PRIVATE_ENDPOINT_SERVICE_TEST"),
		tier:        os.Getenv("OVH_CLOUD_PRIVATE_ENDPOINT_TIER_TEST"),
		networkID:   os.Getenv("OVH_CLOUD_PRIVATE_ENDPOINT_NETWORK_TEST"),
		subnetID:    os.Getenv("OVH_CLOUD_PRIVATE_ENDPOINT_SUBNET_TEST"),
		vipAddress:  os.Getenv("OVH_CLOUD_PRIVATE_ENDPOINT_VIP_ADDRESS_TEST"),
	}
}

// testAccPreCheckCloudPrivateEndpoint skips the acceptance tests unless
// the cloud project credentials and the private-endpoint specific inputs (a
// service, its tier and a subnet) are provided.
func testAccPreCheckCloudPrivateEndpoint(t *testing.T) {
	testAccPreCheckCloud(t)
	testAccCheckCloudProjectExists(t)
	checkEnvOrSkip(t, "OVH_CLOUD_PROJECT_REGION_TEST")
	checkEnvOrSkip(t, "OVH_CLOUD_PRIVATE_ENDPOINT_SERVICE_TEST")
	checkEnvOrSkip(t, "OVH_CLOUD_PRIVATE_ENDPOINT_TIER_TEST")
	checkEnvOrSkip(t, "OVH_CLOUD_PRIVATE_ENDPOINT_SUBNET_TEST")
}

// testAccCloudPrivateEndpointConfig renders the resource. An empty
// description, networkID, subnetID or vipAddress omits that attribute, leaving
// the platform to resolve or allocate it.
func testAccCloudPrivateEndpointConfig(env privateEndpointEnv, name, description, networkID, subnetID, vipAddress string) string {
	var optional string
	for _, attribute := range [][2]string{
		{"description", description},
		{"network_id", networkID},
		{"subnet_id", subnetID},
		{"vip_address", vipAddress},
	} {
		if attribute[1] != "" {
			optional += fmt.Sprintf("  %s = %q\n", attribute[0], attribute[1])
		}
	}

	return fmt.Sprintf(`
resource "ovh_cloud_private_endpoint" "endpoint" {
  service_name       = %q
  name               = %q
  region             = %q
  public_dns_enabled = false
%s
  service = {
    name = %q
    tier = %q
  }
}
`, env.serviceName, name, env.region, optional, env.service, env.tier)
}

func TestAccCloudPrivateEndpoint_basic(t *testing.T) {
	env := privateEndpointTestEnv()
	name := acctest.RandomWithPrefix(testAccResourceCloudPrivateEndpointNamePrefix)

	config := testAccCloudPrivateEndpointConfig(env, name, "initial", "", env.subnetID, "")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckCloudPrivateEndpoint(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudPrivateEndpointDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(testAccPrivateEndpointResource, "service_name", env.serviceName),
					resource.TestCheckResourceAttr(testAccPrivateEndpointResource, "name", name),
					resource.TestCheckResourceAttr(testAccPrivateEndpointResource, "description", "initial"),
					resource.TestCheckResourceAttr(testAccPrivateEndpointResource, "region", env.region),
					resource.TestCheckResourceAttr(testAccPrivateEndpointResource, "subnet_id", env.subnetID),
					resource.TestCheckResourceAttr(testAccPrivateEndpointResource, "service.name", env.service),
					resource.TestCheckResourceAttr(testAccPrivateEndpointResource, "service.tier", env.tier),
					resource.TestCheckResourceAttr(testAccPrivateEndpointResource, "public_dns_enabled", "false"),
					resource.TestCheckResourceAttr(testAccPrivateEndpointResource, "resource_status", "READY"),
					resource.TestCheckResourceAttrSet(testAccPrivateEndpointResource, "id"),
					resource.TestCheckResourceAttrSet(testAccPrivateEndpointResource, "checksum"),
					resource.TestCheckResourceAttrSet(testAccPrivateEndpointResource, "created_at"),
					resource.TestCheckResourceAttrSet(testAccPrivateEndpointResource, "vip_address"),
					resource.TestCheckResourceAttrPair(
						testAccPrivateEndpointResource, "vip_address",
						testAccPrivateEndpointResource, "current_state.vip_address",
					),
					resource.TestCheckResourceAttrSet(testAccPrivateEndpointResource, "current_state.vip_address"),
					resource.TestCheckResourceAttrSet(testAccPrivateEndpointResource, "current_state.network.id"),
					resource.TestCheckResourceAttrSet(testAccPrivateEndpointResource, "current_state.subnet.id"),
					resource.TestCheckResourceAttrSet(testAccPrivateEndpointResource, "current_state.service.port"),
				),
			},
			{
				ResourceName:      testAccPrivateEndpointResource,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccCloudPrivateEndpointImportStateIdFunc(testAccPrivateEndpointResource),
			},
		},
	})
}

func TestAccCloudPrivateEndpoint_update(t *testing.T) {
	env := privateEndpointTestEnv()
	name := acctest.RandomWithPrefix(testAccResourceCloudPrivateEndpointNamePrefix)
	updatedName := acctest.RandomWithPrefix(testAccResourceCloudPrivateEndpointNamePrefix)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckCloudPrivateEndpoint(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudPrivateEndpointDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudPrivateEndpointConfig(env, name, "initial", "", env.subnetID, ""),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(testAccPrivateEndpointResource, "name", name),
					resource.TestCheckResourceAttr(testAccPrivateEndpointResource, "description", "initial"),
				),
			},
			{
				Config: testAccCloudPrivateEndpointConfig(env, updatedName, "updated", "", env.subnetID, ""),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(testAccPrivateEndpointResource, "name", updatedName),
					resource.TestCheckResourceAttr(testAccPrivateEndpointResource, "description", "updated"),
					resource.TestCheckResourceAttr(testAccPrivateEndpointResource, "resource_status", "READY"),
					resource.TestCheckResourceAttrSet(testAccPrivateEndpointResource, "checksum"),
				),
			},
			// The target spec is declarative, so dropping description from the
			// configuration is expected to clear it server-side. If the API merges
			// instead of replacing, the response echoes the previous value and the
			// apply fails on an inconsistent result, which is what this step catches.
			{
				Config: testAccCloudPrivateEndpointConfig(env, updatedName, "", "", env.subnetID, ""),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckNoResourceAttr(testAccPrivateEndpointResource, "description"),
					resource.TestCheckResourceAttr(testAccPrivateEndpointResource, "resource_status", "READY"),
				),
			},
		},
	})
}

func TestAccCloudPrivateEndpoint_pinnedVipAddress(t *testing.T) {
	env := privateEndpointTestEnv()
	name := acctest.RandomWithPrefix(testAccResourceCloudPrivateEndpointNamePrefix)

	config := testAccCloudPrivateEndpointConfig(env, name, "pinned", "", env.subnetID, env.vipAddress)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCloudPrivateEndpoint(t)
			checkEnvOrSkip(t, "OVH_CLOUD_PRIVATE_ENDPOINT_VIP_ADDRESS_TEST")
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudPrivateEndpointDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(testAccPrivateEndpointResource, "vip_address", env.vipAddress),
					resource.TestCheckResourceAttr(testAccPrivateEndpointResource, "current_state.vip_address", env.vipAddress),
					resource.TestCheckResourceAttr(testAccPrivateEndpointResource, "resource_status", "READY"),
					resource.TestCheckResourceAttrSet(testAccPrivateEndpointResource, "current_state.subnet.id"),
				),
			},
		},
	})
}

// TestAccCloudPrivateEndpoint_placementNetworkOnly asserts the platform
// resolves the subnet and that the resolved value surfaces on the root attribute
// without leaving a diff behind.
func TestAccCloudPrivateEndpoint_placementNetworkOnly(t *testing.T) {
	env := privateEndpointTestEnv()
	name := acctest.RandomWithPrefix(testAccResourceCloudPrivateEndpointNamePrefix)

	config := testAccCloudPrivateEndpointConfig(env, name, "", env.networkID, "", "")

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCloudPrivateEndpoint(t)
			checkEnvOrSkip(t, "OVH_CLOUD_PRIVATE_ENDPOINT_NETWORK_TEST")
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudPrivateEndpointDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(testAccPrivateEndpointResource, "network_id", env.networkID),
					resource.TestCheckResourceAttrSet(testAccPrivateEndpointResource, "subnet_id"),
					resource.TestCheckResourceAttrPair(
						testAccPrivateEndpointResource, "subnet_id",
						testAccPrivateEndpointResource, "current_state.subnet.id",
					),
					resource.TestCheckResourceAttrPair(
						testAccPrivateEndpointResource, "network_id",
						testAccPrivateEndpointResource, "current_state.network.id",
					),
					resource.TestCheckResourceAttr(testAccPrivateEndpointResource, "resource_status", "READY"),
				),
			},
			{
				Config:   config,
				PlanOnly: true,
			},
		},
	})
}

// TestAccCloudPrivateEndpoint_placementSubnetOnly asserts the platform
// resolves the parent network onto the root attribute.
func TestAccCloudPrivateEndpoint_placementSubnetOnly(t *testing.T) {
	env := privateEndpointTestEnv()
	name := acctest.RandomWithPrefix(testAccResourceCloudPrivateEndpointNamePrefix)

	config := testAccCloudPrivateEndpointConfig(env, name, "", "", env.subnetID, "")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckCloudPrivateEndpoint(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudPrivateEndpointDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(testAccPrivateEndpointResource, "subnet_id", env.subnetID),
					resource.TestCheckResourceAttrSet(testAccPrivateEndpointResource, "network_id"),
					resource.TestCheckResourceAttrPair(
						testAccPrivateEndpointResource, "network_id",
						testAccPrivateEndpointResource, "current_state.network.id",
					),
					resource.TestCheckResourceAttr(testAccPrivateEndpointResource, "resource_status", "READY"),
				),
			},
			{
				Config:   config,
				PlanOnly: true,
			},
		},
	})
}

// TestAccCloudPrivateEndpoint_placementBoth asserts the both-set placement
// the API supports is accepted by the provider (ExactlyOneOf used to reject it).
func TestAccCloudPrivateEndpoint_placementBoth(t *testing.T) {
	env := privateEndpointTestEnv()
	name := acctest.RandomWithPrefix(testAccResourceCloudPrivateEndpointNamePrefix)

	config := testAccCloudPrivateEndpointConfig(env, name, "", env.networkID, env.subnetID, "")

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCloudPrivateEndpoint(t)
			checkEnvOrSkip(t, "OVH_CLOUD_PRIVATE_ENDPOINT_NETWORK_TEST")
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudPrivateEndpointDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(testAccPrivateEndpointResource, "network_id", env.networkID),
					resource.TestCheckResourceAttr(testAccPrivateEndpointResource, "subnet_id", env.subnetID),
					resource.TestCheckResourceAttr(testAccPrivateEndpointResource, "current_state.network.id", env.networkID),
					resource.TestCheckResourceAttr(testAccPrivateEndpointResource, "current_state.subnet.id", env.subnetID),
					resource.TestCheckResourceAttr(testAccPrivateEndpointResource, "resource_status", "READY"),
				),
			},
			{
				Config:   config,
				PlanOnly: true,
			},
		},
	})
}

// TestAccCloudPrivateEndpoint_placementPinResolvedSubnet asserts that
// promoting the platform-resolved subnet into the configuration is a no-op:
// RequiresReplaceIfConfigured only replaces when a configured value differs from
// state. Assumes OVH_..._NETWORK_TEST holds exactly OVH_..._SUBNET_TEST as its
// single IPv4 subnet, which is also what network-only placement requires.
func TestAccCloudPrivateEndpoint_placementPinResolvedSubnet(t *testing.T) {
	env := privateEndpointTestEnv()
	name := acctest.RandomWithPrefix(testAccResourceCloudPrivateEndpointNamePrefix)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCloudPrivateEndpoint(t)
			checkEnvOrSkip(t, "OVH_CLOUD_PRIVATE_ENDPOINT_NETWORK_TEST")
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudPrivateEndpointDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudPrivateEndpointConfig(env, name, "", env.networkID, "", ""),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(testAccPrivateEndpointResource, "subnet_id", env.subnetID),
				),
			},
			{
				Config: testAccCloudPrivateEndpointConfig(env, name, "", env.networkID, env.subnetID, ""),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(testAccPrivateEndpointResource, plancheck.ResourceActionNoop),
					},
				},
			},
		},
	})
}

func testAccCloudPrivateEndpointImportStateIdFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}
		return fmt.Sprintf("%s/%s", rs.Primary.Attributes["service_name"], rs.Primary.Attributes["id"]), nil
	}
}
