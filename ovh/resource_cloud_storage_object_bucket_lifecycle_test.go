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

const testAccCloudStorageObjectBucketLifecycleName = "ovh_cloud_storage_object_bucket_lifecycle.lifecycle"

func testAccCloudStorageObjectBucketLifecycleImportStateIdFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}
		return fmt.Sprintf("%s/%s", rs.Primary.Attributes["service_name"], rs.Primary.Attributes["bucket_id"]), nil
	}
}

func testAccCheckCloudStorageObjectBucketLifecycleDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ovh_cloud_storage_object_bucket_lifecycle" {
			continue
		}

		endpoint := storageObjectBucketLifecycleEndpoint(rs.Primary.Attributes["service_name"], rs.Primary.Attributes["bucket_id"])
		var res CloudStorageObjectBucketLifecycleAPIRules
		if err := testAccOVHClient.Get(endpoint, &res); err != nil {
			if errOvh, ok := err.(*ovh.APIError); ok && errOvh.Code == 404 {
				continue
			}
			return fmt.Errorf("error checking lifecycle of bucket %s was removed: %w", rs.Primary.Attributes["bucket_id"], err)
		}
		if len(res.Rules) > 0 {
			return fmt.Errorf("bucket %s still has %d lifecycle rules", rs.Primary.Attributes["bucket_id"], len(res.Rules))
		}
	}
	return testAccCheckCloudStorageObjectBucketDestroy(s)
}

func testAccCloudStorageObjectBucketLifecycleConfig(serviceName, bucketName, region, rules string) string {
	return fmt.Sprintf(`
resource "ovh_cloud_storage_object_bucket" "bucket" {
  service_name = "%s"
  name         = "%s"
  region       = "%s"

  versioning = {
    status = "ENABLED"
  }
}

resource "ovh_cloud_storage_object_bucket_lifecycle" "lifecycle" {
  service_name = ovh_cloud_storage_object_bucket.bucket.service_name
  bucket_id    = ovh_cloud_storage_object_bucket.bucket.id

  rules = [%s]
}
`, serviceName, bucketName, region, rules)
}

func TestAccCloudStorageObjectBucketLifecycle_basic(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	region := strings.ToUpper(os.Getenv("OVH_CLOUD_PROJECT_STORAGE_REGION_TEST"))
	bucketName := acctest.RandomWithPrefix(test_prefix)

	// Rule ids are sorted: S3 may return rules sorted by id, and import has no prior order to restore.
	config := testAccCloudStorageObjectBucketLifecycleConfig(serviceName, bucketName, region, `
    {
      id     = "abort-multipart"
      status = "ENABLED"
      abort_incomplete_multipart_upload = {
        days_after_initiation = 7
      }
    },
    {
      id     = "archive-reports"
      status = "ENABLED"
      filter = {
        prefix = "reports/"
        tags = {
          retention = "long"
        }
      }
      transitions = [
        {
          days          = 90
          storage_class = "STANDARD_IA"
        }
      ]
    },
    {
      id     = "expire-logs"
      status = "ENABLED"
      filter = {
        prefix = "logs/"
      }
      expiration = {
        days = 30
      }
    },
    {
      id     = "noncurrent-cleanup"
      status = "ENABLED"
      noncurrent_version_expiration = {
        noncurrent_days           = 60
        newer_noncurrent_versions = 3
      }
      noncurrent_version_transitions = [
        {
          noncurrent_days = 14
          storage_class   = "STANDARD_IA"
        }
      ]
    }
`)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCloudStorageObjectBucket(t)
			testAccCheckCloudProjectExists(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudStorageObjectBucketLifecycleDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(testAccCloudStorageObjectBucketLifecycleName, "service_name", serviceName),
					resource.TestCheckResourceAttr(testAccCloudStorageObjectBucketLifecycleName, "bucket_id", region+"_"+bucketName),
					resource.TestCheckResourceAttr(testAccCloudStorageObjectBucketLifecycleName, "id", region+"_"+bucketName),
					resource.TestCheckResourceAttr(testAccCloudStorageObjectBucketLifecycleName, "rules.#", "4"),
					resource.TestCheckResourceAttr(testAccCloudStorageObjectBucketLifecycleName, "rules.0.id", "abort-multipart"),
					resource.TestCheckResourceAttr(testAccCloudStorageObjectBucketLifecycleName, "rules.0.abort_incomplete_multipart_upload.days_after_initiation", "7"),
					resource.TestCheckNoResourceAttr(testAccCloudStorageObjectBucketLifecycleName, "rules.0.filter.%"),
					resource.TestCheckResourceAttr(testAccCloudStorageObjectBucketLifecycleName, "rules.1.filter.prefix", "reports/"),
					resource.TestCheckResourceAttr(testAccCloudStorageObjectBucketLifecycleName, "rules.1.filter.tags.retention", "long"),
					resource.TestCheckResourceAttr(testAccCloudStorageObjectBucketLifecycleName, "rules.1.transitions.0.days", "90"),
					resource.TestCheckResourceAttr(testAccCloudStorageObjectBucketLifecycleName, "rules.1.transitions.0.storage_class", "STANDARD_IA"),
					resource.TestCheckResourceAttr(testAccCloudStorageObjectBucketLifecycleName, "rules.2.filter.prefix", "logs/"),
					resource.TestCheckResourceAttr(testAccCloudStorageObjectBucketLifecycleName, "rules.2.expiration.days", "30"),
					resource.TestCheckResourceAttr(testAccCloudStorageObjectBucketLifecycleName, "rules.3.noncurrent_version_expiration.noncurrent_days", "60"),
					resource.TestCheckResourceAttr(testAccCloudStorageObjectBucketLifecycleName, "rules.3.noncurrent_version_expiration.newer_noncurrent_versions", "3"),
					resource.TestCheckResourceAttr(testAccCloudStorageObjectBucketLifecycleName, "rules.3.noncurrent_version_transitions.0.noncurrent_days", "14"),
					resource.TestCheckResourceAttr(testAccCloudStorageObjectBucketLifecycleName, "rules.3.noncurrent_version_transitions.0.storage_class", "STANDARD_IA"),
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
				ResourceName:      testAccCloudStorageObjectBucketLifecycleName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccCloudStorageObjectBucketLifecycleImportStateIdFunc(testAccCloudStorageObjectBucketLifecycleName),
			},
		},
	})
}

func TestAccCloudStorageObjectBucketLifecycle_update(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	region := strings.ToUpper(os.Getenv("OVH_CLOUD_PROJECT_STORAGE_REGION_TEST"))
	bucketName := acctest.RandomWithPrefix(test_prefix)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCloudStorageObjectBucket(t)
			testAccCheckCloudProjectExists(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudStorageObjectBucketLifecycleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudStorageObjectBucketLifecycleConfig(serviceName, bucketName, region, `
    {
      id     = "expire-logs"
      status = "ENABLED"
      filter = {
        prefix = "logs/"
      }
      expiration = {
        days = 30
      }
    }
`),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(testAccCloudStorageObjectBucketLifecycleName, "rules.#", "1"),
					resource.TestCheckResourceAttr(testAccCloudStorageObjectBucketLifecycleName, "rules.0.status", "ENABLED"),
					resource.TestCheckResourceAttr(testAccCloudStorageObjectBucketLifecycleName, "rules.0.expiration.days", "30"),
				),
			},
			{
				Config: testAccCloudStorageObjectBucketLifecycleConfig(serviceName, bucketName, region, `
    {
      id     = "expire-logs"
      status = "DISABLED"
      filter = {
        prefix                   = "logs/"
        object_size_greater_than = 1024
      }
      expiration = {
        days = 60
      }
    },
    {
      id     = "expire-tmp"
      status = "ENABLED"
      filter = {
        prefix = "tmp/"
      }
      expiration = {
        date = "2030-01-01"
      }
    }
`),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(testAccCloudStorageObjectBucketLifecycleName, plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(testAccCloudStorageObjectBucketLifecycleName, "rules.#", "2"),
					resource.TestCheckResourceAttr(testAccCloudStorageObjectBucketLifecycleName, "rules.0.status", "DISABLED"),
					resource.TestCheckResourceAttr(testAccCloudStorageObjectBucketLifecycleName, "rules.0.filter.object_size_greater_than", "1024"),
					resource.TestCheckResourceAttr(testAccCloudStorageObjectBucketLifecycleName, "rules.0.expiration.days", "60"),
					resource.TestCheckResourceAttr(testAccCloudStorageObjectBucketLifecycleName, "rules.1.id", "expire-tmp"),
					resource.TestCheckResourceAttr(testAccCloudStorageObjectBucketLifecycleName, "rules.1.expiration.date", "2030-01-01"),
				),
			},
		},
	})
}
