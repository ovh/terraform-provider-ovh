package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func testAccCloudKeymanagerContainerConsumerImportStateIdFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}
		return rs.Primary.Attributes["id"], nil
	}
}

func TestAccCloudKeymanagerContainerConsumer_basic(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	region := os.Getenv("OVH_CLOUD_PROJECT_REGION_TEST")
	containerName := acctest.RandomWithPrefix(test_prefix)
	// Use a fake resource ID for the consumer target — the API accepts any UUID
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
`, serviceName, region, containerName, serviceName, fakeResourceId)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCloudKeymanager(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_keymanager_container_consumer.test", "service_name", serviceName),
					resource.TestCheckResourceAttr("ovh_cloud_keymanager_container_consumer.test", "service", "COMPUTE"),
					resource.TestCheckResourceAttr("ovh_cloud_keymanager_container_consumer.test", "resource_type", "INSTANCE"),
					resource.TestCheckResourceAttr("ovh_cloud_keymanager_container_consumer.test", "resource_id", fakeResourceId),
					resource.TestCheckResourceAttrSet("ovh_cloud_keymanager_container_consumer.test", "id"),
					resource.TestCheckResourceAttrSet("ovh_cloud_keymanager_container_consumer.test", "container_id"),
				),
			},
			// Test import
			{
				ResourceName:      "ovh_cloud_keymanager_container_consumer.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccCloudKeymanagerContainerConsumerImportStateIdFunc("ovh_cloud_keymanager_container_consumer.test"),
			},
		},
	})
}
