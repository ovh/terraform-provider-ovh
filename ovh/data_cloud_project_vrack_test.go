package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCloudProjectVrackDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCloud(t) },
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
resource "ovh_vrack_cloudproject" "attach" {
	service_name = "%s"
	project_id   = "%s"
}

data "ovh_cloud_project_vrack" "vrack" {
  service_name = "%s"
  depends_on = [ovh_vrack_cloudproject.attach]
}
`, os.Getenv("OVH_VRACK_SERVICE_TEST"), os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST"), os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST"))
