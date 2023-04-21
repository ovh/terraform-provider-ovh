package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const testAccDataSourceDbaasLogsInputEngine_deprecated = `
data "ovh_dbaas_logs_input_engine" "logstash" {
 service_name  = "%s"
 name          = "%s"
 version       = "%s"
 is_deprecated = "%s"
}
`

const testAccDataSourceDbaasLogsInputEngine_basic = `
data "ovh_dbaas_logs_input_engine" "logstash" {
 service_name  = "%s"
 name          = "%s"
 version       = "%s"
}
`

func TestAccDbaasLogsInputEngineDataSource_basic(t *testing.T) {
	serviceName := os.Getenv("OVH_DBAAS_LOGS_SERVICE_TEST")
	name := "LOGSTASH"
	// version := "7.x"
	version := os.Getenv("OVH_DBAAS_LOGS_LOGSTASH_VERSION_TEST")

	config := fmt.Sprintf(
		testAccDataSourceDbaasLogsInputEngine_basic,
		serviceName,
		name,
		version,
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testAccPreCheckDbaasLogsInput(t) },

		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.ovh_dbaas_logs_input_engine.logstash",
						"is_deprecated",
						"false",
					),
					resource.TestCheckResourceAttr(
						"data.ovh_dbaas_logs_input_engine.logstash",
						"version",
						version,
					),
				),
			},
		},
	})
}

func TestAccDbaasLogsInputEngineDataSource_deprecated(t *testing.T) {
	serviceName := os.Getenv("OVH_DBAAS_LOGS_SERVICE_TEST")
	name := "LOGSTASH"
	version := "6.8"
	is_deprecated := "true"

	config := fmt.Sprintf(
		testAccDataSourceDbaasLogsInputEngine_deprecated,
		serviceName,
		name,
		version,
		is_deprecated,
	)
	resource.Test(t, resource.TestCase{
		PreCheck: func() { testAccPreCheckDbaasLogs(t) },

		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.ovh_dbaas_logs_input_engine.logstash",
						"is_deprecated",
						"true",
					),
					resource.TestCheckResourceAttr(
						"data.ovh_dbaas_logs_input_engine.logstash",
						"version",
						"6.8",
					),
				),
			},
		},
	})
}
