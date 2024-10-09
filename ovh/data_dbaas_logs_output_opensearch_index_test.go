package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceDbaasLogsOutputOpensearchIndex_basic(t *testing.T) {
	serviceName := os.Getenv("OVH_DBAAS_LOGS_SERVICE_TEST")
	desc := acctest.RandomWithPrefix(test_prefix)
	suffix := acctest.RandomWithPrefix(test_prefix)

	config := fmt.Sprintf(`
		resource "ovh_dbaas_logs_output_opensearch_index" "idx" {
			service_name = "%s"
			description  = "%s"
			suffix        = "%s"
			nb_shard = 1
		}

		data "ovh_dbaas_logs_output_opensearch_index" "idx" {
			service_name = ovh_dbaas_logs_output_opensearch_index.idx.service_name
			name        = ovh_dbaas_logs_output_opensearch_index.idx.name
		}`,
		serviceName,
		desc,
		suffix,
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testAccPreCheckDbaasLogs(t) },

		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.ovh_dbaas_logs_output_opensearch_index.idx",
						"description",
						desc,
					),
				),
			},
		},
	})
}
