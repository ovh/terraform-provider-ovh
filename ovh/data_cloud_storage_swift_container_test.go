package ovh

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceCloudStorageSwiftContainer_basic(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	region := strings.ToUpper(os.Getenv("OVH_CLOUD_PROJECT_STORAGE_REGION_TEST"))

	containerName := acctest.RandomWithPrefix(test_prefix)

	config := fmt.Sprintf(`
resource "ovh_cloud_storage_swift_container" "container" {
  service_name = "%s"
  name         = "%s"
  region       = "%s"
  type         = "PUBLIC"
  cors         = ["https://www.example.com"]
}

data "ovh_cloud_storage_swift_container" "by_id" {
  service_name = "%s"
  id           = ovh_cloud_storage_swift_container.container.id
}
`, serviceName, containerName, region, serviceName)

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
					resource.TestCheckResourceAttr("data.ovh_cloud_storage_swift_container.by_id", "service_name", serviceName),
					resource.TestCheckResourceAttrPair(
						"data.ovh_cloud_storage_swift_container.by_id", "id",
						"ovh_cloud_storage_swift_container.container", "id",
					),
					resource.TestCheckResourceAttr("data.ovh_cloud_storage_swift_container.by_id", "id", region+"_"+containerName),
					resource.TestCheckResourceAttr("data.ovh_cloud_storage_swift_container.by_id", "name", containerName),
					resource.TestCheckResourceAttr("data.ovh_cloud_storage_swift_container.by_id", "location.region", region),
					resource.TestCheckResourceAttr("data.ovh_cloud_storage_swift_container.by_id", "archive", "false"),
					resource.TestCheckResourceAttr("data.ovh_cloud_storage_swift_container.by_id", "type", "PUBLIC"),
					resource.TestCheckResourceAttr("data.ovh_cloud_storage_swift_container.by_id", "cors.#", "1"),
					resource.TestCheckTypeSetElemAttr("data.ovh_cloud_storage_swift_container.by_id", "cors.*", "https://www.example.com"),
					resource.TestCheckResourceAttrSet("data.ovh_cloud_storage_swift_container.by_id", "checksum"),
					resource.TestCheckResourceAttrSet("data.ovh_cloud_storage_swift_container.by_id", "created_at"),
					resource.TestCheckResourceAttrSet("data.ovh_cloud_storage_swift_container.by_id", "resource_status"),
					resource.TestCheckResourceAttr("data.ovh_cloud_storage_swift_container.by_id", "current_state.name", containerName),
					resource.TestCheckResourceAttr("data.ovh_cloud_storage_swift_container.by_id", "current_state.location.region", region),
					resource.TestCheckResourceAttr("data.ovh_cloud_storage_swift_container.by_id", "current_state.type", "PUBLIC"),
					resource.TestCheckResourceAttrSet("data.ovh_cloud_storage_swift_container.by_id", "current_state.static_url"),
					resource.TestCheckResourceAttrSet("data.ovh_cloud_storage_swift_container.by_id", "current_state.objects_count"),
					resource.TestCheckResourceAttrSet("data.ovh_cloud_storage_swift_container.by_id", "current_state.objects_size"),
				),
			},
		},
	})
}
