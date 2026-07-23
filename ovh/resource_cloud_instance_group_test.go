package ovh

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func testAccPreCheckCloudInstanceGroup(t *testing.T) {
	testAccPreCheckCredentials(t)
	for _, v := range []string{
		"OVH_CLOUD_PROJECT_SERVICE_TEST",
		"OVH_CLOUD_PROJECT_REGION_TEST",
	} {
		if os.Getenv(v) == "" {
			t.Skipf("%s must be set for ovh_cloud_instance_group acceptance tests", v)
		}
	}
}

func testAccCloudInstanceGroupConfig(serviceName, region, name string) string {
	return fmt.Sprintf(`
resource "ovh_cloud_instance_group" "test" {
  service_name = "%s"
  region       = "%s"
  name         = "%s"
  policy       = "ANTI_AFFINITY"
}
`, serviceName, region, name)
}

func testAccCloudInstanceGroupImportStateIdFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("resource not found: %s", resourceName)
		}
		return rs.Primary.Attributes["service_name"] + "/" + rs.Primary.Attributes["id"], nil
	}
}

func TestAccCloudInstanceGroup_basic(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	region := os.Getenv("OVH_CLOUD_PROJECT_REGION_TEST")
	name := acctest.RandomWithPrefix("test-grp")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckCloudInstanceGroup(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudInstanceGroupConfig(serviceName, region, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_instance_group.test", "name", name),
					resource.TestCheckResourceAttr("ovh_cloud_instance_group.test", "region", region),
					resource.TestCheckResourceAttr("ovh_cloud_instance_group.test", "policy", "ANTI_AFFINITY"),
					resource.TestCheckResourceAttr("ovh_cloud_instance_group.test", "resource_status", "READY"),
					resource.TestCheckResourceAttrSet("ovh_cloud_instance_group.test", "id"),
					resource.TestCheckResourceAttrSet("ovh_cloud_instance_group.test", "checksum"),
				),
			},
			{
				ResourceName:      "ovh_cloud_instance_group.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccCloudInstanceGroupImportStateIdFunc("ovh_cloud_instance_group.test"),
				// checksum is refreshed on read; ignore volatile fields if import verify complains.
				ImportStateVerifyIgnore: []string{"checksum"},
			},
		},
	})
}

// testAccCloudInstanceGroupConfigPolicy renders an instance group config with an
// explicit placement policy.
func testAccCloudInstanceGroupConfigPolicy(serviceName, region, name, policy string) string {
	return fmt.Sprintf(`
resource "ovh_cloud_instance_group" "test" {
  service_name = "%s"
  region       = "%s"
  name         = "%s"
  policy       = "%s"
}
`, serviceName, region, name, policy)
}

// TestAccCloudInstanceGroup_immutable asserts that instance groups are immutable:
// changing the name (a RequiresReplace attribute) forces a replace.
func TestAccCloudInstanceGroup_immutable(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	region := os.Getenv("OVH_CLOUD_PROJECT_REGION_TEST")
	name := acctest.RandomWithPrefix("test-grp-imm")
	nameUpdated := acctest.RandomWithPrefix("test-grp-imm")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckCloudInstanceGroup(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudInstanceGroupConfig(serviceName, region, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_instance_group.test", "name", name),
					resource.TestCheckResourceAttr("ovh_cloud_instance_group.test", "resource_status", "READY"),
				),
			},
			{
				// Changing the (immutable) name must replace the group.
				Config: testAccCloudInstanceGroupConfig(serviceName, region, nameUpdated),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("ovh_cloud_instance_group.test", plancheck.ResourceActionReplace),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_instance_group.test", "name", nameUpdated),
				),
			},
		},
	})
}

// TestAccCloudInstanceGroup_validators asserts the policy OneOf validator rejects
// an invalid placement policy at plan time.
func TestAccCloudInstanceGroup_validators(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	region := os.Getenv("OVH_CLOUD_PROJECT_REGION_TEST")
	name := acctest.RandomWithPrefix("test-grp-val")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckCloudInstanceGroup(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccCloudInstanceGroupConfigPolicy(serviceName, region, name, "WRONG"),
				ExpectError: regexp.MustCompile(`must be one of`),
			},
		},
	})
}
