package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCloudKeyManagerSecretPayloadDataSource_basic(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	region := os.Getenv("OVH_CLOUD_PROJECT_REGION_TEST")
	name := acctest.RandomWithPrefix(test_prefix)

	config := fmt.Sprintf(`
resource "ovh_cloud_key_manager_secret" "test" {
  service_name         = "%s"
  region               = "%s"
  name                 = "%s"
  secret_type          = "OPAQUE"
  payload              = base64encode("my-secret-payload")
  payload_content_type = "APPLICATION_OCTET_STREAM"
}

data "ovh_cloud_key_manager_secret_payload" "test" {
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
					resource.TestCheckResourceAttr("data.ovh_cloud_key_manager_secret_payload.test", "service_name", serviceName),
					resource.TestCheckResourceAttrSet("data.ovh_cloud_key_manager_secret_payload.test", "secret_id"),
					resource.TestCheckResourceAttrSet("data.ovh_cloud_key_manager_secret_payload.test", "payload"),
				),
			},
		},
	})
}
