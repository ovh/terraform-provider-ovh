package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func testAccPreCheckCloudStorageObjectBucketObjectsDataSource(t *testing.T) {
	testAccPreCheckCredentials(t)
	if os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST") == "" {
		t.Skip("OVH_CLOUD_PROJECT_SERVICE_TEST must be set for acceptance tests")
	}
	if os.Getenv("OVH_CLOUD_PROJECT_S3_BUCKET_TEST") == "" {
		t.Skip("OVH_CLOUD_PROJECT_S3_BUCKET_TEST must be set for acceptance tests")
	}
}

func TestAccCloudStorageObjectBucketObjectsDataSource_basic(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	bucket := os.Getenv("OVH_CLOUD_PROJECT_S3_BUCKET_TEST")

	config := fmt.Sprintf(`
data "ovh_cloud_storage_object_bucket_objects" "test" {
  service_name = "%s"
  bucket_name  = "%s"
}
`, serviceName, bucket)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCloudStorageObjectBucketObjectsDataSource(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.ovh_cloud_storage_object_bucket_objects.test", "service_name", serviceName),
					resource.TestCheckResourceAttr("data.ovh_cloud_storage_object_bucket_objects.test", "bucket_name", bucket),
					resource.TestCheckResourceAttrSet("data.ovh_cloud_storage_object_bucket_objects.test", "is_truncated"),
				),
			},
		},
	})
}
