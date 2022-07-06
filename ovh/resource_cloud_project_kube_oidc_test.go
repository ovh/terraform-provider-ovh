package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

var testAccCloudProjectKubeOIDCConfig = `
	resource "ovh_cloud_project_kube" "cluster" {
		service_name  = "%s"
		name          = "%s"
		region        = "%s"
	}
	resource "ovh_cloud_project_kube_oidc" "my-oidc" {
		service_name  = ovh_cloud_project_kube.cluster.service_name
		kube_id       = ovh_cloud_project_kube.cluster.id
  		client_id    = "%s"
  		issuer_url   = "%s"
	}
`

func TestAccCloudProjectKubeOIDC_full(t *testing.T) {
	region := os.Getenv("OVH_CLOUD_PROJECT_KUBE_REGION_TEST")

	name := acctest.RandomWithPrefix(test_prefix)
	config := fmt.Sprintf(
		testAccCloudProjectKubeOIDCConfig,
		os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST"),
		name,
		region,
		"xxx",
		"https://ovh.com",
	)

	configUpdated := fmt.Sprintf(
		testAccCloudProjectKubeOIDCConfig,
		os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST"),
		name,
		region,
		"yyy",
		"https://docs.ovh.com",
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
						"ovh_cloud_project_kube_oidc.my-oidc", "client_id", "xxx"),
					resource.TestCheckResourceAttr(
						"ovh_cloud_project_kube_oidc.my-oidc", "issuer_url", "https://ovh.com"),
				),
			},
			{
				Config: configUpdated,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"ovh_cloud_project_kube_oidc.my-oidc", "client_id", "yyy"),
					resource.TestCheckResourceAttr(
						"ovh_cloud_project_kube_oidc.my-oidc", "issuer_url", "https://docs.ovh.com"),
				),
			},
			{
				Config:       configUpdated,
				Destroy:      true,
				ResourceName: "ovh_cloud_project_kube_oidc.my-oidc",
			},
		},
	})
}
