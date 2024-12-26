package ovh

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccResourceSavingsPlan_basic(t *testing.T) {
	displayName := acctest.RandomWithPrefix(test_prefix)
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")

	config := fmt.Sprintf(
		`resource "ovh_savings_plan" "sp" {
			service_name = "%s"
			flavor = "Rancher"
			period = "P1M"
			size = 1
			display_name = "%s"
		}`,
		serviceName,
		displayName,
	)

	endDate := time.Now().AddDate(0, 1, 1).Format(time.DateOnly)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckCloud(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_savings_plan.sp", "display_name", displayName),
					resource.TestCheckResourceAttr("ovh_savings_plan.sp", "flavor", "Rancher"),
					resource.TestCheckResourceAttr("ovh_savings_plan.sp", "size", "1"),
					resource.TestCheckResourceAttr("ovh_savings_plan.sp", "end_date", endDate),
					resource.TestCheckResourceAttr("ovh_savings_plan.sp", "period_end_date", endDate),
					resource.TestCheckResourceAttr("ovh_savings_plan.sp", "period_end_action", "TERMINATE"),
				),
			},
		},
	})
}
