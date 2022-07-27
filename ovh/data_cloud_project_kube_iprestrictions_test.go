package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccCloudProjectKubeIPRestrictionsDataSource_basic(t *testing.T) {
	name := acctest.RandomWithPrefix(test_prefix)
	region := os.Getenv("OVH_CLOUD_PROJECT_KUBE_REGION_TEST")

	config := fmt.Sprintf(
		testAccCloudProjectKubeIPRestrictionsDataSourceConfig,
		os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST"),
		name,
		region,
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckKubernetes(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.ovh_cloud_project_kube_iprestrictions.iprestrictionsData", "service_name", os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")),
					resource.TestCheckResourceAttr(
						"data.ovh_cloud_project_kube_iprestrictions.iprestrictionsData", "ips.0", "10.42.0.0/16"),
				),
			},
		},
	})
}

var testAccCloudProjectKubeIPRestrictionsDataSourceConfig = `
resource "ovh_cloud_project_kube" "cluster" {
	service_name  = "%s"
	name          = "%s"
	region        = "%s"
}

resource "ovh_cloud_project_kube_iprestrictions" "iprestrictions" {
	service_name  = ovh_cloud_project_kube.cluster.service_name
	kube_id       = ovh_cloud_project_kube.cluster.id
	ips           = ["10.42.0.0/16"]

	depends_on = [
		ovh_cloud_project_kube.cluster
	]

}

data "ovh_cloud_project_kube_iprestrictions" "iprestrictionsData" {
  service_name = ovh_cloud_project_kube.cluster.service_name
  kube_id = ovh_cloud_project_kube.cluster.id

	depends_on = [
		ovh_cloud_project_kube_iprestrictions.iprestrictions
	]
}
`
