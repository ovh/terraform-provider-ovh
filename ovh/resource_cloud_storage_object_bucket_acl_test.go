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

const testAccCloudStorageObjectBucketACLName = "ovh_cloud_storage_object_bucket_acl.acl"

func testAccCloudStorageObjectBucketACLImportStateIdFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}
		return fmt.Sprintf("%s/%s", rs.Primary.Attributes["service_name"], rs.Primary.Attributes["bucket_id"]), nil
	}
}

// testAccCheckCloudStorageObjectBucketACL checks the canned ACL of the bucket resource on the S3 API.
func testAccCheckCloudStorageObjectBucketACL(bucketResourceName, want string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client, target, err := testAccCloudStorageObjectBucketS3(s, bucketResourceName)
		if err != nil {
			return err
		}

		policy, err := client.GetBucketACL(context.Background(), target.Host, target.Region)
		if err != nil {
			return fmt.Errorf("error reading the ACL of bucket %s: %w", target.Host, err)
		}

		if got := policy.CannedACL(); got != want {
			return fmt.Errorf("bucket %s has the ACL %q (%+v), want %q", target.Host, got, policy.Grants, want)
		}
		return nil
	}
}

func testAccCloudStorageObjectBucketACLConfig(serviceName, bucketName, region, userId, acl string) string {
	config := fmt.Sprintf(`
resource "ovh_cloud_storage_object_bucket" "bucket" {
  service_name  = "%s"
  name          = "%s"
  region        = "%s"
  owner_user_id = "%s"
}
`, serviceName, bucketName, region, userId)

	if acl != "" {
		config += fmt.Sprintf(`
resource "ovh_cloud_storage_object_bucket_acl" "acl" {
  service_name = ovh_cloud_storage_object_bucket.bucket.service_name
  bucket_id    = ovh_cloud_storage_object_bucket.bucket.id
  acl          = "%s"
}
`, acl)
	}

	return config
}

func TestAccCloudStorageObjectBucketACL_basic(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	region := strings.ToUpper(os.Getenv("OVH_CLOUD_PROJECT_STORAGE_REGION_TEST"))
	userId := os.Getenv("OVH_CLOUD_PROJECT_S3_OWNER_USER_ID_TEST")
	bucketName := acctest.RandomWithPrefix(test_prefix)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCloudStorageObjectBucketS3(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudStorageObjectBucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudStorageObjectBucketACLConfig(serviceName, bucketName, region, userId, "authenticated-read"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(testAccCloudStorageObjectBucketACLName, "service_name", serviceName),
					resource.TestCheckResourceAttrPair(testAccCloudStorageObjectBucketACLName, "bucket_id", "ovh_cloud_storage_object_bucket.bucket", "id"),
					resource.TestCheckResourceAttrPair(testAccCloudStorageObjectBucketACLName, "id", "ovh_cloud_storage_object_bucket.bucket", "id"),
					resource.TestCheckResourceAttr(testAccCloudStorageObjectBucketACLName, "acl", "authenticated-read"),
					testAccCheckCloudStorageObjectBucketACL("ovh_cloud_storage_object_bucket.bucket", "authenticated-read"),
				),
			},
			{
				ResourceName:      testAccCloudStorageObjectBucketACLName,
				ImportState:       true,
				ImportStateIdFunc: testAccCloudStorageObjectBucketACLImportStateIdFunc(testAccCloudStorageObjectBucketACLName),
				ImportStateVerify: true,
			},
			{
				Config: testAccCloudStorageObjectBucketACLConfig(serviceName, bucketName, region, userId, "log-delivery-write"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPostRefresh: []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(testAccCloudStorageObjectBucketACLName, "acl", "log-delivery-write"),
					testAccCheckCloudStorageObjectBucketACL("ovh_cloud_storage_object_bucket.bucket", "log-delivery-write"),
				),
			},
			{
				Config: testAccCloudStorageObjectBucketACLConfig(serviceName, bucketName, region, userId, ""),
				Check:  testAccCheckCloudStorageObjectBucketACL("ovh_cloud_storage_object_bucket.bucket", "private"),
			},
		},
	})
}
