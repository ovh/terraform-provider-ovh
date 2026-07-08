package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceCloudGateway_basic(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	region := os.Getenv("OVH_CLOUD_PROJECT_REGION_TEST")
	gatewayModel := os.Getenv("OVH_CLOUD_GATEWAY_MODEL_TEST")
	if gatewayModel == "" {
		gatewayModel = "S"
	}

	gatewayName := acctest.RandomWithPrefix(testAccResourceCloudGatewayNamePrefix)

	config := fmt.Sprintf(`
resource "ovh_cloud_gateway" "test" {
  service_name = "%s"
  name         = "%s"
  region       = "%s"

  external_gateway = {
    enabled = true
    model   = "%s"
  }
}

data "ovh_cloud_gateway" "test" {
  service_name = ovh_cloud_gateway.test.service_name
  id           = ovh_cloud_gateway.test.id
}
`, serviceName, gatewayName, region, gatewayModel)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCloudGateway(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.ovh_cloud_gateway.test", "service_name", serviceName),
					resource.TestCheckResourceAttr("data.ovh_cloud_gateway.test", "name", gatewayName),
					resource.TestCheckResourceAttr("data.ovh_cloud_gateway.test", "location.region", region),
					resource.TestCheckResourceAttr("data.ovh_cloud_gateway.test", "external_gateway.enabled", "true"),
					resource.TestCheckResourceAttr("data.ovh_cloud_gateway.test", "external_gateway.model", gatewayModel),
					resource.TestCheckResourceAttrSet("data.ovh_cloud_gateway.test", "id"),
					resource.TestCheckResourceAttrSet("data.ovh_cloud_gateway.test", "checksum"),
					resource.TestCheckResourceAttrSet("data.ovh_cloud_gateway.test", "created_at"),
					resource.TestCheckResourceAttr("data.ovh_cloud_gateway.test", "resource_status", "READY"),
					resource.TestCheckResourceAttrSet("data.ovh_cloud_gateway.test", "current_state.name"),
					resource.TestCheckResourceAttrSet("data.ovh_cloud_gateway.test", "current_state.status"),
				),
			},
		},
	})
}
