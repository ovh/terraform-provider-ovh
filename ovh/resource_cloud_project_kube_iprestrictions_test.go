package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

var testAccCloudProjectKubeIpRestrictionsConfig = `
resource "ovh_cloud_project_kube" "cluster" {
	service_name  = "%s"
  name          = "%s"
	region        = "%s"
	version       = "%s"
}

resource "ovh_cloud_project_kube_iprestrictions" "iprestrictions" {
	service_name  = ovh_cloud_project_kube.cluster.service_name
	kube_id       = ovh_cloud_project_kube.cluster.id
	ips           = toset(["10.42.0.0/16"])
}
`

func TestAccCloudProjectKubeIpRestrictions_basic(t *testing.T) {
	name := acctest.RandomWithPrefix(test_prefix)
	region := os.Getenv("OVH_CLOUD_PROJECT_KUBE_REGION_TEST")
	version := os.Getenv("OVH_CLOUD_PROJECT_KUBE_VERSION_TEST")
	config := fmt.Sprintf(
		testAccCloudProjectKubeIpRestrictionsConfig,
		os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST"),
		name,
		region,
		version,
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
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"ovh_cloud_project_kube.cluster", "region", region),
					resource.TestCheckResourceAttrSet(
						"ovh_cloud_project_kube.cluster", "kubeconfig"),
					resource.TestCheckResourceAttr(
						"ovh_cloud_project_kube.cluster", "name", name),
					resource.TestCheckResourceAttr(
						"ovh_cloud_project_kube.cluster", "version", version),
					resource.TestCheckResourceAttr(
						"ovh_cloud_project_kube_iprestrictions.iprestrictions", "ips", "[10.42.0.0/16]"),
				),
			},
		},
	})
}
