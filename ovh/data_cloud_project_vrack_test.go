package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccCloudProjectVrackDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCloudWithVrack(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudProjectVrackDatasourceConfig,
				Check:  resource.TestCheckResourceAttrSet("data.ovh_cloud_project_vrack.vrack", "id"),
			},
		},
	})
}

var testAccCloudProjectVrackDatasourceConfig = fmt.Sprintf(`
data "ovh_cloud_project_vrack" "vrack" {
  service_name = "%s"
}
`, os.Getenv("OVH_CLOUD_PROJECT_WITH_VRACK_SERVICE_TEST"))
