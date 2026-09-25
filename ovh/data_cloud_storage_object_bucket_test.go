package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceCloudStorageObjectBucket_basic(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	region := os.Getenv("OVH_CLOUD_PROJECT_STORAGE_REGION_TEST")

	bucketName := acctest.RandomWithPrefix(test_prefix)

	config := fmt.Sprintf(`
resource "ovh_cloud_s3_bucket" "bucket" {
  service_name = "%s"
  name         = "%s"
  region       = "%s"

  versioning = {
    status = "ENABLED"
  }

  tags = {
    env = "test"
  }
}

data "ovh_cloud_storage_object_bucket" "by_id" {
  service_name = "%s"
  id           = ovh_cloud_s3_bucket.bucket.id
}
`, serviceName, bucketName, region, serviceName)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCloudS3Bucket(t)
			testAccCheckCloudProjectExists(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.ovh_cloud_storage_object_bucket.by_id", "service_name", serviceName),
					resource.TestCheckResourceAttrPair(
						"data.ovh_cloud_storage_object_bucket.by_id", "id",
						"ovh_cloud_s3_bucket.bucket", "id",
					),
					resource.TestCheckResourceAttr("data.ovh_cloud_storage_object_bucket.by_id", "id", region+"_"+bucketName),
					resource.TestCheckResourceAttr("data.ovh_cloud_storage_object_bucket.by_id", "name", bucketName),
					resource.TestCheckResourceAttr("data.ovh_cloud_storage_object_bucket.by_id", "location.region", region),
					resource.TestCheckResourceAttr("data.ovh_cloud_storage_object_bucket.by_id", "versioning.status", "ENABLED"),
					resource.TestCheckResourceAttr("data.ovh_cloud_storage_object_bucket.by_id", "tags.env", "test"),
					resource.TestCheckResourceAttrSet("data.ovh_cloud_storage_object_bucket.by_id", "checksum"),
					resource.TestCheckResourceAttrSet("data.ovh_cloud_storage_object_bucket.by_id", "created_at"),
					resource.TestCheckResourceAttrSet("data.ovh_cloud_storage_object_bucket.by_id", "resource_status"),
					resource.TestCheckResourceAttr("data.ovh_cloud_storage_object_bucket.by_id", "current_state.name", bucketName),
					resource.TestCheckResourceAttr("data.ovh_cloud_storage_object_bucket.by_id", "current_state.location.region", region),
					resource.TestCheckResourceAttrSet("data.ovh_cloud_storage_object_bucket.by_id", "current_state.virtual_host"),
					resource.TestCheckResourceAttrSet("data.ovh_cloud_storage_object_bucket.by_id", "current_state.objects_count"),
					resource.TestCheckResourceAttrSet("data.ovh_cloud_storage_object_bucket.by_id", "current_state.objects_size"),
				),
			},
		},
	})
}
