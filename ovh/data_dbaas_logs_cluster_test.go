package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const testAccDataSourceDbaasLogsCluster = `
data "ovh_dbaas_logs_cluster" "ldp" {
  service_name = "%s"
}
`

func TestAccDataSourceDbaasLogsCluster(t *testing.T) {
	serviceName := os.Getenv("OVH_DBAAS_LOGS_SERVICE_TEST")

	config := fmt.Sprintf(
		testAccDataSourceDbaasLogsCluster,
		serviceName,
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testAccPreCheckDbaasLogs(t) },

		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.ovh_dbaas_logs_cluster.ldp",
						"service_name",
						serviceName,
					),
				),
			},
		},
	})
}
