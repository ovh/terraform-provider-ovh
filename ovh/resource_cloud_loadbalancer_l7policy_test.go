package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

const testAccResourceCloudLoadbalancerL7PolicyNamePrefix = "tf-test-l7pol-v2-"

func TestAccCloudLoadbalancerL7Policy_basic(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	loadbalancerId := os.Getenv("OVH_CLOUD_LOADBALANCER_ID_TEST")
	listenerId := os.Getenv("OVH_CLOUD_LOADBALANCER_LISTENER_ID_TEST")

	policyName := acctest.RandomWithPrefix(testAccResourceCloudLoadbalancerL7PolicyNamePrefix)

	config := fmt.Sprintf(`
resource "ovh_cloud_loadbalancer_l7policy" "test" {
  service_name    = "%s"
  loadbalancer_id = "%s"
  listener_id     = "%s"
  name            = "%s"
  action          = "REJECT"
}
`, serviceName, loadbalancerId, listenerId, policyName)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCloudLoadbalancerL7Policy(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_loadbalancer_l7policy.test", "service_name", serviceName),
					resource.TestCheckResourceAttr("ovh_cloud_loadbalancer_l7policy.test", "loadbalancer_id", loadbalancerId),
					resource.TestCheckResourceAttr("ovh_cloud_loadbalancer_l7policy.test", "listener_id", listenerId),
					resource.TestCheckResourceAttr("ovh_cloud_loadbalancer_l7policy.test", "name", policyName),
					resource.TestCheckResourceAttr("ovh_cloud_loadbalancer_l7policy.test", "action", "REJECT"),
					resource.TestCheckResourceAttrSet("ovh_cloud_loadbalancer_l7policy.test", "id"),
					resource.TestCheckResourceAttrSet("ovh_cloud_loadbalancer_l7policy.test", "checksum"),
					resource.TestCheckResourceAttrSet("ovh_cloud_loadbalancer_l7policy.test", "created_at"),
					resource.TestCheckResourceAttr("ovh_cloud_loadbalancer_l7policy.test", "resource_status", "READY"),
					resource.TestCheckResourceAttrSet("ovh_cloud_loadbalancer_l7policy.test", "current_state.action"),
				),
			},
			// Test import
			{
				ResourceName:      "ovh_cloud_loadbalancer_l7policy.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccCloudLoadbalancerL7PolicyImportStateIdFunc("ovh_cloud_loadbalancer_l7policy.test"),
			},
		},
	})
}

func TestAccCloudLoadbalancerL7Policy_update(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	loadbalancerId := os.Getenv("OVH_CLOUD_LOADBALANCER_ID_TEST")
	listenerId := os.Getenv("OVH_CLOUD_LOADBALANCER_LISTENER_ID_TEST")

	policyName := acctest.RandomWithPrefix(testAccResourceCloudLoadbalancerL7PolicyNamePrefix)
	updatedName := acctest.RandomWithPrefix(testAccResourceCloudLoadbalancerL7PolicyNamePrefix)

	config := fmt.Sprintf(`
resource "ovh_cloud_loadbalancer_l7policy" "test" {
  service_name    = "%s"
  loadbalancer_id = "%s"
  listener_id     = "%s"
  name            = "%s"
  description     = "initial description"
  action          = "REJECT"
}
`, serviceName, loadbalancerId, listenerId, policyName)

	updatedConfig := fmt.Sprintf(`
resource "ovh_cloud_loadbalancer_l7policy" "test" {
  service_name       = "%s"
  loadbalancer_id    = "%s"
  listener_id        = "%s"
  name               = "%s"
  description        = "updated description"
  action             = "REDIRECT_TO_URL"
  redirect_url       = "https://example.com"
  redirect_http_code = 302
}
`, serviceName, loadbalancerId, listenerId, updatedName)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCloudLoadbalancerL7Policy(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_loadbalancer_l7policy.test", "name", policyName),
					resource.TestCheckResourceAttr("ovh_cloud_loadbalancer_l7policy.test", "description", "initial description"),
					resource.TestCheckResourceAttr("ovh_cloud_loadbalancer_l7policy.test", "action", "REJECT"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_loadbalancer_l7policy.test", "name", updatedName),
					resource.TestCheckResourceAttr("ovh_cloud_loadbalancer_l7policy.test", "description", "updated description"),
					resource.TestCheckResourceAttr("ovh_cloud_loadbalancer_l7policy.test", "action", "REDIRECT_TO_URL"),
					resource.TestCheckResourceAttr("ovh_cloud_loadbalancer_l7policy.test", "redirect_url", "https://example.com"),
					resource.TestCheckResourceAttr("ovh_cloud_loadbalancer_l7policy.test", "redirect_http_code", "302"),
					resource.TestCheckResourceAttrSet("ovh_cloud_loadbalancer_l7policy.test", "id"),
					resource.TestCheckResourceAttrSet("ovh_cloud_loadbalancer_l7policy.test", "checksum"),
				),
			},
		},
	})
}

func TestAccCloudLoadbalancerL7Policy_withRules(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	loadbalancerId := os.Getenv("OVH_CLOUD_LOADBALANCER_ID_TEST")
	listenerId := os.Getenv("OVH_CLOUD_LOADBALANCER_LISTENER_ID_TEST")

	policyName := acctest.RandomWithPrefix(testAccResourceCloudLoadbalancerL7PolicyNamePrefix)

	config := fmt.Sprintf(`
resource "ovh_cloud_loadbalancer_l7policy" "test" {
  service_name    = "%s"
  loadbalancer_id = "%s"
  listener_id     = "%s"
  name            = "%s"
  action          = "REJECT"

  rules {
    type         = "PATH"
    compare_type = "STARTS_WITH"
    value        = "/admin"
  }

  rules {
    type         = "HEADER"
    compare_type = "EQUAL_TO"
    key          = "X-Custom-Header"
    value        = "blocked"
    invert       = true
  }
}
`, serviceName, loadbalancerId, listenerId, policyName)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCloudLoadbalancerL7Policy(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_loadbalancer_l7policy.test", "name", policyName),
					resource.TestCheckResourceAttr("ovh_cloud_loadbalancer_l7policy.test", "action", "REJECT"),
					resource.TestCheckResourceAttr("ovh_cloud_loadbalancer_l7policy.test", "rules.#", "2"),
					resource.TestCheckResourceAttr("ovh_cloud_loadbalancer_l7policy.test", "resource_status", "READY"),
				),
			},
		},
	})
}

func TestAccCloudLoadbalancerL7Policy_redirectToPool(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	loadbalancerId := os.Getenv("OVH_CLOUD_LOADBALANCER_ID_TEST")
	listenerId := os.Getenv("OVH_CLOUD_LOADBALANCER_LISTENER_ID_TEST")

	policyName := acctest.RandomWithPrefix(testAccResourceCloudLoadbalancerL7PolicyNamePrefix)
	poolName := acctest.RandomWithPrefix("tf-test-pool-v2-")

	config := fmt.Sprintf(`
resource "ovh_cloud_loadbalancer_pool" "redirect_target" {
  service_name    = "%s"
  loadbalancer_id = "%s"
  name            = "%s"
  protocol        = "HTTP"
  algorithm       = "ROUND_ROBIN"
}

resource "ovh_cloud_loadbalancer_l7policy" "test" {
  service_name    = "%s"
  loadbalancer_id = "%s"
  listener_id     = "%s"
  name            = "%s"
  action          = "REDIRECT_TO_POOL"
  redirect_pool_id = ovh_cloud_loadbalancer_pool.redirect_target.id

  rules {
    type         = "PATH"
    compare_type = "STARTS_WITH"
    value        = "/api"
  }
}
`, serviceName, loadbalancerId, poolName,
		serviceName, loadbalancerId, listenerId, policyName)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCloudLoadbalancerL7Policy(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_loadbalancer_l7policy.test", "name", policyName),
					resource.TestCheckResourceAttr("ovh_cloud_loadbalancer_l7policy.test", "action", "REDIRECT_TO_POOL"),
					resource.TestCheckResourceAttrSet("ovh_cloud_loadbalancer_l7policy.test", "redirect_pool_id"),
					resource.TestCheckResourceAttr("ovh_cloud_loadbalancer_l7policy.test", "resource_status", "READY"),
				),
			},
		},
	})
}

func testAccPreCheckCloudLoadbalancerL7Policy(t *testing.T) {
	testAccPreCheckCredentials(t)
	if os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST") == "" {
		t.Skip("OVH_CLOUD_PROJECT_SERVICE_TEST not set")
	}
	if os.Getenv("OVH_CLOUD_LOADBALANCER_ID_TEST") == "" {
		t.Skip("OVH_CLOUD_LOADBALANCER_ID_TEST not set")
	}
	if os.Getenv("OVH_CLOUD_LOADBALANCER_LISTENER_ID_TEST") == "" {
		t.Skip("OVH_CLOUD_LOADBALANCER_LISTENER_ID_TEST not set")
	}
}

func testAccCloudLoadbalancerL7PolicyImportStateIdFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}
		return fmt.Sprintf("%s/%s/%s/%s",
			rs.Primary.Attributes["service_name"],
			rs.Primary.Attributes["loadbalancer_id"],
			rs.Primary.Attributes["listener_id"],
			rs.Primary.Attributes["id"],
		), nil
	}
}
