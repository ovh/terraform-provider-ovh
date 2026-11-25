package ovh

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// Checks that the environment variables needed for the bucket resources going through the S3 API are set:
// the S3 credentials the provider uses, and the id of their user, owner of the test buckets.
func testAccPreCheckCloudStorageObjectBucketS3(t *testing.T) {
	testAccPreCheckCloudStorageObjectBucket(t)
	checkEnvOrSkip(t, "OVH_S3_ACCESS_KEY_ID")
	checkEnvOrSkip(t, "OVH_S3_SECRET_ACCESS_KEY")
	checkEnvOrSkip(t, "OVH_CLOUD_PROJECT_S3_OWNER_USER_ID_TEST")
}

// testAccCloudStorageObjectBucketS3 locates the bucket of the state resource on the S3 API,
// with a client built from the test environment.
func testAccCloudStorageObjectBucketS3(s *terraform.State, resourceName string) (*s3Client, *storageObjectBucketS3Target, error) {
	rs, ok := s.RootModule().Resources[resourceName]
	if !ok {
		return nil, nil, fmt.Errorf("not found: %s", resourceName)
	}

	target, err := storageObjectBucketS3TargetRead(context.Background(), testAccOVHClient, rs.Primary.Attributes["service_name"], rs.Primary.Attributes["id"])
	if err != nil {
		return nil, nil, fmt.Errorf("error locating bucket %s on the S3 API: %w", rs.Primary.Attributes["id"], err)
	}

	client := newS3Client(os.Getenv("OVH_S3_ACCESS_KEY_ID"), os.Getenv("OVH_S3_SECRET_ACCESS_KEY"), "terraform-provider-ovh/testacc")
	return client, target, nil
}
