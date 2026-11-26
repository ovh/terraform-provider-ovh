package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCloudManagedKubernetesPlansDataSource_basic(t *testing.T) {
	config := fmt.Sprintf(`
data "ovh_cloud_managed_kubernetes_plans" "test" {
  service_name = "%s"
  region       = "%s"
}
`, os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST"), os.Getenv("OVH_CLOUD_PROJECT_KUBE_REGION_TEST"))

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckCloudManagedKubernetes(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.ovh_cloud_managed_kubernetes_plans.test", "plans.0.name"),
					resource.TestCheckResourceAttrSet("data.ovh_cloud_managed_kubernetes_plans.test", "plans.0.limits.max_nodes"),
					resource.TestCheckResourceAttrSet("data.ovh_cloud_managed_kubernetes_plans.test", "plans.0.control_plane.etcd.storage.limit"),
				),
			},
		},
	})
}
