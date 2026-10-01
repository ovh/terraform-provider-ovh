package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCloudKeyManagerContainersDataSource_basic(t *testing.T) {
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

data "ovh_cloud_key_manager_containers" "all" {
  service_name = "%s"

  depends_on = [ovh_cloud_key_manager_container.test]
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
					resource.TestCheckResourceAttr("data.ovh_cloud_key_manager_containers.all", "service_name", serviceName),
					resource.TestCheckResourceAttrSet("data.ovh_cloud_key_manager_containers.all", "containers.#"),
				),
			},
		},
	})
}
