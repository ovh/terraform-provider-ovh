package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCloudManagedKubernetesVersionsDataSource_basic(t *testing.T) {
	config := fmt.Sprintf(`
data "ovh_cloud_managed_kubernetes_versions" "test" {
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
					resource.TestCheckResourceAttrSet("data.ovh_cloud_managed_kubernetes_versions.test", "default_version"),
					resource.TestCheckResourceAttrSet("data.ovh_cloud_managed_kubernetes_versions.test", "available_versions.0.version"),
				),
			},
		},
	})
}
