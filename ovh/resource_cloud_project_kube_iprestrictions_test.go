package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

var testAccCloudProjectKubeIpRestrictionsConfig = `
resource "ovh_cloud_project_kube" "cluster" {
	service_name  = "%s"
	name          = "%s"
	region        = "%s"
}

resource "ovh_cloud_project_kube_iprestrictions" "iprestrictions" {
	service_name  = ovh_cloud_project_kube.cluster.service_name
	kube_id       = ovh_cloud_project_kube.cluster.id
	ips           = %s
}
`

func TestAccCloudProjectKubeIpRestrictions_basic(t *testing.T) {
	name := acctest.RandomWithPrefix(test_prefix)
	region := os.Getenv("OVH_CLOUD_PROJECT_KUBE_REGION_TEST")
	resourceName := "ovh_cloud_project_kube_iprestrictions.iprestrictions"

	ips1 := `["10.42.0.0/16"]`
	config1 := fmt.Sprintf(
		testAccCloudProjectKubeIpRestrictionsConfig,
		os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST"),
		name,
		region,
		ips1,
	)

	ips2 := `["10.42.0.0/16","10.43.0.0/16"]`
	config2 := fmt.Sprintf(
		testAccCloudProjectKubeIpRestrictionsConfig,
		os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST"),
		name,
		region,
		ips2,
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCloud(t)
			testAccCheckCloudProjectExists(t)
			testAccPreCheckKubernetes(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config1,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", "region", region),
					resource.TestCheckResourceAttrSet("ovh_cloud_project_kube.cluster", "kubeconfig"),
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", "name", name),
					resource.TestCheckResourceAttr("ovh_cloud_project_kube_iprestrictions.iprestrictions", "ips.0", "10.42.0.0/16"),
				),
			},
			{
				Config: config2,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", "region", region),
					resource.TestCheckResourceAttrSet("ovh_cloud_project_kube.cluster", "kubeconfig"),
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", "name", name),
					resource.TestCheckResourceAttr("ovh_cloud_project_kube_iprestrictions.iprestrictions", "ips.1", "10.43.0.0/16"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: func(state *terraform.State) (string, error) {
					return fmt.Sprintf("%s/%s", os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST"), state.RootModule().Resources[resourceName].Primary.ID), nil
				},
			},
		},
	})
}
