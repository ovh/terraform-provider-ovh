package ovh

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// The provider cannot upload objects, so the data source is exercised on an
// empty container: the route answers and the list is empty.
func TestAccDataSourceCloudStorageSwiftContainerObjects_basic(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	region := strings.ToUpper(os.Getenv("OVH_CLOUD_PROJECT_STORAGE_REGION_TEST"))

	containerName := acctest.RandomWithPrefix(test_prefix)

	config := fmt.Sprintf(`
resource "ovh_cloud_storage_swift_container" "container" {
  service_name = "%s"
  name         = "%s"
  region       = "%s"
}

data "ovh_cloud_storage_swift_container_objects" "objects" {
  service_name = "%s"
  container_id = ovh_cloud_storage_swift_container.container.id
}

data "ovh_cloud_storage_swift_container_objects" "filtered" {
  service_name = "%s"
  container_id = ovh_cloud_storage_swift_container.container.id
  prefix       = "documents/"
  limit        = 10
}
`, serviceName, containerName, region, serviceName, serviceName)

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
					resource.TestCheckResourceAttr("data.ovh_cloud_storage_swift_container_objects.objects", "service_name", serviceName),
					resource.TestCheckResourceAttr("data.ovh_cloud_storage_swift_container_objects.objects", "container_id", region+"_"+containerName),
					resource.TestCheckResourceAttr("data.ovh_cloud_storage_swift_container_objects.objects", "objects.#", "0"),
					resource.TestCheckResourceAttr("data.ovh_cloud_storage_swift_container_objects.filtered", "prefix", "documents/"),
					resource.TestCheckResourceAttr("data.ovh_cloud_storage_swift_container_objects.filtered", "limit", "10"),
					resource.TestCheckResourceAttr("data.ovh_cloud_storage_swift_container_objects.filtered", "objects.#", "0"),
				),
			},
		},
	})
}
