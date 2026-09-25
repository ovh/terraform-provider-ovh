package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceCloudS3Buckets_basic(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	region := os.Getenv("OVH_CLOUD_PROJECT_STORAGE_REGION_TEST")

	bucketName := acctest.RandomWithPrefix(test_prefix)

	config := fmt.Sprintf(`
resource "ovh_cloud_s3_bucket" "bucket" {
  service_name = "%s"
  name         = "%s"
  region       = "%s"
}

data "ovh_cloud_s3_buckets" "buckets" {
  service_name = "%s"
  region       = "%s"

  depends_on = [ovh_cloud_s3_bucket.bucket]
}
`, serviceName, bucketName, region, serviceName, region)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCloudS3Bucket(t)
			testAccCheckCloudProjectExists(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudS3BucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.ovh_cloud_s3_buckets.buckets", "service_name", serviceName),
					resource.TestCheckResourceAttr("data.ovh_cloud_s3_buckets.buckets", "region", region),
					resource.TestCheckResourceAttrSet("data.ovh_cloud_s3_buckets.buckets", "buckets.0.id"),
					resource.TestCheckResourceAttrSet("data.ovh_cloud_s3_buckets.buckets", "buckets.0.name"),
					resource.TestCheckResourceAttrSet("data.ovh_cloud_s3_buckets.buckets", "buckets.0.checksum"),
					resource.TestCheckResourceAttrSet("data.ovh_cloud_s3_buckets.buckets", "buckets.0.resource_status"),
					resource.TestCheckResourceAttr("data.ovh_cloud_s3_buckets.buckets", "buckets.0.location.region", region),
					resource.TestCheckResourceAttr("data.ovh_cloud_s3_buckets.buckets", "buckets.0.current_state.location.region", region),
				),
			},
		},
	})
}
