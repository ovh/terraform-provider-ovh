package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

const testAccResourceCloudLoadbalancerNamePrefix = "tf-test-lb-v2-"

func TestAccCloudLoadbalancer_basic(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	region := os.Getenv("OVH_CLOUD_PROJECT_LOADBALANCER_REGION_TEST")
	vipNetworkId := os.Getenv("OVH_CLOUD_PROJECT_LOADBALANCER_VIP_NETWORK_ID_TEST")
	vipSubnetId := os.Getenv("OVH_CLOUD_PROJECT_LOADBALANCER_VIP_SUBNET_ID_TEST")
	flavorName := "SMALL"

	lbName := acctest.RandomWithPrefix(testAccResourceCloudLoadbalancerNamePrefix)

	config := fmt.Sprintf(`
resource "ovh_cloud_loadbalancer" "test" {
  service_name = "%s"
  name         = "%s"
  region       = "%s"
  flavor_name  = "%s"

  network = {
    id        = "%s"
    subnet_id = "%s"
  }
}
`, serviceName, lbName, region, flavorName, vipNetworkId, vipSubnetId)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCloudLoadbalancer(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_loadbalancer.test", "service_name", serviceName),
					resource.TestCheckResourceAttr("ovh_cloud_loadbalancer.test", "name", lbName),
					resource.TestCheckResourceAttr("ovh_cloud_loadbalancer.test", "region", region),
					resource.TestCheckResourceAttr("ovh_cloud_loadbalancer.test", "network.id", vipNetworkId),
					resource.TestCheckResourceAttr("ovh_cloud_loadbalancer.test", "network.subnet_id", vipSubnetId),
					resource.TestCheckNoResourceAttr("ovh_cloud_loadbalancer.test", "network.ip"),
					resource.TestCheckResourceAttr("ovh_cloud_loadbalancer.test", "flavor_name", flavorName),
					resource.TestCheckResourceAttrSet("ovh_cloud_loadbalancer.test", "id"),
					resource.TestCheckResourceAttrSet("ovh_cloud_loadbalancer.test", "checksum"),
					resource.TestCheckResourceAttrSet("ovh_cloud_loadbalancer.test", "created_at"),
					resource.TestCheckResourceAttr("ovh_cloud_loadbalancer.test", "resource_status", "READY"),
					resource.TestCheckResourceAttrSet("ovh_cloud_loadbalancer.test", "current_state.name"),
					resource.TestCheckResourceAttr("ovh_cloud_loadbalancer.test", "current_state.network.id", vipNetworkId),
					resource.TestCheckResourceAttr("ovh_cloud_loadbalancer.test", "current_state.network.subnet_id", vipSubnetId),
					resource.TestCheckResourceAttr("ovh_cloud_loadbalancer.test", "current_state.network.addresses.0.type", "FIXED"),
					resource.TestCheckResourceAttrSet("ovh_cloud_loadbalancer.test", "current_state.network.addresses.0.ip"),
					resource.TestCheckResourceAttrSet("ovh_cloud_loadbalancer.test", "current_state.operating_status"),
					resource.TestCheckResourceAttrSet("ovh_cloud_loadbalancer.test", "current_state.provisioning_status"),
				),
			},
			// Test import
			{
				ResourceName:      "ovh_cloud_loadbalancer.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccCloudLoadbalancerImportStateIdFunc("ovh_cloud_loadbalancer.test"),
			},
		},
	})
}

func TestAccCloudLoadbalancer_pinnedVipIp(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	region := os.Getenv("OVH_CLOUD_PROJECT_LOADBALANCER_REGION_TEST")
	vipNetworkId := os.Getenv("OVH_CLOUD_PROJECT_LOADBALANCER_VIP_NETWORK_ID_TEST")
	vipSubnetId := os.Getenv("OVH_CLOUD_PROJECT_LOADBALANCER_VIP_SUBNET_ID_TEST")
	vipIp := os.Getenv("OVH_CLOUD_PROJECT_LOADBALANCER_VIP_IP_TEST")
	flavorName := "SMALL"

	lbName := acctest.RandomWithPrefix(testAccResourceCloudLoadbalancerNamePrefix)

	config := fmt.Sprintf(`
resource "ovh_cloud_loadbalancer" "test" {
  service_name = "%s"
  name         = "%s"
  region       = "%s"
  flavor_name  = "%s"

  network = {
    id        = "%s"
    subnet_id = "%s"
    ip        = "%s"
  }
}
`, serviceName, lbName, region, flavorName, vipNetworkId, vipSubnetId, vipIp)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCloudLoadbalancer(t)
			if vipIp == "" {
				t.Skip("OVH_CLOUD_PROJECT_LOADBALANCER_VIP_IP_TEST not set")
			}
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_loadbalancer.test", "network.ip", vipIp),
					resource.TestCheckResourceAttr("ovh_cloud_loadbalancer.test", "current_state.network.addresses.0.ip", vipIp),
					resource.TestCheckResourceAttr("ovh_cloud_loadbalancer.test", "current_state.network.addresses.0.type", "FIXED"),
					resource.TestCheckResourceAttr("ovh_cloud_loadbalancer.test", "resource_status", "READY"),
				),
			},
		},
	})
}

func TestAccCloudLoadbalancer_update(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	region := os.Getenv("OVH_CLOUD_PROJECT_LOADBALANCER_REGION_TEST")
	vipNetworkId := os.Getenv("OVH_CLOUD_PROJECT_LOADBALANCER_VIP_NETWORK_ID_TEST")
	vipSubnetId := os.Getenv("OVH_CLOUD_PROJECT_LOADBALANCER_VIP_SUBNET_ID_TEST")
	flavorName := "SMALL"

	lbName := acctest.RandomWithPrefix(testAccResourceCloudLoadbalancerNamePrefix)
	updatedName := acctest.RandomWithPrefix(testAccResourceCloudLoadbalancerNamePrefix)

	configTemplate := `
resource "ovh_cloud_loadbalancer" "test" {
  service_name = "%s"
  name         = "%s"
  region       = "%s"
  flavor_name  = "%s"
  description  = "%s"

  network = {
    id        = "%s"
    subnet_id = "%s"
  }
}
`

	config := fmt.Sprintf(configTemplate, serviceName, lbName, region, flavorName, "initial description", vipNetworkId, vipSubnetId)
	updatedConfig := fmt.Sprintf(configTemplate, serviceName, updatedName, region, flavorName, "updated description", vipNetworkId, vipSubnetId)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCloudLoadbalancer(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_loadbalancer.test", "name", lbName),
					resource.TestCheckResourceAttr("ovh_cloud_loadbalancer.test", "description", "initial description"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_loadbalancer.test", "name", updatedName),
					resource.TestCheckResourceAttr("ovh_cloud_loadbalancer.test", "description", "updated description"),
					resource.TestCheckResourceAttr("ovh_cloud_loadbalancer.test", "network.id", vipNetworkId),
					resource.TestCheckResourceAttr("ovh_cloud_loadbalancer.test", "network.subnet_id", vipSubnetId),
					resource.TestCheckResourceAttrSet("ovh_cloud_loadbalancer.test", "id"),
					resource.TestCheckResourceAttrSet("ovh_cloud_loadbalancer.test", "checksum"),
				),
			},
		},
	})
}

func testAccPreCheckCloudLoadbalancer(t *testing.T) {
	testAccPreCheckCredentials(t)
	if os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST") == "" {
		t.Skip("OVH_CLOUD_PROJECT_SERVICE_TEST not set")
	}
	if os.Getenv("OVH_CLOUD_PROJECT_LOADBALANCER_REGION_TEST") == "" {
		t.Skip("OVH_CLOUD_PROJECT_LOADBALANCER_REGION_TEST not set")
	}
	if os.Getenv("OVH_CLOUD_PROJECT_LOADBALANCER_VIP_NETWORK_ID_TEST") == "" {
		t.Skip("OVH_CLOUD_PROJECT_LOADBALANCER_VIP_NETWORK_ID_TEST not set")
	}
	if os.Getenv("OVH_CLOUD_PROJECT_LOADBALANCER_VIP_SUBNET_ID_TEST") == "" {
		t.Skip("OVH_CLOUD_PROJECT_LOADBALANCER_VIP_SUBNET_ID_TEST not set")
	}
}

func testAccCloudLoadbalancerImportStateIdFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}
		return fmt.Sprintf("%s/%s", rs.Primary.Attributes["service_name"], rs.Primary.Attributes["id"]), nil
	}
}
