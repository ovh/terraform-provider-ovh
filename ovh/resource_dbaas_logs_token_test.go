package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccResourceDbaasLogsToken_basic(t *testing.T) {
	serviceName := os.Getenv("OVH_DBAAS_LOGS_SERVICE_TEST")

	config := fmt.Sprintf(`
		resource "ovh_dbaas_logs_token" "tok" {
			service_name = "%s"
			name         = "TestToken"
		}
	`, serviceName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckDbaasLogs(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_dbaas_logs_token.tok", "name", "TestToken"),
					resource.TestCheckResourceAttr("ovh_dbaas_logs_token.tok", "service_name", serviceName),
					resource.TestCheckResourceAttrSet("ovh_dbaas_logs_token.tok", "value"),
					resource.TestCheckResourceAttrSet("ovh_dbaas_logs_token.tok", "token_id"),
					resource.TestCheckResourceAttrSet("ovh_dbaas_logs_token.tok", "cluster_id"),
				),
			},
		},
	})
}
