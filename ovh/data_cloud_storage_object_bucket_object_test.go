package ovh

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// testAccCloudStorageObjectBucketTestId builds the <REGION>_<name> id of the pre-existing test bucket.
func testAccCloudStorageObjectBucketTestId() string {
	return strings.ToUpper(os.Getenv("OVH_CLOUD_PROJECT_STORAGE_REGION_TEST")) + "_" + os.Getenv("OVH_CLOUD_PROJECT_STORAGE_BUCKET_NAME_TEST")
}

func TestAccDataSourceCloudStorageObjectBucketObject_basic(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	bucketId := testAccCloudStorageObjectBucketTestId()
	key := os.Getenv("OVH_CLOUD_PROJECT_STORAGE_OBJECT_TEST")

	config := fmt.Sprintf(`
data "ovh_cloud_storage_object_bucket_object" "object" {
  service_name = %q
  bucket_id    = %q
  key          = %q
}
`, serviceName, bucketId, key)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckCloudStorage(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.ovh_cloud_storage_object_bucket_object.object", "bucket_id", bucketId),
					resource.TestCheckResourceAttr("data.ovh_cloud_storage_object_bucket_object.object", "key", key),
					resource.TestCheckResourceAttrSet("data.ovh_cloud_storage_object_bucket_object.object", "size"),
					resource.TestCheckResourceAttrSet("data.ovh_cloud_storage_object_bucket_object.object", "etag"),
					resource.TestCheckResourceAttrSet("data.ovh_cloud_storage_object_bucket_object.object", "last_modified_at"),
					resource.TestCheckResourceAttrSet("data.ovh_cloud_storage_object_bucket_object.object", "storage_class"),
					resource.TestCheckResourceAttr("data.ovh_cloud_storage_object_bucket_object.object", "is_delete_marker", "false"),
				),
			},
		},
	})
}
