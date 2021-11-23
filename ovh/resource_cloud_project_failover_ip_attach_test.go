package ovh

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"os"
	"testing"
)

const testAccCloudProjectFailoverIpAttach = `
resource "ovh_cloud_project_failover_ip_attach" "myfailoverip" {
 service_name = "%s"
 ip = "%s"
 routed_to = "%s"
}
`

func TestAccResourceCloudProjectFailoverIpAttach(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	ipAddress := os.Getenv("OVH_CLOUD_PROJECT_FAILOVER_IP_TEST")
	routedTo1 := os.Getenv("OVH_CLOUD_PROJECT_FAILOVER_IP_ROUTED_TO_1_TEST")
	config := fmt.Sprintf(
		testAccCloudProjectFailoverIpAttach,
		serviceName,
		ipAddress,
		routedTo1,
	)
	resource.Test(t, resource.TestCase{
		PreCheck: func() { testAccPreCheckFailoverIpAttach(t) },

		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"ovh_cloud_project_failover_ip_attach.myfailoverip",
						"service_name",
						serviceName,
					),
					resource.TestCheckResourceAttr(
						"ovh_cloud_project_failover_ip_attach.myfailoverip",
						"ip",
						ipAddress,
					),
					resource.TestCheckResourceAttr(
						"ovh_cloud_project_failover_ip_attach.myfailoverip",
						"routed_to",
						routedTo1,
					),
				),
			},
		},
	})
	routedTo2 := os.Getenv("OVH_CLOUD_PROJECT_FAILOVER_IP_ROUTED_TO_2_TEST")
	config = fmt.Sprintf(
		testAccCloudProjectFailoverIpAttach,
		serviceName,
		ipAddress,
		routedTo2,
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testAccPreCheckFailoverIpAttach(t) },

		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"ovh_cloud_project_failover_ip_attach.myfailoverip",
						"service_name",
						serviceName,
					),
					resource.TestCheckResourceAttr(
						"ovh_cloud_project_failover_ip_attach.myfailoverip",
						"ip",
						ipAddress,
					),
					resource.TestCheckResourceAttr(
						"ovh_cloud_project_failover_ip_attach.myfailoverip",
						"routed_to",
						routedTo2,
					),
				),
			},
		},
	})
}
