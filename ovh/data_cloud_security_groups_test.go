package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceCloudSecurityGroups_basic(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	region := os.Getenv("OVH_CLOUD_PROJECT_REGION_TEST")

	name := acctest.RandomWithPrefix(testAccResourceCloudSecurityGroupNamePrefix)

	config := fmt.Sprintf(`
resource "ovh_cloud_security_group" "test" {
  service_name = "%s"
  region       = "%s"
  name         = "%s"
  description  = "Test security group"

  rule = [
    {
      direction        = "INGRESS"
      ethernet_type    = "IPV4"
      protocol         = "TCP"
      port_range_min   = 22
      port_range_max   = 22
      remote_ip_prefix = "0.0.0.0/0"
      description      = "SSH"
    },
  ]
}

data "ovh_cloud_security_groups" "test" {
  service_name = ovh_cloud_security_group.test.service_name
}
`, serviceName, region, name)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCloudSecurityGroup(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.ovh_cloud_security_groups.test", "service_name", serviceName),
					resource.TestCheckResourceAttrSet("data.ovh_cloud_security_groups.test", "security_groups.#"),
				),
			},
		},
	})
}
