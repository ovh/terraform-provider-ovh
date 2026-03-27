package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCloudKeymanagerContainerConsumers_basic(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	region := os.Getenv("OVH_CLOUD_PROJECT_REGION_TEST")
	containerName := acctest.RandomWithPrefix(test_prefix)
	fakeResourceId := "00000000-0000-0000-0000-000000000001"

	config := fmt.Sprintf(`
resource "ovh_cloud_keymanager_container" "test" {
  service_name = "%s"
  region       = "%s"
  name         = "%s"
  type         = "GENERIC"
}

resource "ovh_cloud_keymanager_container_consumer" "test" {
  service_name  = "%s"
  container_id  = ovh_cloud_keymanager_container.test.id
  service       = "COMPUTE"
  resource_type = "INSTANCE"
  resource_id   = "%s"
}

data "ovh_cloud_keymanager_container_consumers" "test" {
  service_name = "%s"
  container_id = ovh_cloud_keymanager_container.test.id

  depends_on = [ovh_cloud_keymanager_container_consumer.test]
}
`, serviceName, region, containerName, serviceName, fakeResourceId, serviceName)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCloudKeymanager(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.ovh_cloud_keymanager_container_consumers.test", "consumers.#", "1"),
					resource.TestCheckResourceAttr("data.ovh_cloud_keymanager_container_consumers.test", "consumers.0.service", "COMPUTE"),
					resource.TestCheckResourceAttr("data.ovh_cloud_keymanager_container_consumers.test", "consumers.0.resource_type", "INSTANCE"),
					resource.TestCheckResourceAttr("data.ovh_cloud_keymanager_container_consumers.test", "consumers.0.resource_id", fakeResourceId),
					resource.TestCheckResourceAttrSet("data.ovh_cloud_keymanager_container_consumers.test", "consumers.0.id"),
				),
			},
		},
	})
}
