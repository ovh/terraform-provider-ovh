package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

var testAccCloudProjectGatewayConfig = `

resource "ovh_cloud_project_network_private" "mypriv" {
  service_name  = "%s"
  vlan_id       = "%d"
  name          = "%s"
  regions       = ["%s"]
}

resource "ovh_cloud_project_network_private_subnet" "myprivsub" {
  service_name  = ovh_cloud_project_network_private.mypriv.service_name
  network_id    = ovh_cloud_project_network_private.mypriv.id
  region        = "%s"
  start         = "10.0.0.2"
  end           = "10.0.255.254"
  network       = "10.0.0.0/16"
  dhcp          = true
}

resource "ovh_cloud_project_gateway" "gateway" {
  service_name = ovh_cloud_project_network_private.mypriv.service_name
  name          = "%s"
  model         = "s"
  region        = ovh_cloud_project_network_private_subnet.myprivsub.region
  network_id    = tolist(ovh_cloud_project_network_private.mypriv.regions_attributes[*].openstackid)[0]
  subnet_id     = ovh_cloud_project_network_private_subnet.myprivsub.id
}
`

func TestAccCloudProjectGateway(t *testing.T) {
	name := acctest.RandomWithPrefix(test_prefix)
	gatewayName := acctest.RandomWithPrefix(test_prefix)
	vlanId := acctest.RandIntRange(100, 200)
	region := os.Getenv("OVH_CLOUD_PROJECT_KUBE_REGION_TEST")
	resourcePath := "ovh_cloud_project_gateway.gateway"

	config := fmt.Sprintf(
		testAccCloudProjectGatewayConfig,
		os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST"),
		vlanId,
		name,
		region,
		region,
		gatewayName,
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCloud(t)
			testAccCheckCloudProjectExists(t)
		},

		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourcePath, "service_name"),
					resource.TestCheckResourceAttrSet(resourcePath, "network_id"),
					resource.TestCheckResourceAttrSet(resourcePath, "subnet_id"),
					resource.TestCheckResourceAttrSet(resourcePath, "model"),
					resource.TestCheckResourceAttrSet(resourcePath, "external_information.0.network_id"),
					resource.TestCheckResourceAttrSet(resourcePath, "external_information.0.ips.0.ip"),
					resource.TestCheckResourceAttrSet(resourcePath, "external_information.0.ips.0.subnet_id"),
					resource.TestCheckResourceAttrSet(resourcePath, "interfaces.0.id"),
					resource.TestCheckResourceAttrSet(resourcePath, "interfaces.0.ip"),
					resource.TestCheckResourceAttrSet(resourcePath, "interfaces.0.subnet_id"),
					resource.TestCheckResourceAttrSet(resourcePath, "interfaces.0.network_id"),
					resource.TestCheckResourceAttr(resourcePath, "region", region),
					resource.TestCheckResourceAttr(resourcePath, "name", gatewayName),
					resource.TestCheckResourceAttr(resourcePath, "model", "s"),
				),
			},
			{
				ResourceName:            "ovh_cloud_project_gateway.gateway",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"network_id", "subnet_id"},
				ImportStateIdFunc:       testAccCloudProjectGatewayImportId("ovh_cloud_project_gateway.gateway"),
			},
			{
				Config: config,
			},
		},
	})
}

func testAccCloudProjectGatewayImportId(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		gateway, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("gateway not found: %s", resourceName)
		}

		return fmt.Sprintf(
			"%s/%s/%s",
			gateway.Primary.Attributes["service_name"],
			gateway.Primary.Attributes["region"],
			gateway.Primary.ID,
		), nil
	}
}
