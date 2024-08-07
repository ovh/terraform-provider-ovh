package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCloudProjectLoadBalancerGetLogSubscription_basic(t *testing.T) {
	config := fmt.Sprintf(testAccCloudProjectSubscription,
		os.Getenv("OVH_DBAAS_LOGS_SERVICE_TEST"),
		acctest.RandomWithPrefix(test_prefix),
		acctest.RandomWithPrefix(test_prefix),
		os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST"),
		os.Getenv("OVH_CLOUD_LOADBALANCER_ID_TEST"),
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCloud(t)
			testAccCheckCloudProjectExists(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.ovh_cloud_project_region_loadbalancer_log_subscription.test", "created_at"),
					resource.TestCheckResourceAttrSet("data.ovh_cloud_project_region_loadbalancer_log_subscription.test", "resource_name"),
					resource.TestCheckResourceAttr("data.ovh_cloud_project_region_loadbalancer_log_subscription.test", "kind", "haproxy"),
				),
			},
		},
	})
}

var testAccCloudProjectSubscription = `
resource "ovh_dbaas_logs_output_graylog_stream" "stream" {
	service_name = "%s"
	title        = "%s"
	description  = "%s"
}

resource "ovh_cloud_project_region_loadbalancer_log_subscription" "subscription" {
	service_name = "%s"
	region_name = "GRA11"
	loadbalancer_id = "%s"
	kind = "haproxy"
	stream_id = ovh_dbaas_logs_output_graylog_stream.stream.stream_id
}

data "ovh_cloud_project_region_loadbalancer_log_subscription" "test" {
	service_name = ovh_cloud_project_region_loadbalancer_log_subscription.subscription.service_name
	region_name = ovh_cloud_project_region_loadbalancer_log_subscription.subscription.region_name
	loadbalancer_id = ovh_cloud_project_region_loadbalancer_log_subscription.subscription.loadbalancer_id
	subscription_id = ovh_cloud_project_region_loadbalancer_log_subscription.subscription.id
}
`
