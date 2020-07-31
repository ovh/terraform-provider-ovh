package ovh

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"os"
	"testing"
)

func TestAccDedicatedCephDatasource(t *testing.T) {
	dedicated_ceph := os.Getenv("OVH_DEDICATED_CEPH")
	config := fmt.Sprintf(testAccDedicatedCephDatasourceConfig, dedicated_ceph)
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckDedicatedCeph(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.ovh_dedicated_ceph.ceph", "service_name", dedicated_ceph),
					resource.TestCheckResourceAttr(
						"data.ovh_dedicated_ceph.ceph", "status", "INSTALLED"),
				),
			},
		},
	})
}

const testAccDedicatedCephDatasourceConfig = `
data "ovh_dedicated_ceph" "ceph" {
  service_name = "%s"
}
`
