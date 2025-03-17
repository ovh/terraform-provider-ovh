package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceCloudProjectImage_basic(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	imageID := os.Getenv("OVH_CLOUD_PROJECT_IMAGE_TEST")
	config := fmt.Sprintf(`
		data "ovh_cloud_project_image" "image" {
			service_name = "%s"
			image_id     = "%s"
		}
	`, serviceName, imageID)

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
					resource.TestCheckResourceAttr("data.ovh_cloud_project_image.image", "service_name", serviceName),
					resource.TestCheckResourceAttr("data.ovh_cloud_project_image.image", "id", imageID),
					resource.TestCheckResourceAttrSet("data.ovh_cloud_project_image.image", "creation_date"),
					resource.TestCheckResourceAttrSet("data.ovh_cloud_project_image.image", "name"),
					resource.TestCheckResourceAttrSet("data.ovh_cloud_project_image.image", "region"),
					resource.TestCheckResourceAttrSet("data.ovh_cloud_project_image.image", "size"),
					resource.TestCheckResourceAttrSet("data.ovh_cloud_project_image.image", "type"),
				),
			},
		},
	})
}
