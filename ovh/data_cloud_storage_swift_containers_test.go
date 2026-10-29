package ovh

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceCloudStorageSwiftContainers_basic(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	region := strings.ToUpper(os.Getenv("OVH_CLOUD_PROJECT_STORAGE_REGION_TEST"))

	containerName := acctest.RandomWithPrefix(test_prefix)

	config := fmt.Sprintf(`
resource "ovh_cloud_storage_swift_container" "container" {
  service_name = "%s"
  name         = "%s"
  region       = "%s"
}

data "ovh_cloud_storage_swift_containers" "containers" {
  service_name = "%s"
  region       = "%s"

  depends_on = [ovh_cloud_storage_swift_container.container]
}
`, serviceName, containerName, region, serviceName, region)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCloudStorageSwiftContainer(t)
			testAccCheckCloudProjectExists(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudStorageSwiftContainerDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.ovh_cloud_storage_swift_containers.containers", "service_name", serviceName),
					resource.TestCheckResourceAttr("data.ovh_cloud_storage_swift_containers.containers", "region", region),
					resource.TestCheckResourceAttrSet("data.ovh_cloud_storage_swift_containers.containers", "containers.0.id"),
					resource.TestCheckResourceAttrSet("data.ovh_cloud_storage_swift_containers.containers", "containers.0.name"),
					resource.TestCheckResourceAttrSet("data.ovh_cloud_storage_swift_containers.containers", "containers.0.checksum"),
					resource.TestCheckResourceAttrSet("data.ovh_cloud_storage_swift_containers.containers", "containers.0.resource_status"),
					resource.TestCheckResourceAttr("data.ovh_cloud_storage_swift_containers.containers", "containers.0.location.region", region),
					resource.TestCheckResourceAttr("data.ovh_cloud_storage_swift_containers.containers", "containers.0.current_state.location.region", region),
					resource.TestCheckTypeSetElemNestedAttrs("data.ovh_cloud_storage_swift_containers.containers", "containers.*", map[string]string{
						"id":   region + "_" + containerName,
						"name": containerName,
						"type": "PRIVATE",
					}),
				),
			},
		},
	})
}
