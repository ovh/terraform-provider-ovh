package ovh

import (
	"fmt"
	"net/url"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func testAccCloudStorageObjectBucketMultipartUploadCreate(t *testing.T, serviceName, bucketId, key string) string {
	bucketPath := "/v2/publicCloud/project/" + url.PathEscape(serviceName) +
		"/storage/object/bucket/" + url.PathEscape(bucketId)

	var upload CloudStorageObjectBucketMultipartUploadAPIResponse
	if err := testAccOVHClient.Post(bucketPath+"/multipartUpload", map[string]string{"key": key}, &upload); err != nil {
		t.Fatalf("creating multipart upload for key %q: %s", key, err)
	}

	t.Cleanup(func() {
		abortPath := bucketPath + "/object/" + url.PathEscape(key) + "/multipartUpload/" + url.PathEscape(upload.Id)
		if err := testAccOVHClient.Delete(abortPath, nil); err != nil {
			t.Errorf("aborting multipart upload %q: %s", upload.Id, err)
		}
	})

	return upload.Id
}

func TestAccDataSourceCloudStorageObjectBucketMultipartUploads_basic(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	bucketId := testAccCloudStorageObjectBucketTestId()
	key := acctest.RandomWithPrefix(test_prefix) + ".bin"
	var uploadId string

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
		PreCheck: func() {
			testAccPreCheckCloudStorage(t)
			uploadId = testAccCloudStorageObjectBucketMultipartUploadCreate(t, serviceName, bucketId, key)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeAggregateTestCheckFunc(
					func(s *terraform.State) error {
						return resource.TestCheckTypeSetElemNestedAttrs("data.ovh_cloud_storage_object_bucket_multipart_uploads.all", "uploads.*", map[string]string{
							"id":  uploadId,
							"key": key,
						})(s)
					},
					resource.TestCheckResourceAttr("data.ovh_cloud_storage_object_bucket_multipart_uploads.by_key", "key", key),
					resource.TestCheckResourceAttr("data.ovh_cloud_storage_object_bucket_multipart_uploads.by_key", "uploads.#", "1"),
					func(s *terraform.State) error {
						return resource.TestCheckResourceAttr("data.ovh_cloud_storage_object_bucket_multipart_uploads.by_key", "uploads.0.id", uploadId)(s)
					},
					resource.TestCheckResourceAttr("data.ovh_cloud_storage_object_bucket_multipart_uploads.by_key", "uploads.0.key", key),
					resource.TestCheckResourceAttrSet("data.ovh_cloud_storage_object_bucket_multipart_uploads.by_key", "uploads.0.storage_class"),
					resource.TestCheckResourceAttrSet("data.ovh_cloud_storage_object_bucket_multipart_uploads.by_key", "uploads.0.created_at"),
				),
			},
		},
	})
}
