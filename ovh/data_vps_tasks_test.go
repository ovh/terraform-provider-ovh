package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccVPSTasksDataSource_basic(t *testing.T) {
	vps := os.Getenv("OVH_VPS")
	config := fmt.Sprintf(testAccVPSTasksDatasourceConfig_Basic, vps)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckVPS(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.ovh_vps_tasks.tasks", "service_name", vps),
					resource.TestCheckResourceAttrSet(
						"data.ovh_vps_tasks.tasks", "task_ids.#"),
				),
			},
		},
	})
}

const testAccVPSTasksDatasourceConfig_Basic = `
data "ovh_vps_tasks" "tasks" {
  service_name = "%s"
}
`
