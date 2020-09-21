package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccCloudDBEnterprise(t *testing.T) {
	clusterId := os.Getenv("OVH_CLOUDDB_ENTERPRISE")
	config := fmt.Sprintf(testAccCloudDBEnterpriseDatasourceConfig, clusterId)
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCloudDBEnterprise(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.ovh_clouddb_enterprise.db", "cluster_id", clusterId),
					resource.TestCheckResourceAttr(
						"data.ovh_clouddb_enterprise.db", "status", string(CloudDBEnterpriseStatusCreated)),
				),
			},
		},
	})
}

const testAccCloudDBEnterpriseDatasourceConfig = `
data "ovh_clouddb_enterprise" "db" {
  cluster_id = "%s"
}
`
