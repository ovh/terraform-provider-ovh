package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCloudKeyManagerContainerDataSource_basic(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	region := os.Getenv("OVH_CLOUD_PROJECT_REGION_TEST")
	name := acctest.RandomWithPrefix(test_prefix)

	config := fmt.Sprintf(`
resource "ovh_cloud_key_manager_container" "test" {
  service_name = "%s"
  region       = "%s"
  name         = "%s"
  type         = "GENERIC"
}

data "ovh_cloud_key_manager_container" "test" {
  service_name = "%s"
  container_id = ovh_cloud_key_manager_container.test.id
}
`, serviceName, region, name, serviceName)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCloudKeyManager(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.ovh_cloud_key_manager_container.test", "service_name", serviceName),
					resource.TestCheckResourceAttr("data.ovh_cloud_key_manager_container.test", "location.region", region),
					resource.TestCheckResourceAttr("data.ovh_cloud_key_manager_container.test", "name", name),
					resource.TestCheckResourceAttr("data.ovh_cloud_key_manager_container.test", "type", "GENERIC"),
					resource.TestCheckResourceAttrSet("data.ovh_cloud_key_manager_container.test", "id"),
					resource.TestCheckResourceAttrSet("data.ovh_cloud_key_manager_container.test", "checksum"),
					resource.TestCheckResourceAttrSet("data.ovh_cloud_key_manager_container.test", "created_at"),
					resource.TestCheckResourceAttr("data.ovh_cloud_key_manager_container.test", "resource_status", "READY"),
				),
			},
		},
	})
}
