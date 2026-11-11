package ovh

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/ovh/go-ovh/ovh"
)

const testAccCloudStorageObjectBucketReplicationName = "ovh_cloud_storage_object_bucket_replication.replication"

func testAccPreCheckCloudStorageObjectBucketReplication(t *testing.T) {
	testAccPreCheckCloudStorageObjectBucket(t)
	checkEnvOrSkip(t, "OVH_CLOUD_PROJECT_STORAGE_REPLICATION_REGION_TEST")
}

func testAccCloudStorageObjectBucketReplicationImportStateIdFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}
		return fmt.Sprintf("%s/%s", rs.Primary.Attributes["service_name"], rs.Primary.Attributes["bucket_id"]), nil
	}
}

func testAccCheckCloudStorageObjectBucketReplicationDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ovh_cloud_storage_object_bucket_replication" {
			continue
		}

		endpoint := storageObjectBucketReplicationEndpoint(rs.Primary.Attributes["service_name"], rs.Primary.Attributes["bucket_id"])
		var res CloudStorageObjectBucketReplicationAPIRules
		if err := testAccOVHClient.Get(endpoint, &res); err != nil {
			if errOvh, ok := err.(*ovh.APIError); ok && errOvh.Code == 404 {
				continue
			}
			return fmt.Errorf("error checking replication of bucket %s was removed: %w", rs.Primary.Attributes["bucket_id"], err)
		}
		if len(res.Rules) > 0 {
			return fmt.Errorf("bucket %s still has %d replication rules", rs.Primary.Attributes["bucket_id"], len(res.Rules))
		}
	}
	return testAccCheckCloudStorageObjectBucketDestroy(s)
}

// testAccCloudStorageObjectBucketReplicationBuckets declares a versioned source bucket and a
// versioned replica bucket, both required by S3 replication.
func testAccCloudStorageObjectBucketReplicationBuckets(serviceName, sourceName, sourceRegion, replicaName, replicaRegion string) string {
	return fmt.Sprintf(`
resource "ovh_cloud_storage_object_bucket" "replica" {
  service_name = %q
  name         = %q
  region       = %q

  versioning = {
    status = "ENABLED"
  }
}

resource "ovh_cloud_storage_object_bucket" "source" {
  service_name = %q
  name         = %q
  region       = %q

  versioning = {
    status = "ENABLED"
  }
}
`, serviceName, replicaName, replicaRegion, serviceName, sourceName, sourceRegion)
}

func testAccCloudStorageObjectBucketReplicationConfig(buckets, rules string) string {
	return buckets + fmt.Sprintf(`
resource "ovh_cloud_storage_object_bucket_replication" "replication" {
  service_name = ovh_cloud_storage_object_bucket.source.service_name
  bucket_id    = ovh_cloud_storage_object_bucket.source.id

  rules = [%s]
}
`, rules)
}

func TestAccCloudStorageObjectBucketReplication_basic(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	sourceRegion := strings.ToUpper(os.Getenv("OVH_CLOUD_PROJECT_STORAGE_REGION_TEST"))
	replicaRegion := strings.ToUpper(os.Getenv("OVH_CLOUD_PROJECT_STORAGE_REPLICATION_REGION_TEST"))
	sourceName := acctest.RandomWithPrefix(test_prefix)
	replicaName := acctest.RandomWithPrefix(test_prefix)
	buckets := testAccCloudStorageObjectBucketReplicationBuckets(serviceName, sourceName, sourceRegion, replicaName, replicaRegion)

	// Rule ids are sorted: S3 may return rules sorted by id, and import has no prior order to restore.
	config := testAccCloudStorageObjectBucketReplicationConfig(buckets, `
    {
      id                        = "replicate-all"
      status                    = "ENABLED"
      priority                  = 1
      delete_marker_replication = "DISABLED"
      destination = {
        name   = ovh_cloud_storage_object_bucket.replica.name
        region = ovh_cloud_storage_object_bucket.replica.region
      }
    },
    {
      id                        = "replicate-logs"
      status                    = "ENABLED"
      priority                  = 2
      delete_marker_replication = "ENABLED"
      filter = {
        prefix = "logs/"
        tags = {
          env = "prod"
        }
      }
      destination = {
        name          = ovh_cloud_storage_object_bucket.replica.name
        region        = ovh_cloud_storage_object_bucket.replica.region
        storage_class = "STANDARD_IA"
      }
    }
`)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCloudStorageObjectBucketReplication(t)
			testAccCheckCloudProjectExists(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudStorageObjectBucketReplicationDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(testAccCloudStorageObjectBucketReplicationName, "service_name", serviceName),
					resource.TestCheckResourceAttr(testAccCloudStorageObjectBucketReplicationName, "bucket_id", sourceRegion+"_"+sourceName),
					resource.TestCheckResourceAttr(testAccCloudStorageObjectBucketReplicationName, "id", sourceRegion+"_"+sourceName),
					resource.TestCheckResourceAttr(testAccCloudStorageObjectBucketReplicationName, "rules.#", "2"),
					resource.TestCheckResourceAttr(testAccCloudStorageObjectBucketReplicationName, "rules.0.id", "replicate-all"),
					resource.TestCheckResourceAttr(testAccCloudStorageObjectBucketReplicationName, "rules.0.priority", "1"),
					resource.TestCheckResourceAttr(testAccCloudStorageObjectBucketReplicationName, "rules.0.delete_marker_replication", "DISABLED"),
					resource.TestCheckResourceAttr(testAccCloudStorageObjectBucketReplicationName, "rules.0.destination.name", replicaName),
					resource.TestCheckResourceAttr(testAccCloudStorageObjectBucketReplicationName, "rules.0.destination.region", replicaRegion),
					resource.TestCheckNoResourceAttr(testAccCloudStorageObjectBucketReplicationName, "rules.0.filter.%"),
					resource.TestCheckResourceAttr(testAccCloudStorageObjectBucketReplicationName, "rules.1.id", "replicate-logs"),
					resource.TestCheckResourceAttr(testAccCloudStorageObjectBucketReplicationName, "rules.1.delete_marker_replication", "ENABLED"),
					resource.TestCheckResourceAttr(testAccCloudStorageObjectBucketReplicationName, "rules.1.filter.prefix", "logs/"),
					resource.TestCheckResourceAttr(testAccCloudStorageObjectBucketReplicationName, "rules.1.filter.tags.env", "prod"),
					resource.TestCheckResourceAttr(testAccCloudStorageObjectBucketReplicationName, "rules.1.destination.storage_class", "STANDARD_IA"),
				),
			},
			{
				Config: config,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
			{
				ResourceName:      testAccCloudStorageObjectBucketReplicationName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccCloudStorageObjectBucketReplicationImportStateIdFunc(testAccCloudStorageObjectBucketReplicationName),
			},
		},
	})
}

func TestAccCloudStorageObjectBucketReplication_update(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	sourceRegion := strings.ToUpper(os.Getenv("OVH_CLOUD_PROJECT_STORAGE_REGION_TEST"))
	replicaRegion := strings.ToUpper(os.Getenv("OVH_CLOUD_PROJECT_STORAGE_REPLICATION_REGION_TEST"))
	sourceName := acctest.RandomWithPrefix(test_prefix)
	replicaName := acctest.RandomWithPrefix(test_prefix)
	buckets := testAccCloudStorageObjectBucketReplicationBuckets(serviceName, sourceName, sourceRegion, replicaName, replicaRegion)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCloudStorageObjectBucketReplication(t)
			testAccCheckCloudProjectExists(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudStorageObjectBucketReplicationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudStorageObjectBucketReplicationConfig(buckets, `
    {
      id                        = "replicate-logs"
      status                    = "ENABLED"
      priority                  = 1
      delete_marker_replication = "DISABLED"
      filter = {
        prefix = "logs/"
      }
      destination = {
        name   = ovh_cloud_storage_object_bucket.replica.name
        region = ovh_cloud_storage_object_bucket.replica.region
      }
    }
`),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(testAccCloudStorageObjectBucketReplicationName, "rules.#", "1"),
					resource.TestCheckResourceAttr(testAccCloudStorageObjectBucketReplicationName, "rules.0.status", "ENABLED"),
					resource.TestCheckResourceAttr(testAccCloudStorageObjectBucketReplicationName, "rules.0.filter.prefix", "logs/"),
				),
			},
			{
				Config: testAccCloudStorageObjectBucketReplicationConfig(buckets, `
    {
      id                        = "replicate-logs"
      status                    = "DISABLED"
      priority                  = 2
      delete_marker_replication = "ENABLED"
      filter = {
        prefix = "archive/"
      }
      destination = {
        name          = ovh_cloud_storage_object_bucket.replica.name
        region        = ovh_cloud_storage_object_bucket.replica.region
        storage_class = "STANDARD_IA"
      }
    },
    {
      id                        = "replicate-tmp"
      status                    = "ENABLED"
      priority                  = 3
      delete_marker_replication = "DISABLED"
      filter = {
        prefix = "tmp/"
      }
      destination = {
        name   = ovh_cloud_storage_object_bucket.replica.name
        region = ovh_cloud_storage_object_bucket.replica.region
      }
    }
`),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(testAccCloudStorageObjectBucketReplicationName, plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(testAccCloudStorageObjectBucketReplicationName, "rules.#", "2"),
					resource.TestCheckResourceAttr(testAccCloudStorageObjectBucketReplicationName, "rules.0.status", "DISABLED"),
					resource.TestCheckResourceAttr(testAccCloudStorageObjectBucketReplicationName, "rules.0.priority", "2"),
					resource.TestCheckResourceAttr(testAccCloudStorageObjectBucketReplicationName, "rules.0.delete_marker_replication", "ENABLED"),
					resource.TestCheckResourceAttr(testAccCloudStorageObjectBucketReplicationName, "rules.0.filter.prefix", "archive/"),
					resource.TestCheckResourceAttr(testAccCloudStorageObjectBucketReplicationName, "rules.0.destination.storage_class", "STANDARD_IA"),
					resource.TestCheckResourceAttr(testAccCloudStorageObjectBucketReplicationName, "rules.1.id", "replicate-tmp"),
				),
			},
		},
	})
}
