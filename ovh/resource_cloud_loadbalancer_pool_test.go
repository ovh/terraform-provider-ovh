package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

const testAccResourceCloudLoadbalancerPoolNamePrefix = "tf-test-pool-v2-"

func TestAccCloudLoadbalancerPool_basic(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	loadbalancerId := os.Getenv("OVH_CLOUD_LOADBALANCER_ID_TEST")

	poolName := acctest.RandomWithPrefix(testAccResourceCloudLoadbalancerPoolNamePrefix)

	config := fmt.Sprintf(`
resource "ovh_cloud_loadbalancer_pool" "test" {
  service_name    = "%s"
  loadbalancer_id = "%s"
  name            = "%s"
  protocol        = "HTTP"
  algorithm       = "ROUND_ROBIN"
}
`, serviceName, loadbalancerId, poolName)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCloudLoadbalancerPool(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_loadbalancer_pool.test", "service_name", serviceName),
					resource.TestCheckResourceAttr("ovh_cloud_loadbalancer_pool.test", "loadbalancer_id", loadbalancerId),
					resource.TestCheckResourceAttr("ovh_cloud_loadbalancer_pool.test", "name", poolName),
					resource.TestCheckResourceAttr("ovh_cloud_loadbalancer_pool.test", "protocol", "HTTP"),
					resource.TestCheckResourceAttr("ovh_cloud_loadbalancer_pool.test", "algorithm", "ROUND_ROBIN"),
					resource.TestCheckResourceAttrSet("ovh_cloud_loadbalancer_pool.test", "id"),
					resource.TestCheckResourceAttrSet("ovh_cloud_loadbalancer_pool.test", "checksum"),
					resource.TestCheckResourceAttrSet("ovh_cloud_loadbalancer_pool.test", "created_at"),
					resource.TestCheckResourceAttr("ovh_cloud_loadbalancer_pool.test", "resource_status", "READY"),
					resource.TestCheckResourceAttrSet("ovh_cloud_loadbalancer_pool.test", "current_state.protocol"),
					resource.TestCheckResourceAttrSet("ovh_cloud_loadbalancer_pool.test", "current_state.algorithm"),
				),
			},
			// Test import
			{
				ResourceName:      "ovh_cloud_loadbalancer_pool.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccCloudLoadbalancerPoolImportStateIdFunc("ovh_cloud_loadbalancer_pool.test"),
			},
		},
	})
}

func TestAccCloudLoadbalancerPool_update(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	loadbalancerId := os.Getenv("OVH_CLOUD_LOADBALANCER_ID_TEST")

	poolName := acctest.RandomWithPrefix(testAccResourceCloudLoadbalancerPoolNamePrefix)
	updatedName := acctest.RandomWithPrefix(testAccResourceCloudLoadbalancerPoolNamePrefix)

	config := fmt.Sprintf(`
resource "ovh_cloud_loadbalancer_pool" "test" {
  service_name    = "%s"
  loadbalancer_id = "%s"
  name            = "%s"
  description     = "initial description"
  protocol        = "HTTP"
  algorithm       = "ROUND_ROBIN"
}
`, serviceName, loadbalancerId, poolName)

	updatedConfig := fmt.Sprintf(`
resource "ovh_cloud_loadbalancer_pool" "test" {
  service_name    = "%s"
  loadbalancer_id = "%s"
  name            = "%s"
  description     = "updated description"
  protocol        = "HTTP"
  algorithm       = "LEAST_CONNECTIONS"
}
`, serviceName, loadbalancerId, updatedName)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCloudLoadbalancerPool(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_loadbalancer_pool.test", "name", poolName),
					resource.TestCheckResourceAttr("ovh_cloud_loadbalancer_pool.test", "description", "initial description"),
					resource.TestCheckResourceAttr("ovh_cloud_loadbalancer_pool.test", "algorithm", "ROUND_ROBIN"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_loadbalancer_pool.test", "name", updatedName),
					resource.TestCheckResourceAttr("ovh_cloud_loadbalancer_pool.test", "description", "updated description"),
					resource.TestCheckResourceAttr("ovh_cloud_loadbalancer_pool.test", "algorithm", "LEAST_CONNECTIONS"),
					resource.TestCheckResourceAttrSet("ovh_cloud_loadbalancer_pool.test", "id"),
					resource.TestCheckResourceAttrSet("ovh_cloud_loadbalancer_pool.test", "checksum"),
				),
			},
		},
	})
}

func TestAccCloudLoadbalancerPool_withPersistence(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	loadbalancerId := os.Getenv("OVH_CLOUD_LOADBALANCER_ID_TEST")

	poolName := acctest.RandomWithPrefix(testAccResourceCloudLoadbalancerPoolNamePrefix)

	config := fmt.Sprintf(`
resource "ovh_cloud_loadbalancer_pool" "test" {
  service_name    = "%s"
  loadbalancer_id = "%s"
  name            = "%s"
  protocol        = "HTTP"
  algorithm       = "ROUND_ROBIN"

  persistence {
    type        = "APP_COOKIE"
    cookie_name = "JSESSIONID"
  }
}
`, serviceName, loadbalancerId, poolName)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCloudLoadbalancerPool(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_loadbalancer_pool.test", "name", poolName),
					resource.TestCheckResourceAttr("ovh_cloud_loadbalancer_pool.test", "protocol", "HTTP"),
					resource.TestCheckResourceAttr("ovh_cloud_loadbalancer_pool.test", "algorithm", "ROUND_ROBIN"),
					resource.TestCheckResourceAttr("ovh_cloud_loadbalancer_pool.test", "persistence.type", "APP_COOKIE"),
					resource.TestCheckResourceAttr("ovh_cloud_loadbalancer_pool.test", "persistence.cookie_name", "JSESSIONID"),
					resource.TestCheckResourceAttr("ovh_cloud_loadbalancer_pool.test", "resource_status", "READY"),
				),
			},
		},
	})
}

func testAccPreCheckCloudLoadbalancerPool(t *testing.T) {
	testAccPreCheckCredentials(t)
	if os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST") == "" {
		t.Skip("OVH_CLOUD_PROJECT_SERVICE_TEST not set")
	}
	if os.Getenv("OVH_CLOUD_LOADBALANCER_ID_TEST") == "" {
		t.Skip("OVH_CLOUD_LOADBALANCER_ID_TEST not set")
	}
}

func testAccCloudLoadbalancerPoolImportStateIdFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}
		return fmt.Sprintf("%s/%s/%s",
			rs.Primary.Attributes["service_name"],
			rs.Primary.Attributes["loadbalancer_id"],
			rs.Primary.Attributes["id"],
		), nil
	}
}
