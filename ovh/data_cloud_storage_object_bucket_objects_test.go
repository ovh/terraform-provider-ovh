package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceCloudStorageObjectBucketObjects_basic(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	bucketId := testAccCloudStorageObjectBucketTestId()
	key := os.Getenv("OVH_CLOUD_PROJECT_STORAGE_OBJECT_TEST")

	config := fmt.Sprintf(`
data "ovh_cloud_storage_object_bucket_objects" "all" {
  service_name = %q
  bucket_id    = %q
}

data "ovh_cloud_storage_object_bucket_objects" "by_key" {
  service_name = %q
  bucket_id    = %q
  prefix       = %q
}

data "ovh_cloud_storage_object_bucket_objects" "versions" {
  service_name  = %q
  bucket_id     = %q
  prefix        = %q
  with_versions = true
}
`, serviceName, bucketId, serviceName, bucketId, key, serviceName, bucketId, key)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckCloudStorage(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.ovh_cloud_storage_object_bucket_objects.all", "objects.#"),
					resource.TestCheckTypeSetElemNestedAttrs("data.ovh_cloud_storage_object_bucket_objects.all", "objects.*", map[string]string{
						"key":              key,
						"is_common_prefix": "false",
					}),
					resource.TestCheckTypeSetElemNestedAttrs("data.ovh_cloud_storage_object_bucket_objects.by_key", "objects.*", map[string]string{
						"key": key,
					}),
					resource.TestCheckResourceAttrSet("data.ovh_cloud_storage_object_bucket_objects.by_key", "objects.0.size"),
					resource.TestCheckResourceAttrSet("data.ovh_cloud_storage_object_bucket_objects.by_key", "objects.0.etag"),
					resource.TestCheckTypeSetElemNestedAttrs("data.ovh_cloud_storage_object_bucket_objects.versions", "objects.*", map[string]string{
						"key":       key,
						"is_latest": "true",
					}),
				),
			},
		},
	})
}
