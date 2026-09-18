package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

const testAccResourceCloudLoadbalancerListenerNamePrefix = "tf-test-listener-v2-"

func TestAccCloudLoadbalancerListener_basic(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	loadbalancerId := os.Getenv("OVH_CLOUD_LOADBALANCER_ID_TEST")

	listenerName := acctest.RandomWithPrefix(testAccResourceCloudLoadbalancerListenerNamePrefix)

	config := fmt.Sprintf(`
resource "ovh_cloud_loadbalancer_listener" "test" {
  service_name    = "%s"
  loadbalancer_id = "%s"
  name            = "%s"
  protocol        = "HTTP"
  protocol_port   = 80
}
`, serviceName, loadbalancerId, listenerName)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCloudLoadbalancerListener(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_loadbalancer_listener.test", "service_name", serviceName),
					resource.TestCheckResourceAttr("ovh_cloud_loadbalancer_listener.test", "loadbalancer_id", loadbalancerId),
					resource.TestCheckResourceAttr("ovh_cloud_loadbalancer_listener.test", "name", listenerName),
					resource.TestCheckResourceAttr("ovh_cloud_loadbalancer_listener.test", "protocol", "HTTP"),
					resource.TestCheckResourceAttr("ovh_cloud_loadbalancer_listener.test", "protocol_port", "80"),
					resource.TestCheckResourceAttrSet("ovh_cloud_loadbalancer_listener.test", "id"),
					resource.TestCheckResourceAttrSet("ovh_cloud_loadbalancer_listener.test", "checksum"),
					resource.TestCheckResourceAttrSet("ovh_cloud_loadbalancer_listener.test", "created_at"),
					resource.TestCheckResourceAttr("ovh_cloud_loadbalancer_listener.test", "resource_status", "READY"),
					resource.TestCheckResourceAttrSet("ovh_cloud_loadbalancer_listener.test", "current_state.name"),
					resource.TestCheckResourceAttrSet("ovh_cloud_loadbalancer_listener.test", "current_state.protocol"),
					resource.TestCheckResourceAttrSet("ovh_cloud_loadbalancer_listener.test", "current_state.protocol_port"),
				),
			},
			// Test import
			{
				ResourceName:      "ovh_cloud_loadbalancer_listener.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccCloudLoadbalancerListenerImportStateIdFunc("ovh_cloud_loadbalancer_listener.test"),
			},
		},
	})
}

func TestAccCloudLoadbalancerListener_update(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	loadbalancerId := os.Getenv("OVH_CLOUD_LOADBALANCER_ID_TEST")

	listenerName := acctest.RandomWithPrefix(testAccResourceCloudLoadbalancerListenerNamePrefix)
	updatedName := acctest.RandomWithPrefix(testAccResourceCloudLoadbalancerListenerNamePrefix)

	config := fmt.Sprintf(`
resource "ovh_cloud_loadbalancer_listener" "test" {
  service_name    = "%s"
  loadbalancer_id = "%s"
  name            = "%s"
  protocol        = "HTTP"
  protocol_port   = 80
  description     = "initial description"
}
`, serviceName, loadbalancerId, listenerName)

	updatedConfig := fmt.Sprintf(`
resource "ovh_cloud_loadbalancer_listener" "test" {
  service_name    = "%s"
  loadbalancer_id = "%s"
  name            = "%s"
  protocol        = "HTTP"
  protocol_port   = 80
  description     = "updated description"
}
`, serviceName, loadbalancerId, updatedName)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCloudLoadbalancerListener(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_loadbalancer_listener.test", "name", listenerName),
					resource.TestCheckResourceAttr("ovh_cloud_loadbalancer_listener.test", "description", "initial description"),
					resource.TestCheckResourceAttr("ovh_cloud_loadbalancer_listener.test", "protocol", "HTTP"),
					resource.TestCheckResourceAttr("ovh_cloud_loadbalancer_listener.test", "protocol_port", "80"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_loadbalancer_listener.test", "name", updatedName),
					resource.TestCheckResourceAttr("ovh_cloud_loadbalancer_listener.test", "description", "updated description"),
					resource.TestCheckResourceAttrSet("ovh_cloud_loadbalancer_listener.test", "id"),
					resource.TestCheckResourceAttrSet("ovh_cloud_loadbalancer_listener.test", "checksum"),
				),
			},
		},
	})
}

func TestAccCloudLoadbalancerListener_withInsertHeaders(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	loadbalancerId := os.Getenv("OVH_CLOUD_LOADBALANCER_ID_TEST")

	listenerName := acctest.RandomWithPrefix(testAccResourceCloudLoadbalancerListenerNamePrefix)

	config := fmt.Sprintf(`
resource "ovh_cloud_loadbalancer_listener" "test" {
  service_name    = "%s"
  loadbalancer_id = "%s"
  name            = "%s"
  protocol        = "HTTP"
  protocol_port   = 8080

  insert_headers {
    x_forwarded_for   = true
    x_forwarded_port  = true
    x_forwarded_proto = true
  }
}
`, serviceName, loadbalancerId, listenerName)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCloudLoadbalancerListener(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_loadbalancer_listener.test", "name", listenerName),
					resource.TestCheckResourceAttr("ovh_cloud_loadbalancer_listener.test", "insert_headers.x_forwarded_for", "true"),
					resource.TestCheckResourceAttr("ovh_cloud_loadbalancer_listener.test", "insert_headers.x_forwarded_port", "true"),
					resource.TestCheckResourceAttr("ovh_cloud_loadbalancer_listener.test", "insert_headers.x_forwarded_proto", "true"),
					resource.TestCheckResourceAttr("ovh_cloud_loadbalancer_listener.test", "resource_status", "READY"),
				),
			},
		},
	})
}

func testAccPreCheckCloudLoadbalancerListener(t *testing.T) {
	testAccPreCheckCredentials(t)
	if os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST") == "" {
		t.Skip("OVH_CLOUD_PROJECT_SERVICE_TEST not set")
	}
	if os.Getenv("OVH_CLOUD_LOADBALANCER_ID_TEST") == "" {
		t.Skip("OVH_CLOUD_LOADBALANCER_ID_TEST not set")
	}
}

func testAccCloudLoadbalancerListenerImportStateIdFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}
		return fmt.Sprintf("%s/%s/%s", rs.Primary.Attributes["service_name"], rs.Primary.Attributes["loadbalancer_id"], rs.Primary.Attributes["id"]), nil
	}
}
