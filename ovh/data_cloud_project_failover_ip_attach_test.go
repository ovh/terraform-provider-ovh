package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const testAccDataSourceCloudProjectFailoverIpAttach = `
data "ovh_cloud_project_failover_ip_attach" "myfailoverip" {
 service_name = "%s"
 ip = "%s"
}
`

func TestAccDataSourceCloudProjectFailoverIpAttach(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	ipAddress := os.Getenv("OVH_CLOUD_PROJECT_FAILOVER_IP_TEST")
	config := fmt.Sprintf(
		testAccDataSourceCloudProjectFailoverIpAttach,
		serviceName,
		ipAddress,
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testAccPreCheckFailoverIpAttach(t) },

		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.ovh_cloud_project_failover_ip_attach.myfailoverip",
						"service_name",
						serviceName,
					),
					resource.TestCheckResourceAttr(
						"data.ovh_cloud_project_failover_ip_attach.myfailoverip",
						"ip",
						ipAddress,
					),
					resource.TestCheckResourceAttrSet(
						"data.ovh_cloud_project_failover_ip_attach.myfailoverip",
						"id",
					),
					resource.TestCheckResourceAttrSet(
						"data.ovh_cloud_project_failover_ip_attach.myfailoverip",
						"routed_to",
					),
				),
			},
		},
	})
}
