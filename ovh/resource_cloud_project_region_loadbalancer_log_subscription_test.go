package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCloudProjectLoadBalancerLogSubscription_basic(t *testing.T) {
	dbaasLogsServiceName := os.Getenv("OVH_DBAAS_LOGS_SERVICE_TEST")
	cloudProject := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	loadBalancerId := os.Getenv("OVH_CLOUD_LOADBALANCER_ID_TEST")

	var testCreateLoadBalancerLogSubscription = fmt.Sprintf(`
		resource "ovh_dbaas_logs_output_graylog_stream" "stream" {
			service_name = "%s"
			title        = "%s"
			description  = "%s"
		}

		resource "ovh_cloud_project_region_loadbalancer_log_subscription" "CreateLogSubscription" {
			service_name = "%s"
			region_name = "GRA11"
			loadbalancer_id = "%s"
			kind = "haproxy"
			stream_id = ovh_dbaas_logs_output_graylog_stream.stream.stream_id
		}
	`, dbaasLogsServiceName, acctest.RandomWithPrefix(test_prefix), acctest.RandomWithPrefix(test_prefix), cloudProject, loadBalancerId)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCloud(t)
			testAccCheckCloudProjectExists(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testCreateLoadBalancerLogSubscription,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("ovh_cloud_project_region_loadbalancer_log_subscription.CreateLogSubscription", "created_at"),
					resource.TestCheckResourceAttrSet("ovh_cloud_project_region_loadbalancer_log_subscription.CreateLogSubscription", "subscription_id"),
					resource.TestCheckResourceAttrSet("ovh_cloud_project_region_loadbalancer_log_subscription.CreateLogSubscription", "resource_type"),
					resource.TestCheckResourceAttrSet("ovh_cloud_project_region_loadbalancer_log_subscription.CreateLogSubscription", "resource_name"),
					resource.TestCheckResourceAttrSet("ovh_cloud_project_region_loadbalancer_log_subscription.CreateLogSubscription", "operation_id"),
				),
			},
		},
	})
}
