package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceCloudImage_basic(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	imageID := os.Getenv("OVH_INSTANCE_IMAGE_ID_TEST")

	config := fmt.Sprintf(`
data "ovh_cloud_image" "test" {
  service_name = "%s"
  id           = "%s"
}
`, serviceName, imageID)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckCloudInstance(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.ovh_cloud_image.test", "id", imageID),
					resource.TestCheckResourceAttrSet("data.ovh_cloud_image.test", "name"),
					resource.TestCheckResourceAttrSet("data.ovh_cloud_image.test", "size"),
					resource.TestCheckResourceAttrSet("data.ovh_cloud_image.test", "status"),
				),
			},
		},
	})
}
