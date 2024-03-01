package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const testAccCloudProjectAlertingConfig = `
resource "ovh_cloud_project_alerting" "alert" {
	service_name = "%s"
	delay = 259200
	monthly_threshold = 3000
	email = "some.test@ovhcloud.com"
}
`

const testAccCloudProjectAlertingUpdatedConfig = `
resource "ovh_cloud_project_alerting" "alert" {
	service_name = "%s"
	delay = 604800
	monthly_threshold = 100
	email = "some.test@ovhcloud.com"
}
`

func TestAccCloudProjectAlerting_basic(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")

	config := fmt.Sprintf(testAccCloudProjectAlertingConfig, serviceName)
	updatedConfig := fmt.Sprintf(testAccCloudProjectAlertingUpdatedConfig, serviceName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckCloud(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"ovh_cloud_project_alerting.alert", "delay", "259200"),
					resource.TestCheckResourceAttr(
						"ovh_cloud_project_alerting.alert", "monthly_threshold", "3000"),
					resource.TestCheckResourceAttr(
						"ovh_cloud_project_alerting.alert", "email", "some.test@ovhcloud.com"),
					resource.TestCheckResourceAttr(
						"ovh_cloud_project_alerting.alert", "formatted_monthly_threshold.currency_code", "EUR"),
					resource.TestCheckResourceAttr(
						"ovh_cloud_project_alerting.alert", "formatted_monthly_threshold.text", "3000.00 €"),
					resource.TestCheckResourceAttr(
						"ovh_cloud_project_alerting.alert", "formatted_monthly_threshold.value", "3000"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"ovh_cloud_project_alerting.alert", "delay", "604800"),
					resource.TestCheckResourceAttr(
						"ovh_cloud_project_alerting.alert", "monthly_threshold", "100"),
					resource.TestCheckResourceAttr(
						"ovh_cloud_project_alerting.alert", "email", "some.test@ovhcloud.com"),
					resource.TestCheckResourceAttr(
						"ovh_cloud_project_alerting.alert", "formatted_monthly_threshold.currency_code", "EUR"),
					resource.TestCheckResourceAttr(
						"ovh_cloud_project_alerting.alert", "formatted_monthly_threshold.text", "100.00 €"),
					resource.TestCheckResourceAttr(
						"ovh_cloud_project_alerting.alert", "formatted_monthly_threshold.value", "100"),
				),
			},
			{
				Config:  updatedConfig,
				Destroy: true,
			},
		},
	})
}
