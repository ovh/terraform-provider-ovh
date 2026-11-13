package ovh

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

const testAccCloudStorageObjectBucketReplicationJobName = "ovh_cloud_storage_object_bucket_replication_job.job"

func testAccCloudStorageObjectBucketReplicationJobImportStateIdFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}
		return fmt.Sprintf("%s/%s/%s", rs.Primary.Attributes["service_name"], rs.Primary.Attributes["bucket_id"], rs.Primary.Attributes["id"]), nil
	}
}

// testAccCloudStorageObjectBucketReplicationJobConfig declares a source bucket replicated to a
// replica bucket and one batch replication job on the source bucket.
func testAccCloudStorageObjectBucketReplicationJobConfig(serviceName, sourceName, sourceRegion, replicaName, replicaRegion string) string {
	buckets := testAccCloudStorageObjectBucketReplicationBuckets(serviceName, sourceName, sourceRegion, replicaName, replicaRegion)
	return testAccCloudStorageObjectBucketReplicationConfig(buckets, `
    {
      id                        = "replicate-all"
      status                    = "ENABLED"
      priority                  = 1
      delete_marker_replication = "DISABLED"
      destination = {
        name   = ovh_cloud_storage_object_bucket.replica.name
        region = ovh_cloud_storage_object_bucket.replica.region
      }
    }
`) + `
resource "ovh_cloud_storage_object_bucket_replication_job" "job" {
  service_name = ovh_cloud_storage_object_bucket_replication.replication.service_name
  bucket_id    = ovh_cloud_storage_object_bucket_replication.replication.bucket_id
}
`
}

func TestAccCloudStorageObjectBucketReplicationJob_basic(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	sourceRegion := strings.ToUpper(os.Getenv("OVH_CLOUD_PROJECT_STORAGE_REGION_TEST"))
	replicaRegion := strings.ToUpper(os.Getenv("OVH_CLOUD_PROJECT_STORAGE_REPLICATION_REGION_TEST"))
	sourceName := acctest.RandomWithPrefix(test_prefix)
	replicaName := acctest.RandomWithPrefix(test_prefix)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCloudStorageObjectBucketReplication(t)
			testAccCheckCloudProjectExists(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudStorageObjectBucketReplicationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudStorageObjectBucketReplicationJobConfig(serviceName, sourceName, sourceRegion, replicaName, replicaRegion),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(testAccCloudStorageObjectBucketReplicationJobName, "service_name", serviceName),
					resource.TestCheckResourceAttr(testAccCloudStorageObjectBucketReplicationJobName, "bucket_id", sourceRegion+"_"+sourceName),
					resource.TestCheckResourceAttrSet(testAccCloudStorageObjectBucketReplicationJobName, "id"),
					resource.TestCheckResourceAttrSet(testAccCloudStorageObjectBucketReplicationJobName, "status"),
					resource.TestCheckResourceAttrSet(testAccCloudStorageObjectBucketReplicationJobName, "total"),
					resource.TestCheckResourceAttrSet(testAccCloudStorageObjectBucketReplicationJobName, "created_at"),
				),
			},
			{
				ResourceName:      testAccCloudStorageObjectBucketReplicationJobName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccCloudStorageObjectBucketReplicationJobImportStateIdFunc(testAccCloudStorageObjectBucketReplicationJobName),
				// The job progresses between the apply and the import.
				ImportStateVerifyIgnore: []string{"status", "total", "completed", "failed", "stuck", "created_at"},
			},
		},
	})
}
