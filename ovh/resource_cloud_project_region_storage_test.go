package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCloudProjectRegionStorage_basic(t *testing.T) {
	config := fmt.Sprintf(`
	resource "ovh_cloud_project_region_storage" "storage" {
		service_name = "%s"
		region_name = "GRA"
		name = "storage-test"
		versioning = {
			status = "enabled"
		}
	}
	`, os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST"))

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckCloud(t); testAccCheckCloudProjectExists(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_project_region_storage.storage", "name", "storage-test"),
					resource.TestCheckResourceAttr("ovh_cloud_project_region_storage.storage", "region", "GRA"),
					resource.TestCheckResourceAttr("ovh_cloud_project_region_storage.storage", "versioning.status", "enabled"),
					resource.TestCheckResourceAttr("ovh_cloud_project_region_storage.storage", "encryption.sse_algorithm", "plaintext"),
					resource.TestCheckResourceAttrSet("ovh_cloud_project_region_storage.storage", "virtual_host"),
				),
			},
			{
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "name",
				ResourceName:                         "ovh_cloud_project_region_storage.storage",
				ImportStateId:                        fmt.Sprintf("%s/GRA/storage-test", os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")),
				ImportStateVerifyIgnore:              []string{"created_at"}, // Ignore created_at since its value is invalid in response of the POST.
			},
		},
	})
}
