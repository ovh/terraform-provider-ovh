package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccDbaasLogsOutputOpensearchIndex_importBasic(t *testing.T) {
	serviceName := os.Getenv("OVH_DBAAS_LOGS_SERVICE_TEST")
	desc := acctest.RandomWithPrefix(test_prefix)
	suffix := acctest.RandomWithPrefix(test_prefix)

	config := fmt.Sprintf(
		testAccDbaasLogsOutputOpensearchIndexConfig,
		serviceName,
		desc,
		suffix,
	)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckDbaasLogs(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
			},
			{
				ResourceName:      "ovh_dbaas_logs_output_opensearch_index.index",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccDbaasLogsOutputOpensearchIndexImportId("ovh_dbaas_logs_output_opensearch_index.index"),
			},
		},
	})
}

func testAccDbaasLogsOutputOpensearchIndexImportId(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		testIndex, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("ovh_dbaas_logs_output_opensearch_index not found: %s", resourceName)
		}
		return fmt.Sprintf(
			"%s/%s",
			testIndex.Primary.Attributes["service_name"],
			testIndex.Primary.Attributes["id"],
		), nil
	}
}

const testAccDbaasLogsOutputOpensearchIndexConfig = `
resource "ovh_dbaas_logs_output_opensearch_index" "index" {
	service_name = "%s"
	description  = "%s"
	suffix       = "%s"
	nb_shard     = 1
}
`
