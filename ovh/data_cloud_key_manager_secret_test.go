package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCloudKeyManagerSecretDataSource_basic(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	region := os.Getenv("OVH_CLOUD_PROJECT_REGION_TEST")
	name := acctest.RandomWithPrefix(test_prefix)

	config := fmt.Sprintf(`
resource "ovh_cloud_key_manager_secret" "test" {
  service_name = "%s"
  region       = "%s"
  name         = "%s"
  secret_type  = "OPAQUE"
}

data "ovh_cloud_key_manager_secret" "test" {
  service_name = "%s"
  secret_id    = ovh_cloud_key_manager_secret.test.id
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
					resource.TestCheckResourceAttr("data.ovh_cloud_key_manager_secret.test", "service_name", serviceName),
					resource.TestCheckResourceAttr("data.ovh_cloud_key_manager_secret.test", "location.region", region),
					resource.TestCheckResourceAttr("data.ovh_cloud_key_manager_secret.test", "name", name),
					resource.TestCheckResourceAttr("data.ovh_cloud_key_manager_secret.test", "secret_type", "OPAQUE"),
					resource.TestCheckResourceAttrSet("data.ovh_cloud_key_manager_secret.test", "id"),
					resource.TestCheckResourceAttrSet("data.ovh_cloud_key_manager_secret.test", "checksum"),
					resource.TestCheckResourceAttrSet("data.ovh_cloud_key_manager_secret.test", "created_at"),
					resource.TestCheckResourceAttr("data.ovh_cloud_key_manager_secret.test", "resource_status", "READY"),
				),
			},
		},
	})
}
