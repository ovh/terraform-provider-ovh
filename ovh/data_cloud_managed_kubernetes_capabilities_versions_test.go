package ovh

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCloudManagedKubernetesCapabilitiesVersionsDataSource_basic(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	region := os.Getenv("OVH_CLOUD_PROJECT_KUBE_REGION_TEST")
	name := acctest.RandomWithPrefix("test-mks")

	config := testAccCloudManagedKubernetesConfig(serviceName, region, name, true) + `
data "ovh_cloud_managed_kubernetes_capabilities_versions" "test" {
  service_name          = ovh_cloud_managed_kubernetes.test.service_name
  managed_kubernetes_id = ovh_cloud_managed_kubernetes.test.id
}
`

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckCloudManagedKubernetes(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.ovh_cloud_managed_kubernetes_capabilities_versions.test", "current_version"),
				),
			},
		},
	})
}
