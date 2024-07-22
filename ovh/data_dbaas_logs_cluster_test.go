package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceDbaasLogsCluster(t *testing.T) {
	serviceName := os.Getenv("OVH_DBAAS_LOGS_SERVICE_TEST")
	clusterId := os.Getenv("OVH_DBAAS_LOGS_CLUSTER_ID")

	config := fmt.Sprintf(`
		data "ovh_dbaas_logs_cluster" "ldp" {
			service_name = "%s"
			cluster_id   = "%s"
		}`,
		serviceName,
		clusterId,
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testAccPreCheckDbaasLogsCluster(t) },

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
					resource.TestCheckResourceAttr(
						"data.ovh_dbaas_logs_cluster.ldp",
						"cluster_id",
						clusterId,
					),
				),
			},
		},
	})
}

func TestAccDataSourceDbaasLogsClusterDefault(t *testing.T) {
	serviceName := os.Getenv("OVH_DBAAS_LOGS_SERVICE_TEST")
	clusterId := os.Getenv("OVH_DBAAS_LOGS_CLUSTER_ID")

	config := fmt.Sprintf(`
		data "ovh_dbaas_logs_cluster" "ldp" {
			service_name = "%s"
		}`,
		serviceName,
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testAccPreCheckDbaasLogsCluster(t) },

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
					resource.TestCheckResourceAttr(
						"data.ovh_dbaas_logs_cluster.ldp",
						"cluster_id",
						clusterId,
					),
				),
			},
		},
	})
}
