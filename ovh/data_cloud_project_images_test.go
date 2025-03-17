package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceCloudProjectImages_basic(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	config := fmt.Sprintf(`
		data "ovh_cloud_project_images" "images" {
			service_name = "%s"
		}
	`, serviceName)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCloud(t)
			testAccCheckCloudProjectExists(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.ovh_cloud_project_images.images", "service_name", serviceName),
					resource.TestCheckResourceAttrSet("data.ovh_cloud_project_images.images", "images.#"),
					resource.TestCheckResourceAttrSet("data.ovh_cloud_project_images.images", "images.0.id"),
					resource.TestCheckResourceAttrSet("data.ovh_cloud_project_images.images", "images.0.creation_date"),
					resource.TestCheckResourceAttrSet("data.ovh_cloud_project_images.images", "images.0.name"),
					resource.TestCheckResourceAttrSet("data.ovh_cloud_project_images.images", "images.0.region"),
					resource.TestCheckResourceAttrSet("data.ovh_cloud_project_images.images", "images.0.size"),
					resource.TestCheckResourceAttrSet("data.ovh_cloud_project_images.images", "images.0.type"),
				),
			},
		},
	})
}
