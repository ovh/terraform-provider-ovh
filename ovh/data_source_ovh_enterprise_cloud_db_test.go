package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccEnterpriseCloudDB(t *testing.T) {
	clusterId := os.Getenv("OVH_ENTERPRISE_CLOUD_DB")
	config := fmt.Sprintf(testAccEnterpriseCloudDBDatasourceConfig, clusterId)
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckEnterpriseCloudDB(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.ovh_enterprise_cloud_db.db", "cluster_id", clusterId),
					resource.TestCheckResourceAttr(
						"data.ovh_enterprise_cloud_db.db", "status", string(EnterpriseCloudDBStatusCreated)),
				),
			},
		},
	})
}

const testAccEnterpriseCloudDBDatasourceConfig = `
data "ovh_enterprise_cloud_db" "db" {
  cluster_id = "%s"
}
`
