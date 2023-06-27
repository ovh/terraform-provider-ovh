package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const ACTIVE_MONTHLY_BILLING_INSTANCE_ID_TEST = "OVH_CLOUD_PROJECT_ACTIVE_MONTHLY_BILLING_INSTANCE_ID_TEST"

func TestAccCloudProjectInstanceActiveMonthlyBilling_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckActiveMonthlyBilling(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudProjectInstanceActiveMonthlyBillingConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"ovh_cloud_project_instance_amb.monthly", "monthly_billing_status", "ok"),
				),
			},
		},
	})
}

func testAccCloudProjectInstanceActiveMonthlyBillingConfig() string {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	instanceId := os.Getenv(ACTIVE_MONTHLY_BILLING_INSTANCE_ID_TEST)

	return fmt.Sprintf(
		testAccCloudProjectInstanceActiveMonthlyBillingConfig_Basic,
		serviceName,
		instanceId,
	)
}

const testAccCloudProjectInstanceActiveMonthlyBillingConfig_Basic = `
resource "ovh_cloud_project_instance_amb" "monthly" {
  service_name  = "%s"
  instance_id  = "%s"
  wait_activation = true
}
`
