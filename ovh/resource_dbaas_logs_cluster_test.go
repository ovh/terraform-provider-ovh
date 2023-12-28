package ovh

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func init() {
	resource.AddTestSweepers("ovh_dbaas_logs_cluster", &resource.Sweeper{
		Name: "ovh_dbaas_logs_cluster",
		F:    testSweepDbaasLogsCluster,
	})
}

func testSweepDbaasLogsCluster(region string) error {
	serviceName := os.Getenv("OVH_DBAAS_LOGS_SERVICE_TEST")
	if serviceName == "" {
		log.Print("[DEBUG] OVH_DBAAS_LOGS_SERVICE_TEST is not set. No LDP cluster to sweep")
		return nil
	}

	// Nothing to sweep as LDP dedicated cluster can't be created/deleted thru API

	return nil
}

func TestAccDbaasLogsCluster(t *testing.T) {
	serviceName := os.Getenv("OVH_DBAAS_LOGS_SERVICE_TEST")
	clusterId := os.Getenv("OVH_DBAAS_LOGS_CLUSTER_ID")

	config := fmt.Sprintf(
		testAccDbaasLogsClusterConfig,
		serviceName,
		clusterId,
	)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckDbaasLogsCluster(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"ovh_dbaas_logs_cluster.ldp", "service_name", serviceName),
					resource.TestCheckResourceAttrSet(
						"ovh_dbaas_logs_cluster.ldp", "dedicated_input_pem"),
					resource.TestCheckResourceAttrSet(
						"ovh_dbaas_logs_cluster.ldp", "direct_input_pem"),
					resource.TestCheckTypeSetElemAttr(
						"ovh_dbaas_logs_cluster.ldp", "archive_allowed_networks.*", "10.0.0.0/16",
					),
					resource.TestCheckTypeSetElemAttr(
						"ovh_dbaas_logs_cluster.ldp", "direct_input_allowed_networks.*", "10.0.0.0/16",
					),
					resource.TestCheckTypeSetElemAttr(
						"ovh_dbaas_logs_cluster.ldp", "query_allowed_networks.*", "10.0.0.0/16",
					),
				),
			},
		},
	})
}

const testAccDbaasLogsClusterConfig = `
resource "ovh_dbaas_logs_cluster" "ldp" {
	service_name = "%s"
	cluster_id = "%s"

	archive_allowed_networks       = ["10.0.0.0/16"]
	direct_input_allowed_networks  = ["10.0.0.0/16"]
	query_allowed_networks         = ["10.0.0.0/16"]
}
`
