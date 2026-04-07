package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

const testAccResourceCloudLoadbalancerPoolMemberNamePrefix = "tf-test-member-v2-"

func TestAccCloudLoadbalancerPoolMember_basic(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	loadbalancerId := os.Getenv("OVH_CLOUD_LOADBALANCER_ID_TEST")
	poolId := os.Getenv("OVH_CLOUD_LOADBALANCER_POOL_ID_TEST")
	memberAddress := os.Getenv("OVH_CLOUD_LOADBALANCER_MEMBER_ADDRESS_TEST")

	memberName := acctest.RandomWithPrefix(testAccResourceCloudLoadbalancerPoolMemberNamePrefix)

	config := fmt.Sprintf(`
resource "ovh_cloud_loadbalancer_pool_member" "test" {
  service_name    = "%s"
  loadbalancer_id = "%s"
  pool_id         = "%s"
  name            = "%s"
  address         = "%s"
  protocol_port   = 8080
}
`, serviceName, loadbalancerId, poolId, memberName, memberAddress)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCloudLoadbalancerPoolMember(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_loadbalancer_pool_member.test", "service_name", serviceName),
					resource.TestCheckResourceAttr("ovh_cloud_loadbalancer_pool_member.test", "loadbalancer_id", loadbalancerId),
					resource.TestCheckResourceAttr("ovh_cloud_loadbalancer_pool_member.test", "pool_id", poolId),
					resource.TestCheckResourceAttr("ovh_cloud_loadbalancer_pool_member.test", "name", memberName),
					resource.TestCheckResourceAttr("ovh_cloud_loadbalancer_pool_member.test", "address", memberAddress),
					resource.TestCheckResourceAttr("ovh_cloud_loadbalancer_pool_member.test", "protocol_port", "8080"),
					resource.TestCheckResourceAttrSet("ovh_cloud_loadbalancer_pool_member.test", "id"),
					resource.TestCheckResourceAttrSet("ovh_cloud_loadbalancer_pool_member.test", "checksum"),
					resource.TestCheckResourceAttrSet("ovh_cloud_loadbalancer_pool_member.test", "created_at"),
					resource.TestCheckResourceAttr("ovh_cloud_loadbalancer_pool_member.test", "resource_status", "READY"),
					resource.TestCheckResourceAttrSet("ovh_cloud_loadbalancer_pool_member.test", "current_state.address"),
					resource.TestCheckResourceAttrSet("ovh_cloud_loadbalancer_pool_member.test", "current_state.protocol_port"),
				),
			},
			// Test import
			{
				ResourceName:      "ovh_cloud_loadbalancer_pool_member.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccCloudLoadbalancerPoolMemberImportStateIdFunc("ovh_cloud_loadbalancer_pool_member.test"),
			},
		},
	})
}

func TestAccCloudLoadbalancerPoolMember_update(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	loadbalancerId := os.Getenv("OVH_CLOUD_LOADBALANCER_ID_TEST")
	poolId := os.Getenv("OVH_CLOUD_LOADBALANCER_POOL_ID_TEST")
	memberAddress := os.Getenv("OVH_CLOUD_LOADBALANCER_MEMBER_ADDRESS_TEST")

	memberName := acctest.RandomWithPrefix(testAccResourceCloudLoadbalancerPoolMemberNamePrefix)
	updatedName := acctest.RandomWithPrefix(testAccResourceCloudLoadbalancerPoolMemberNamePrefix)

	config := fmt.Sprintf(`
resource "ovh_cloud_loadbalancer_pool_member" "test" {
  service_name    = "%s"
  loadbalancer_id = "%s"
  pool_id         = "%s"
  name            = "%s"
  address         = "%s"
  protocol_port   = 8080
  weight          = 10
}
`, serviceName, loadbalancerId, poolId, memberName, memberAddress)

	updatedConfig := fmt.Sprintf(`
resource "ovh_cloud_loadbalancer_pool_member" "test" {
  service_name    = "%s"
  loadbalancer_id = "%s"
  pool_id         = "%s"
  name            = "%s"
  address         = "%s"
  protocol_port   = 8080
  weight          = 50
}
`, serviceName, loadbalancerId, poolId, updatedName, memberAddress)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCloudLoadbalancerPoolMember(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_loadbalancer_pool_member.test", "name", memberName),
					resource.TestCheckResourceAttr("ovh_cloud_loadbalancer_pool_member.test", "weight", "10"),
					resource.TestCheckResourceAttr("ovh_cloud_loadbalancer_pool_member.test", "resource_status", "READY"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_loadbalancer_pool_member.test", "name", updatedName),
					resource.TestCheckResourceAttr("ovh_cloud_loadbalancer_pool_member.test", "weight", "50"),
					resource.TestCheckResourceAttrSet("ovh_cloud_loadbalancer_pool_member.test", "id"),
					resource.TestCheckResourceAttrSet("ovh_cloud_loadbalancer_pool_member.test", "checksum"),
				),
			},
		},
	})
}

func TestAccCloudLoadbalancerPoolMember_withMonitor(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	loadbalancerId := os.Getenv("OVH_CLOUD_LOADBALANCER_ID_TEST")
	poolId := os.Getenv("OVH_CLOUD_LOADBALANCER_POOL_ID_TEST")
	memberAddress := os.Getenv("OVH_CLOUD_LOADBALANCER_MEMBER_ADDRESS_TEST")
	monitorAddress := os.Getenv("OVH_CLOUD_LOADBALANCER_MEMBER_MONITOR_ADDRESS_TEST")

	memberName := acctest.RandomWithPrefix(testAccResourceCloudLoadbalancerPoolMemberNamePrefix)

	config := fmt.Sprintf(`
resource "ovh_cloud_loadbalancer_pool_member" "test" {
  service_name    = "%s"
  loadbalancer_id = "%s"
  pool_id         = "%s"
  name            = "%s"
  address         = "%s"
  protocol_port   = 8080

  monitor {
    address = "%s"
    port    = 9090
  }
}
`, serviceName, loadbalancerId, poolId, memberName, memberAddress, monitorAddress)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCloudLoadbalancerPoolMemberWithMonitor(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_loadbalancer_pool_member.test", "name", memberName),
					resource.TestCheckResourceAttr("ovh_cloud_loadbalancer_pool_member.test", "address", memberAddress),
					resource.TestCheckResourceAttr("ovh_cloud_loadbalancer_pool_member.test", "protocol_port", "8080"),
					resource.TestCheckResourceAttr("ovh_cloud_loadbalancer_pool_member.test", "monitor.address", monitorAddress),
					resource.TestCheckResourceAttr("ovh_cloud_loadbalancer_pool_member.test", "monitor.port", "9090"),
					resource.TestCheckResourceAttr("ovh_cloud_loadbalancer_pool_member.test", "resource_status", "READY"),
				),
			},
		},
	})
}

func testAccPreCheckCloudLoadbalancerPoolMember(t *testing.T) {
	testAccPreCheckCredentials(t)
	if os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST") == "" {
		t.Skip("OVH_CLOUD_PROJECT_SERVICE_TEST not set")
	}
	if os.Getenv("OVH_CLOUD_LOADBALANCER_ID_TEST") == "" {
		t.Skip("OVH_CLOUD_LOADBALANCER_ID_TEST not set")
	}
	if os.Getenv("OVH_CLOUD_LOADBALANCER_POOL_ID_TEST") == "" {
		t.Skip("OVH_CLOUD_LOADBALANCER_POOL_ID_TEST not set")
	}
	if os.Getenv("OVH_CLOUD_LOADBALANCER_MEMBER_ADDRESS_TEST") == "" {
		t.Skip("OVH_CLOUD_LOADBALANCER_MEMBER_ADDRESS_TEST not set")
	}
}

func testAccPreCheckCloudLoadbalancerPoolMemberWithMonitor(t *testing.T) {
	testAccPreCheckCloudLoadbalancerPoolMember(t)
	if os.Getenv("OVH_CLOUD_LOADBALANCER_MEMBER_MONITOR_ADDRESS_TEST") == "" {
		t.Skip("OVH_CLOUD_LOADBALANCER_MEMBER_MONITOR_ADDRESS_TEST not set")
	}
}

func testAccCloudLoadbalancerPoolMemberImportStateIdFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}
		return fmt.Sprintf("%s/%s/%s/%s",
			rs.Primary.Attributes["service_name"],
			rs.Primary.Attributes["loadbalancer_id"],
			rs.Primary.Attributes["pool_id"],
			rs.Primary.Attributes["id"],
		), nil
	}
}
