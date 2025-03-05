package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCloudProjectStorageObjectDataSource_basic(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	regionName := os.Getenv("OVH_CLOUD_PROJECT_STORAGE_REGION_TEST")
	bucketName := os.Getenv("OVH_CLOUD_PROJECT_STORAGE_BUCKET_NAME_TEST")
	objectName := os.Getenv("OVH_CLOUD_PROJECT_STORAGE_OBJECT_TEST")

	config := fmt.Sprintf(`
		data "ovh_cloud_project_storage_object" "obj" {
			service_name = %q
			region_name  = %q
			name         = %q
			key          = %q
		}
	`, serviceName, regionName, bucketName, objectName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckCloudStorage(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.ovh_cloud_project_storage_object.obj", "key", objectName),
					resource.TestCheckResourceAttrSet("data.ovh_cloud_project_storage_object.obj", "etag"),
					resource.TestCheckResourceAttrSet("data.ovh_cloud_project_storage_object.obj", "size"),
					resource.TestCheckResourceAttrSet("data.ovh_cloud_project_storage_object.obj", "storage_class"),
				),
			},
		},
	})
}
