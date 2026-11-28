package ovh

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

const testAccCloudStorageObjectBucketLoggingName = "ovh_cloud_storage_object_bucket_logging.logging"

func testAccCloudStorageObjectBucketLoggingImportStateIdFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}
		return fmt.Sprintf("%s/%s", rs.Primary.Attributes["service_name"], rs.Primary.Attributes["bucket_id"]), nil
	}
}

// testAccCheckCloudStorageObjectBucketLogging checks the logging of the bucket resource on the S3 API, nil wants it disabled.
func testAccCheckCloudStorageObjectBucketLogging(bucketResourceName string, want *s3LoggingEnabled) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client, target, err := testAccCloudStorageObjectBucketS3(s, bucketResourceName)
		if err != nil {
			return err
		}

		got, err := client.GetBucketLogging(context.Background(), target.Host, target.Region)
		if err != nil {
			return fmt.Errorf("error reading the logging of bucket %s: %w", target.Host, err)
		}

		switch {
		case want == nil && got != nil:
			return fmt.Errorf("bucket %s still logs to %s", target.Host, got.TargetBucket)
		case want != nil && got == nil:
			return fmt.Errorf("bucket %s does not log", target.Host)
		case want != nil && *got != *want:
			return fmt.Errorf("bucket %s logs to %+v, want %+v", target.Host, *got, *want)
		}
		return nil
	}
}

func testAccCloudStorageObjectBucketLoggingConfig(serviceName, bucketName, region, userId, logging string) string {
	return fmt.Sprintf(`
resource "ovh_cloud_storage_object_bucket" "logs" {
  service_name  = "%[1]s"
  name          = "%[2]s-logs"
  region        = "%[3]s"
  owner_user_id = "%[4]s"
}

resource "ovh_cloud_storage_object_bucket_acl" "logs" {
  service_name = ovh_cloud_storage_object_bucket.logs.service_name
  bucket_id    = ovh_cloud_storage_object_bucket.logs.id
  acl          = "log-delivery-write"
}

resource "ovh_cloud_storage_object_bucket" "bucket" {
  service_name  = "%[1]s"
  name          = "%[2]s"
  region        = "%[3]s"
  owner_user_id = "%[4]s"
}
%[5]s
`, serviceName, bucketName, region, userId, logging)
}

const testAccCloudStorageObjectBucketLoggingPrefixed = `
resource "ovh_cloud_storage_object_bucket_logging" "logging" {
  service_name  = ovh_cloud_storage_object_bucket.bucket.service_name
  bucket_id     = ovh_cloud_storage_object_bucket.bucket.id
  target_bucket = ovh_cloud_storage_object_bucket.logs.name
  target_prefix = "access/"

  depends_on = [ovh_cloud_storage_object_bucket_acl.logs]
}
`

const testAccCloudStorageObjectBucketLoggingUnprefixed = `
resource "ovh_cloud_storage_object_bucket_logging" "logging" {
  service_name  = ovh_cloud_storage_object_bucket.bucket.service_name
  bucket_id     = ovh_cloud_storage_object_bucket.bucket.id
  target_bucket = ovh_cloud_storage_object_bucket.logs.name

  depends_on = [ovh_cloud_storage_object_bucket_acl.logs]
}
`

func TestAccCloudStorageObjectBucketLogging_basic(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	region := strings.ToUpper(os.Getenv("OVH_CLOUD_PROJECT_STORAGE_REGION_TEST"))
	userId := os.Getenv("OVH_CLOUD_PROJECT_S3_OWNER_USER_ID_TEST")
	bucketName := acctest.RandomWithPrefix(test_prefix)
	logsBucketName := bucketName + "-logs"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCloudStorageObjectBucketS3(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudStorageObjectBucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudStorageObjectBucketLoggingConfig(serviceName, bucketName, region, userId, testAccCloudStorageObjectBucketLoggingPrefixed),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(testAccCloudStorageObjectBucketLoggingName, "service_name", serviceName),
					resource.TestCheckResourceAttrPair(testAccCloudStorageObjectBucketLoggingName, "bucket_id", "ovh_cloud_storage_object_bucket.bucket", "id"),
					resource.TestCheckResourceAttrPair(testAccCloudStorageObjectBucketLoggingName, "id", "ovh_cloud_storage_object_bucket.bucket", "id"),
					resource.TestCheckResourceAttr(testAccCloudStorageObjectBucketLoggingName, "target_bucket", logsBucketName),
					resource.TestCheckResourceAttr(testAccCloudStorageObjectBucketLoggingName, "target_prefix", "access/"),
					testAccCheckCloudStorageObjectBucketLogging("ovh_cloud_storage_object_bucket.bucket", &s3LoggingEnabled{TargetBucket: logsBucketName, TargetPrefix: "access/"}),
				),
			},
			{
				ResourceName:      testAccCloudStorageObjectBucketLoggingName,
				ImportState:       true,
				ImportStateIdFunc: testAccCloudStorageObjectBucketLoggingImportStateIdFunc(testAccCloudStorageObjectBucketLoggingName),
				ImportStateVerify: true,
			},
			{
				Config: testAccCloudStorageObjectBucketLoggingConfig(serviceName, bucketName, region, userId, testAccCloudStorageObjectBucketLoggingUnprefixed),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPostRefresh: []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(testAccCloudStorageObjectBucketLoggingName, "target_bucket", logsBucketName),
					resource.TestCheckResourceAttr(testAccCloudStorageObjectBucketLoggingName, "target_prefix", ""),
					testAccCheckCloudStorageObjectBucketLogging("ovh_cloud_storage_object_bucket.bucket", &s3LoggingEnabled{TargetBucket: logsBucketName}),
				),
			},
			{
				Config: testAccCloudStorageObjectBucketLoggingConfig(serviceName, bucketName, region, userId, ""),
				Check:  testAccCheckCloudStorageObjectBucketLogging("ovh_cloud_storage_object_bucket.bucket", nil),
			},
		},
	})
}
