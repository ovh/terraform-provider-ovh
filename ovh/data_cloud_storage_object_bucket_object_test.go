package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func testAccPreCheckCloudStorageObjectBucketObjectDataSource(t *testing.T) {
	testAccPreCheckCredentials(t)
	if os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST") == "" {
		t.Skip("OVH_CLOUD_PROJECT_SERVICE_TEST must be set for acceptance tests")
	}
	if os.Getenv("OVH_CLOUD_PROJECT_S3_BUCKET_TEST") == "" {
		t.Skip("OVH_CLOUD_PROJECT_S3_BUCKET_TEST must be set for acceptance tests")
	}
	if os.Getenv("OVH_CLOUD_PROJECT_S3_OBJECT_KEY_TEST") == "" {
		t.Skip("OVH_CLOUD_PROJECT_S3_OBJECT_KEY_TEST must be set for acceptance tests")
	}
}

func TestAccCloudStorageObjectBucketObjectDataSource_basic(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	bucket := os.Getenv("OVH_CLOUD_PROJECT_S3_BUCKET_TEST")
	objectKey := os.Getenv("OVH_CLOUD_PROJECT_S3_OBJECT_KEY_TEST")

	config := fmt.Sprintf(`
data "ovh_cloud_storage_object_bucket_object" "test" {
  service_name = "%s"
  bucket_name  = "%s"
  object_key   = "%s"
}
`, serviceName, bucket, objectKey)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCloudStorageObjectBucketObjectDataSource(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.ovh_cloud_storage_object_bucket_object.test", "service_name", serviceName),
					resource.TestCheckResourceAttr("data.ovh_cloud_storage_object_bucket_object.test", "bucket_name", bucket),
					resource.TestCheckResourceAttr("data.ovh_cloud_storage_object_bucket_object.test", "object_key", objectKey),
					resource.TestCheckResourceAttr("data.ovh_cloud_storage_object_bucket_object.test", "key", objectKey),
					resource.TestCheckResourceAttrSet("data.ovh_cloud_storage_object_bucket_object.test", "size"),
					resource.TestCheckResourceAttrSet("data.ovh_cloud_storage_object_bucket_object.test", "storage_class"),
				),
			},
		},
	})
}
