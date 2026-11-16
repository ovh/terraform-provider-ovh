package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceCloudStorageObjectBucketMultipartUploads_basic(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	bucketId := testAccCloudStorageObjectBucketTestId()
	key := os.Getenv("OVH_CLOUD_PROJECT_STORAGE_OBJECT_TEST")

	config := fmt.Sprintf(`
data "ovh_cloud_storage_object_bucket_multipart_uploads" "all" {
  service_name = %q
  bucket_id    = %q
}

data "ovh_cloud_storage_object_bucket_multipart_uploads" "by_key" {
  service_name = %q
  bucket_id    = %q
  key          = %q
}
`, serviceName, bucketId, serviceName, bucketId, key)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckCloudStorage(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.ovh_cloud_storage_object_bucket_multipart_uploads.all", "uploads.#"),
					resource.TestCheckResourceAttr("data.ovh_cloud_storage_object_bucket_multipart_uploads.by_key", "key", key),
					resource.TestCheckResourceAttrSet("data.ovh_cloud_storage_object_bucket_multipart_uploads.by_key", "uploads.#"),
				),
			},
		},
	})
}
