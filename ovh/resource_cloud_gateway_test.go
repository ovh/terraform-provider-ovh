package ovh

import (
	"fmt"
	"net/url"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/ovh/go-ovh/ovh"
)

const testAccResourceCloudGatewayNamePrefix = "tf-test-gateway-v2-"

func testAccPreCheckCloudGateway(t *testing.T) {
	testAccPreCheckCredentials(t)
	if os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST") == "" {
		t.Skip("OVH_CLOUD_PROJECT_SERVICE_TEST must be set for acceptance tests")
	}
	if os.Getenv("OVH_CLOUD_PROJECT_REGION_TEST") == "" {
		t.Skip("OVH_CLOUD_PROJECT_REGION_TEST must be set for acceptance tests")
	}
}

func testAccCloudGatewayImportStateIdFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}
		return fmt.Sprintf("%s/%s", rs.Primary.Attributes["service_name"], rs.Primary.Attributes["id"]), nil
	}
}

// testAccCheckCloudGatewayDestroy verifies that every gateway managed by the
// test has actually been removed from the API once the test completes (the GET
// must return a 404).
func testAccCheckCloudGatewayDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ovh_cloud_gateway" {
			continue
		}

		serviceName := rs.Primary.Attributes["service_name"]
		gatewayID := rs.Primary.Attributes["id"]

		endpoint := "/v2/publicCloud/project/" + url.PathEscape(serviceName) + "/gateway/" + url.PathEscape(gatewayID)

		err := testAccOVHClient.Get(endpoint, nil)
		if err == nil {
			return fmt.Errorf("cloud gateway %s still exists", gatewayID)
		}
		if errOvh, ok := err.(*ovh.APIError); ok && errOvh.Code == 404 {
			continue
		}
		return fmt.Errorf("error checking that cloud gateway %s was destroyed: %w", gatewayID, err)
	}

	return nil
}

// testAccCloudGatewayConfig builds an HCL config that provisions a private
// network + subnet and a gateway attached to that subnet.
//
// Note: the schema uses terraform-plugin-framework SingleNestedAttribute /
// ListAttribute, so the HCL MUST use attribute assignment (external_gateway =
// { ... }, subnet_ids = [ ... ]) rather than block syntax.
func testAccCloudGatewayConfig(serviceName, region, networkName, subnetName, gatewayName, model string) string {
	return fmt.Sprintf(`
resource "ovh_cloud_network_private_vrack" "network" {
  service_name = "%s"
  name         = "%s"
  region       = "%s"
}

resource "ovh_cloud_network_private_vrack_subnet" "subnet" {
  service_name = ovh_cloud_network_private_vrack.network.service_name
  network_id   = ovh_cloud_network_private_vrack.network.id
  name         = "%s"
  cidr         = "10.0.0.0/24"
  region       = "%s"
  dhcp_enabled = true
}

resource "ovh_cloud_gateway" "test" {
  service_name = "%s"
  name         = "%s"
  region       = "%s"

  external_gateway = {
    enabled = true
    model   = "%s"
  }

  subnet_ids = [ovh_cloud_network_private_vrack_subnet.subnet.id]
}
`, serviceName, networkName, region, subnetName, region, serviceName, gatewayName, region, model)
}

func TestAccCloudGateway_basic(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	region := os.Getenv("OVH_CLOUD_PROJECT_REGION_TEST")
	gatewayModel := os.Getenv("OVH_CLOUD_GATEWAY_MODEL_TEST")
	if gatewayModel == "" {
		gatewayModel = "S"
	}

	networkName := acctest.RandomWithPrefix(testAccResourceCloudNetworkPrivateVrackNamePrefix)
	subnetName := acctest.RandomWithPrefix(testAccResourceCloudNetworkPrivateVrackSubnetNamePrefix)
	gatewayName := acctest.RandomWithPrefix(testAccResourceCloudGatewayNamePrefix)

	config := testAccCloudGatewayConfig(serviceName, region, networkName, subnetName, gatewayName, gatewayModel)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCloudGateway(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudGatewayDestroy,
		Steps: []resource.TestStep{
			// Step 1: create the gateway with external_gateway.model = "S".
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_gateway.test", "service_name", serviceName),
					resource.TestCheckResourceAttr("ovh_cloud_gateway.test", "name", gatewayName),
					resource.TestCheckResourceAttr("ovh_cloud_gateway.test", "region", region),
					resource.TestCheckResourceAttr("ovh_cloud_gateway.test", "external_gateway.enabled", "true"),
					resource.TestCheckResourceAttr("ovh_cloud_gateway.test", "external_gateway.model", gatewayModel),
					resource.TestCheckResourceAttr("ovh_cloud_gateway.test", "subnet_ids.#", "1"),
					resource.TestCheckResourceAttrSet("ovh_cloud_gateway.test", "id"),
					resource.TestCheckResourceAttrSet("ovh_cloud_gateway.test", "checksum"),
					resource.TestCheckResourceAttrSet("ovh_cloud_gateway.test", "created_at"),
					resource.TestCheckResourceAttr("ovh_cloud_gateway.test", "resource_status", "READY"),
					resource.TestCheckResourceAttrSet("ovh_cloud_gateway.test", "current_state.name"),
					resource.TestCheckResourceAttr("ovh_cloud_gateway.test", "current_state.external_gateway.enabled", "true"),
					resource.TestCheckResourceAttrSet("ovh_cloud_gateway.test", "current_state.status"),
				),
			},
			// Step 2: REGRESSION GUARD for fix/cloud-gateway-model-perpetual-diff.
			//
			// Re-applying the SAME config (external_gateway.model = "S") must
			// produce an EMPTY plan. Before the fix, MergeWith clobbered
			// external_gateway.model to null because the API does not always
			// echo the model back, which caused a perpetual diff and a gateway
			// update on every apply. ExpectEmptyPlan locks the fix in.
			{
				Config: config,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_gateway.test", "external_gateway.model", gatewayModel),
				),
			},
			// Step 3: import.
			{
				ResourceName:      "ovh_cloud_gateway.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccCloudGatewayImportStateIdFunc("ovh_cloud_gateway.test"),
			},
		},
	})
}

// TestAccCloudGateway_noPerpetualDiff is a focused regression test for
// fix/cloud-gateway-model-perpetual-diff. It creates the gateway then runs a
// PlanOnly step on the identical config: a non-empty plan fails the step.
// This proves external_gateway.model no longer drifts to null between applies.
func TestAccCloudGateway_noPerpetualDiff(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	region := os.Getenv("OVH_CLOUD_PROJECT_REGION_TEST")
	gatewayModel := os.Getenv("OVH_CLOUD_GATEWAY_MODEL_TEST")
	if gatewayModel == "" {
		gatewayModel = "S"
	}

	networkName := acctest.RandomWithPrefix(testAccResourceCloudNetworkPrivateVrackNamePrefix)
	subnetName := acctest.RandomWithPrefix(testAccResourceCloudNetworkPrivateVrackSubnetNamePrefix)
	gatewayName := acctest.RandomWithPrefix(testAccResourceCloudGatewayNamePrefix)

	config := testAccCloudGatewayConfig(serviceName, region, networkName, subnetName, gatewayName, gatewayModel)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCloudGateway(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_gateway.test", "external_gateway.model", gatewayModel),
					resource.TestCheckResourceAttr("ovh_cloud_gateway.test", "resource_status", "READY"),
				),
			},
			{
				// PlanOnly re-plans the identical config and fails if the plan
				// is not empty — the canonical "no perpetual diff" assertion.
				Config:             config,
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
		},
	})
}

// TestAccCloudGateway_update renames the gateway in place (name is mutable) and
// verifies the gateway stays READY and keeps the same id.
func TestAccCloudGateway_update(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	region := os.Getenv("OVH_CLOUD_PROJECT_REGION_TEST")
	gatewayModel := os.Getenv("OVH_CLOUD_GATEWAY_MODEL_TEST")
	if gatewayModel == "" {
		gatewayModel = "S"
	}

	networkName := acctest.RandomWithPrefix(testAccResourceCloudNetworkPrivateVrackNamePrefix)
	subnetName := acctest.RandomWithPrefix(testAccResourceCloudNetworkPrivateVrackSubnetNamePrefix)
	gatewayName := acctest.RandomWithPrefix(testAccResourceCloudGatewayNamePrefix)
	updatedName := acctest.RandomWithPrefix(testAccResourceCloudGatewayNamePrefix)

	config := testAccCloudGatewayConfig(serviceName, region, networkName, subnetName, gatewayName, gatewayModel)
	updatedConfig := testAccCloudGatewayConfig(serviceName, region, networkName, subnetName, updatedName, gatewayModel)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCloudGateway(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_gateway.test", "name", gatewayName),
					resource.TestCheckResourceAttr("ovh_cloud_gateway.test", "external_gateway.enabled", "true"),
					resource.TestCheckResourceAttr("ovh_cloud_gateway.test", "external_gateway.model", gatewayModel),
					resource.TestCheckResourceAttrSet("ovh_cloud_gateway.test", "id"),
					resource.TestCheckResourceAttrSet("ovh_cloud_gateway.test", "checksum"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_gateway.test", "name", updatedName),
					resource.TestCheckResourceAttr("ovh_cloud_gateway.test", "external_gateway.model", gatewayModel),
					resource.TestCheckResourceAttr("ovh_cloud_gateway.test", "resource_status", "READY"),
					resource.TestCheckResourceAttrSet("ovh_cloud_gateway.test", "id"),
					resource.TestCheckResourceAttrSet("ovh_cloud_gateway.test", "checksum"),
				),
			},
		},
	})
}
