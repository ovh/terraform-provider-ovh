package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const testAccDataSourceDbaasLogsOutputGraylogStream_basic = `
resource "ovh_dbaas_logs_output_graylog_stream" "stream" {
 service_name = "%s"
 title        = "%s"
 description  = "%s"
}

data "ovh_dbaas_logs_output_graylog_stream" "stream" {
 service_name = ovh_dbaas_logs_output_graylog_stream.stream.service_name
 title        = ovh_dbaas_logs_output_graylog_stream.stream.title
}
`

func TestAccDataSourceDbaasLogsOutputGraylogStream_basic(t *testing.T) {
	serviceName := os.Getenv("OVH_DBAAS_LOGS_SERVICE_TEST")
	title := acctest.RandomWithPrefix(test_prefix)
	desc := acctest.RandomWithPrefix(test_prefix)

	config := fmt.Sprintf(
		testAccDataSourceDbaasLogsOutputGraylogStream_basic,
		serviceName,
		title,
		desc,
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testAccPreCheckDbaasLogs(t) },

		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.ovh_dbaas_logs_output_graylog_stream.stream",
						"description",
						desc,
					),
					resource.TestCheckResourceAttr(
						"data.ovh_dbaas_logs_output_graylog_stream.stream",
						"title",
						title,
					),
				),
			},
		},
	})
}
