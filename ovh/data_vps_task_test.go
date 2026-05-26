package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccVPSTaskDataSource_basic(t *testing.T) {
	vps := os.Getenv("OVH_VPS")
	taskId := os.Getenv("OVH_VPS_TASK_ID")
	config := fmt.Sprintf(testAccVPSTaskDatasourceConfig_Basic, vps, taskId)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckVPS(t)
			if taskId == "" {
				t.Skip("OVH_VPS_TASK_ID is not set")
			}
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.ovh_vps_task.task", "service_name", vps),
					resource.TestCheckResourceAttrSet(
						"data.ovh_vps_task.task", "state"),
					resource.TestCheckResourceAttrSet(
						"data.ovh_vps_task.task", "type"),
				),
			},
		},
	})
}

const testAccVPSTaskDatasourceConfig_Basic = `
data "ovh_vps_task" "task" {
  service_name = "%s"
  id           = %s
}
`
