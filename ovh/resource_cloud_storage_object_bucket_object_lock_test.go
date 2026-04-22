package ovh

import (
	"fmt"
	"net/url"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// testAccPreCheckCloudStorageObjectBucketObjectLock skips the test unless the
// target project, bucket, object key and retain-until date are all configured
// through environment variables. The object must exist in a bucket with object
// lock enabled.
func testAccPreCheckCloudStorageObjectBucketObjectLock(t *testing.T) {
	testAccPreCheckCredentials(t)
	if os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST") == "" {
		t.Skip("OVH_CLOUD_PROJECT_SERVICE_TEST must be set for acceptance tests")
	}
	if os.Getenv("OVH_CLOUD_PROJECT_S3_BUCKET_TEST") == "" {
		t.Skip("OVH_CLOUD_PROJECT_S3_BUCKET_TEST must be set for acceptance tests")
	}
	if os.Getenv("OVH_CLOUD_PROJECT_S3_OBJECT_KEY_TEST") == "" {
		t.Skip("OVH_CLOUD_PROJECT_S3_OBJECT_KEY_TEST must be set for acceptance tests")
	}
	if os.Getenv("OVH_CLOUD_PROJECT_S3_OBJECT_RETAIN_UNTIL_TEST") == "" {
		t.Skip("OVH_CLOUD_PROJECT_S3_OBJECT_RETAIN_UNTIL_TEST must be set (RFC3339) for acceptance tests")
	}
}

func testAccCloudStorageObjectBucketObjectLockImportStateIdFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}
		serviceName := rs.Primary.Attributes["service_name"]
		bucket := rs.Primary.Attributes["bucket_name"]
		key := rs.Primary.Attributes["object_key"]
		versionId := rs.Primary.Attributes["version_id"]
		id := fmt.Sprintf("%s/%s/%s", serviceName, bucket, url.PathEscape(key))
		if versionId != "" {
			id += "/" + versionId
		}
		return id, nil
	}
}

func TestAccCloudStorageObjectBucketObjectLock_basic(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	bucket := os.Getenv("OVH_CLOUD_PROJECT_S3_BUCKET_TEST")
	objectKey := os.Getenv("OVH_CLOUD_PROJECT_S3_OBJECT_KEY_TEST")
	retainUntil := os.Getenv("OVH_CLOUD_PROJECT_S3_OBJECT_RETAIN_UNTIL_TEST")

	config := fmt.Sprintf(`
resource "ovh_cloud_storage_object_bucket_object_lock" "test" {
  service_name = "%s"
  bucket_name  = "%s"
  object_key   = "%s"

  retention = {
    mode              = "GOVERNANCE"
    retain_until_date = "%s"
  }

  legal_hold = {
    status = "ON"
  }
}
`, serviceName, bucket, objectKey, retainUntil)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCloudStorageObjectBucketObjectLock(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_storage_object_bucket_object_lock.test", "service_name", serviceName),
					resource.TestCheckResourceAttr("ovh_cloud_storage_object_bucket_object_lock.test", "bucket_name", bucket),
					resource.TestCheckResourceAttr("ovh_cloud_storage_object_bucket_object_lock.test", "object_key", objectKey),
					resource.TestCheckResourceAttr("ovh_cloud_storage_object_bucket_object_lock.test", "retention.mode", "GOVERNANCE"),
					resource.TestCheckResourceAttr("ovh_cloud_storage_object_bucket_object_lock.test", "retention.retain_until_date", retainUntil),
					resource.TestCheckResourceAttr("ovh_cloud_storage_object_bucket_object_lock.test", "legal_hold.status", "ON"),
					resource.TestCheckResourceAttrSet("ovh_cloud_storage_object_bucket_object_lock.test", "id"),
				),
			},
			{
				ResourceName:      "ovh_cloud_storage_object_bucket_object_lock.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccCloudStorageObjectBucketObjectLockImportStateIdFunc("ovh_cloud_storage_object_bucket_object_lock.test"),
			},
		},
	})
}

func TestAccCloudStorageObjectBucketObjectLock_legalHoldUpdate(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	bucket := os.Getenv("OVH_CLOUD_PROJECT_S3_BUCKET_TEST")
	objectKey := os.Getenv("OVH_CLOUD_PROJECT_S3_OBJECT_KEY_TEST")
	retainUntil := os.Getenv("OVH_CLOUD_PROJECT_S3_OBJECT_RETAIN_UNTIL_TEST")

	configOn := fmt.Sprintf(`
resource "ovh_cloud_storage_object_bucket_object_lock" "test" {
  service_name = "%s"
  bucket_name  = "%s"
  object_key   = "%s"

  retention = {
    mode              = "GOVERNANCE"
    retain_until_date = "%s"
  }

  legal_hold = {
    status = "ON"
  }
}
`, serviceName, bucket, objectKey, retainUntil)

	configOff := fmt.Sprintf(`
resource "ovh_cloud_storage_object_bucket_object_lock" "test" {
  service_name = "%s"
  bucket_name  = "%s"
  object_key   = "%s"

  retention = {
    mode              = "GOVERNANCE"
    retain_until_date = "%s"
  }

  legal_hold = {
    status = "OFF"
  }
}
`, serviceName, bucket, objectKey, retainUntil)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCloudStorageObjectBucketObjectLock(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: configOn,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_storage_object_bucket_object_lock.test", "legal_hold.status", "ON"),
				),
			},
			{
				Config: configOff,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_storage_object_bucket_object_lock.test", "legal_hold.status", "OFF"),
				),
			},
		},
	})
}
