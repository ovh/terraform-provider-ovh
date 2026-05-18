package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

var testAccCloudProjectKubeLogSubscriptionConfig = `
resource "ovh_cloud_project_kube" "cluster" {
  service_name = "%s"
  name         = "%s"
  region       = "%s"
}

resource "ovh_dbaas_logs_output_graylog_stream" "stream" {
  service_name = "%s"
  title        = "%s"
  description  = "%s"
}

resource "ovh_cloud_project_kube_log_subscription" "sub" {
  service_name = ovh_cloud_project_kube.cluster.service_name
  kube_id      = ovh_cloud_project_kube.cluster.id
  stream_id    = ovh_dbaas_logs_output_graylog_stream.stream.stream_id
  kind         = "audit"
}
`

func TestAccCloudProjectKubeLogSubscription_basic(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	region := os.Getenv("OVH_CLOUD_PROJECT_KUBE_REGION_TEST")
	ldpServiceName := os.Getenv("OVH_DBAAS_LOGS_SERVICE_TEST")
	clusterName := acctest.RandomWithPrefix(test_prefix)
	streamTitle := acctest.RandomWithPrefix(test_prefix)
	streamDesc := acctest.RandomWithPrefix(test_prefix)

	resourceName := "ovh_cloud_project_kube_log_subscription.sub"

	config := fmt.Sprintf(
		testAccCloudProjectKubeLogSubscriptionConfig,
		serviceName,
		clusterName,
		region,
		ldpServiceName,
		streamTitle,
		streamDesc,
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCloud(t)
			testAccCheckCloudProjectExists(t)
			testAccPreCheckKubernetes(t)
			testAccPreCheckDbaasLogs(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			// Step 1: Creation
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "subscription_id"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
					resource.TestCheckResourceAttrSet(resourceName, "updated_at"),
					resource.TestCheckResourceAttr(resourceName, "kind", "audit"),
					resource.TestCheckResourceAttrSet(resourceName, "stream_id"),
					resource.TestCheckResourceAttrSet(resourceName, "resource.0.name"),
					resource.TestCheckResourceAttrSet(resourceName, "resource.0.type"),
				),
			},
			// Step 2: No re-creation (idempotency)
			{
				Config: config,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPostRefresh: []plancheck.PlanCheck{
						ExpectEmptyPlan(),
					},
				},
			},
			// Step 3: Import
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: func(state *terraform.State) (string, error) {
					res, ok := state.RootModule().Resources[resourceName]
					if !ok {
						return "", fmt.Errorf("resource %s not found in state", resourceName)
					}
					return fmt.Sprintf(
						"%s/%s/%s",
						res.Primary.Attributes["service_name"],
						res.Primary.Attributes["kube_id"],
						res.Primary.ID,
					), nil
				},
			},
		},
	})
}
