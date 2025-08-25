package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccResourceDbaasLogsToken_basic(t *testing.T) {
	serviceName := os.Getenv("OVH_DBAAS_LOGS_SERVICE_TEST")
	tokenName := acctest.RandomWithPrefix(test_prefix)

	config := fmt.Sprintf(`
		resource "ovh_dbaas_logs_token" "tok" {
			service_name = "%s"
			name         = "%s"
		}
	`, serviceName, tokenName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckDbaasLogs(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_dbaas_logs_token.tok", "name", tokenName),
					resource.TestCheckResourceAttr("ovh_dbaas_logs_token.tok", "service_name", serviceName),
					resource.TestCheckResourceAttrSet("ovh_dbaas_logs_token.tok", "value"),
					resource.TestCheckResourceAttrSet("ovh_dbaas_logs_token.tok", "token_id"),
					resource.TestCheckResourceAttrSet("ovh_dbaas_logs_token.tok", "cluster_id"),
					resource.TestCheckResourceAttrSet("ovh_dbaas_logs_token.tok", "id"),
				),
			},
			{
				Config:            config,
				ImportState:       true,
				ImportStateVerify: true,
				ResourceName:      "ovh_dbaas_logs_token.tok",
				ImportStateIdFunc: testAccDbaasLogsTokenImportId("ovh_dbaas_logs_token.tok"),
			},
		},
	})
}

func testAccDbaasLogsTokenImportId(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		testDatabase, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("resource not found: %s", resourceName)
		}
		return fmt.Sprintf(
			"%s/%s",
			testDatabase.Primary.Attributes["service_name"],
			testDatabase.Primary.Attributes["token_id"],
		), nil
	}
}
