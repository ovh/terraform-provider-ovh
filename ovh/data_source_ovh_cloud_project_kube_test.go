package ovh

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccCloudProjectKubeDataSource_basic(t *testing.T) {
	name := acctest.RandomWithPrefix(test_prefix)
	region := os.Getenv("OVH_CLOUD_PROJECT_KUBE_REGION_TEST")
	version := os.Getenv("OVH_CLOUD_PROJECT_KUBE_VERSION_TEST")
	config := fmt.Sprintf(
		testAccCloudProjectKubeDatasourceConfig,
		os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST"),
		name,
		region,
		version,
	)
	matchVersion := regexp.MustCompile(`^` + version + `\..*$`)

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
						"data.ovh_cloud_project_kube.cluster", "region", region),
					resource.TestCheckResourceAttr(
						"data.ovh_cloud_project_kube.cluster", "name", name),
					resource.TestMatchResourceAttr(
						"data.ovh_cloud_project_kube.cluster", "version", matchVersion),
				),
			},
		},
	})
}

var testAccCloudProjectKubeDatasourceConfig = `
resource "ovh_cloud_project_kube" "cluster" {
	service_name  = "%s"
    name          = "%s"
	region        = "%s"
	version = "%s"
}

data "ovh_cloud_project_kube" "cluster" {
  service_name = ovh_cloud_project_kube.cluster.service_name
  kube_id = ovh_cloud_project_kube.cluster.id
}
`
