package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCloudProjectLoadBalancerGetLogSubscriptions_basic(t *testing.T) {

	config := fmt.Sprintf(testAccCloudProjectSubscriptions,
		os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST"),
		os.Getenv("OVH_CLOUD_PROJECT_WORKFLOW_BACKUP_REGION_TEST"),
		os.Getenv("OVH_CLOUD_LOADBALANCER_ID_TEST"),
		os.Getenv("OVH_CLOUD_STREAM_ID_TEST"),
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
					resource.TestCheckResourceAttrSet(
						"data.ovh_cloud_project_region_loadbalancer_log_subscriptions.test", "subscription_ids.#"),
				),
			},
		},
	})
}

func testAccCheckGetSubscriptions(t *testing.T) {
	testAccPreCheckCloud(t)
	testAccCheckCloudProjectExists(t)
}

var testAccCloudProjectSubscriptions = `
resource "ovh_cloud_project_region_loadbalancer_log_subscription" "subscription" {
	service_name = "%s"
	region_name = "%s"
	loadbalancer_id = "%s"
	kind = "haproxy"
	stream_id = "%s"
}

data "ovh_cloud_project_region_loadbalancer_log_subscriptions" "test" {
	service_name = ovh_cloud_project_region_loadbalancer_log_subscription.subscription.service_name
    region_name = ovh_cloud_project_region_loadbalancer_log_subscription.subscription.region_name
    loadbalancer_id = ovh_cloud_project_region_loadbalancer_log_subscription.subscription.loadbalancer_id
}
`
