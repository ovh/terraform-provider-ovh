package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCloudManagedKubernetesRegionsDataSource_basic(t *testing.T) {
	config := fmt.Sprintf(`
data "ovh_cloud_managed_kubernetes_regions" "test" {
  service_name = "%s"
}
`, os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST"))

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckCloudManagedKubernetes(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.ovh_cloud_managed_kubernetes_regions.test", "regions.0.name"),
					resource.TestCheckResourceAttrSet("data.ovh_cloud_managed_kubernetes_regions.test", "regions.0.available_plans.#"),
				),
			},
		},
	})
}
