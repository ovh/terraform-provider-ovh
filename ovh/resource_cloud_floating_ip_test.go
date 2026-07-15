package ovh

import (
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

const testAccResourceCloudFloatingIPDescriptionPrefix = "tf-test-public-ip-floating-"

func TestAccCloudFloatingIP_basic(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	region := os.Getenv("OVH_CLOUD_PROJECT_REGION_TEST")

	description := acctest.RandomWithPrefix(testAccResourceCloudFloatingIPDescriptionPrefix)

	config := fmt.Sprintf(`
resource "ovh_cloud_floating_ip" "test" {
  service_name = "%s"
  region       = "%s"
  description  = "%s"
}
`, serviceName, region, description)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCloudPublicIP(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_floating_ip.test", "service_name", serviceName),
					resource.TestCheckResourceAttr("ovh_cloud_floating_ip.test", "region", region),
					resource.TestCheckResourceAttr("ovh_cloud_floating_ip.test", "description", description),
					resource.TestCheckResourceAttr("ovh_cloud_floating_ip.test", "resource_status", "READY"),
					resource.TestCheckResourceAttrSet("ovh_cloud_floating_ip.test", "id"),
					resource.TestCheckResourceAttrSet("ovh_cloud_floating_ip.test", "checksum"),
					resource.TestCheckResourceAttrSet("ovh_cloud_floating_ip.test", "created_at"),
					resource.TestCheckResourceAttrSet("ovh_cloud_floating_ip.test", "current_state.ip"),
					resource.TestCheckResourceAttrSet("ovh_cloud_floating_ip.test", "current_state.status"),
					resource.TestCheckResourceAttrSet("ovh_cloud_floating_ip.test", "current_state.id"),
					resource.TestCheckResourceAttr("ovh_cloud_floating_ip.test", "current_state.location.region", region),
				),
			},
			// Test import
			{
				ResourceName:      "ovh_cloud_floating_ip.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccCloudFloatingIPImportStateIdFunc("ovh_cloud_floating_ip.test"),
			},
		},
	})
}

func TestAccCloudFloatingIP_updateDescription(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	region := os.Getenv("OVH_CLOUD_PROJECT_REGION_TEST")

	configTemplate := `
resource "ovh_cloud_floating_ip" "test" {
  service_name = "%s"
  region       = "%s"
  description  = "%s"
}
`

	config := fmt.Sprintf(configTemplate, serviceName, region, "initial")
	updatedConfig := fmt.Sprintf(configTemplate, serviceName, region, "updated")

	// Captured in step 1 and compared in step 2 to prove the floating IP was
	// updated in place, not replaced.
	var floatingIPID string

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCloudPublicIP(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_floating_ip.test", "description", "initial"),
					resource.TestCheckResourceAttr("ovh_cloud_floating_ip.test", "resource_status", "READY"),
					resource.TestCheckResourceAttrWith("ovh_cloud_floating_ip.test", "id", func(value string) error {
						if value == "" {
							return fmt.Errorf("expected floating IP id to be set")
						}
						floatingIPID = value
						return nil
					}),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_floating_ip.test", "description", "updated"),
					resource.TestCheckResourceAttr("ovh_cloud_floating_ip.test", "resource_status", "READY"),
					resource.TestCheckResourceAttrSet("ovh_cloud_floating_ip.test", "checksum"),
					resource.TestCheckResourceAttrWith("ovh_cloud_floating_ip.test", "id", func(value string) error {
						if value != floatingIPID {
							return fmt.Errorf("floating IP was replaced during update: id changed from %q to %q", floatingIPID, value)
						}
						return nil
					}),
				),
			},
		},
	})
}

// TestAccCloudFloatingIP_serviceNameFromEnv validates that when service_name is
// omitted from the resource configuration, the provider falls back to the
// OVH_CLOUD_PROJECT_SERVICE environment variable at plan time (via the
// EnvDefaultString plan modifier) and that this does not produce a perpetual
// diff / phantom replace on subsequent plans.
func TestAccCloudFloatingIP_serviceNameFromEnv(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	region := os.Getenv("OVH_CLOUD_PROJECT_REGION_TEST")

	// Make the project id resolvable from the environment so the resource can be
	// configured without an explicit service_name.
	t.Setenv("OVH_CLOUD_PROJECT_SERVICE", serviceName)

	description := acctest.RandomWithPrefix(testAccResourceCloudFloatingIPDescriptionPrefix)

	// service_name is intentionally omitted from the config: it must be resolved
	// from OVH_CLOUD_PROJECT_SERVICE by the EnvDefaultString plan modifier.
	config := fmt.Sprintf(`
resource "ovh_cloud_floating_ip" "test" {
  region      = "%s"
  description = "%s"
}
`, region, description)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCloudPublicIP(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					// The env default was applied even though service_name was
					// absent from the config.
					resource.TestCheckResourceAttr("ovh_cloud_floating_ip.test", "service_name", serviceName),
					resource.TestCheckResourceAttr("ovh_cloud_floating_ip.test", "region", region),
					resource.TestCheckResourceAttr("ovh_cloud_floating_ip.test", "description", description),
					resource.TestCheckResourceAttr("ovh_cloud_floating_ip.test", "resource_status", "READY"),
					resource.TestCheckResourceAttrSet("ovh_cloud_floating_ip.test", "id"),
				),
			},
			{
				// Re-planning with service_name still omitted must be a no-op:
				// the EnvDefaultString modifier injects the env value so the plan
				// matches state and RequiresReplace does not fire (regression guard).
				Config: config,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
		},
	})
}

// testAccPreCheckCloudPublicIP is the shared PreCheck for the public IP family
// (floating, extNet, additional, aggregate) acceptance tests.
func testAccPreCheckCloudPublicIP(t *testing.T) {
	testAccPreCheckCredentials(t)
	if os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST") == "" {
		t.Skip("OVH_CLOUD_PROJECT_SERVICE_TEST not set")
	}
	if os.Getenv("OVH_CLOUD_PROJECT_REGION_TEST") == "" {
		t.Skip("OVH_CLOUD_PROJECT_REGION_TEST not set")
	}
}

func testAccCloudFloatingIPImportStateIdFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}
		return fmt.Sprintf("%s/%s", rs.Primary.Attributes["service_name"], rs.Primary.Attributes["id"]), nil
	}
}

// testAccCheckCloudPublicIPListNotEmpty checks that a list count attribute
// (e.g. "floating_ips.#") parses to at least one element.
func testAccCheckCloudPublicIPListNotEmpty(value string) error {
	count, err := strconv.Atoi(value)
	if err != nil {
		return fmt.Errorf("failed to parse list length %q: %w", value, err)
	}
	if count < 1 {
		return fmt.Errorf("expected at least one element in the list, got %d", count)
	}
	return nil
}
