package ovh

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const testAccDataSourceDbaasLogsInputEngine_deprecated = `
data "ovh_dbaas_logs_input_engine" "logstash" {
 name          = "logstash"
 version       = "6.8"
 is_deprecated = true
}
`
const testAccDataSourceDbaasLogsInputEngine_basic = `
data "ovh_dbaas_logs_input_engine" "logstash" {
 name          = "logstash"
 version       = "7.x"
}
`

func TestAccDataSourceDbaasLogsInputEngine_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() { testAccPreCheckCredentials(t) },

		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceDbaasLogsInputEngine_basic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.ovh_dbaas_logs_input_engine.logstash",
						"is_deprecated",
						"false",
					),
					resource.TestCheckResourceAttr(
						"data.ovh_dbaas_logs_input_engine.logstash",
						"version",
						"7.x",
					),
				),
			},
		},
	})
}

func TestAccDataSourceDbaasLogsInputEngine_deprecated(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() { testAccPreCheckCredentials(t) },

		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceDbaasLogsInputEngine_deprecated,
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
