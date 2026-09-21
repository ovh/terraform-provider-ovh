package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

const testAccCloudS3BucketNamePrefix = "tf-test-s3-bucket-v2"

func testAccPreCheckCloudS3Bucket(t *testing.T) {
	testAccPreCheckCloud(t)
	checkEnvOrSkip(t, "OVH_CLOUD_PROJECT_STORAGE_REGION_TEST")
}

func testAccCloudS3BucketImportStateIdFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}
		return fmt.Sprintf("%s/%s", rs.Primary.Attributes["service_name"], rs.Primary.Attributes["id"]), nil
	}
}

func TestAccCloudS3Bucket_basic(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	region := os.Getenv("OVH_CLOUD_PROJECT_STORAGE_REGION_TEST")

	bucketName := acctest.RandomWithPrefix(testAccCloudS3BucketNamePrefix)

	config := fmt.Sprintf(`
resource "ovh_cloud_s3_bucket" "bucket" {
  service_name = "%s"
  name         = "%s"
  region       = "%s"
}
`, serviceName, bucketName, region)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCloudS3Bucket(t)
			testAccCheckCloudProjectExists(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_s3_bucket.bucket", "service_name", serviceName),
					resource.TestCheckResourceAttr("ovh_cloud_s3_bucket.bucket", "name", bucketName),
					resource.TestCheckResourceAttr("ovh_cloud_s3_bucket.bucket", "region", region),
					resource.TestCheckResourceAttr("ovh_cloud_s3_bucket.bucket", "resource_status", "READY"),
					resource.TestCheckResourceAttrSet("ovh_cloud_s3_bucket.bucket", "id"),
					resource.TestCheckResourceAttrSet("ovh_cloud_s3_bucket.bucket", "checksum"),
					resource.TestCheckResourceAttrSet("ovh_cloud_s3_bucket.bucket", "created_at"),
					resource.TestCheckResourceAttrSet("ovh_cloud_s3_bucket.bucket", "updated_at"),
					resource.TestCheckResourceAttr("ovh_cloud_s3_bucket.bucket", "current_state.name", bucketName),
					resource.TestCheckResourceAttr("ovh_cloud_s3_bucket.bucket", "current_state.location.region", region),
					resource.TestCheckResourceAttrSet("ovh_cloud_s3_bucket.bucket", "current_state.virtual_host"),
				),
			},
			// Test import
			{
				ResourceName:      "ovh_cloud_s3_bucket.bucket",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccCloudS3BucketImportStateIdFunc("ovh_cloud_s3_bucket.bucket"),
			},
		},
	})
}

func TestAccCloudS3Bucket_update(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	region := os.Getenv("OVH_CLOUD_PROJECT_STORAGE_REGION_TEST")

	bucketName := acctest.RandomWithPrefix(testAccCloudS3BucketNamePrefix)

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
`, serviceName, bucketName, region)

	updatedConfig := fmt.Sprintf(`
resource "ovh_cloud_s3_bucket" "bucket" {
  service_name = "%s"
  name         = "%s"
  region       = "%s"

  versioning = {
    status = "SUSPENDED"
  }

  encryption = {
    algorithm = "AES256"
  }

  tags = {
    env   = "test"
    owner = "terraform"
  }
}
`, serviceName, bucketName, region)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCloudS3Bucket(t)
			testAccCheckCloudProjectExists(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_s3_bucket.bucket", "versioning.status", "ENABLED"),
					resource.TestCheckResourceAttr("ovh_cloud_s3_bucket.bucket", "tags.env", "test"),
					resource.TestCheckResourceAttr("ovh_cloud_s3_bucket.bucket", "current_state.versioning.status", "ENABLED"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_s3_bucket.bucket", "versioning.status", "SUSPENDED"),
					resource.TestCheckResourceAttr("ovh_cloud_s3_bucket.bucket", "encryption.algorithm", "AES256"),
					resource.TestCheckResourceAttr("ovh_cloud_s3_bucket.bucket", "tags.env", "test"),
					resource.TestCheckResourceAttr("ovh_cloud_s3_bucket.bucket", "tags.owner", "terraform"),
					resource.TestCheckResourceAttr("ovh_cloud_s3_bucket.bucket", "resource_status", "READY"),
					resource.TestCheckResourceAttrSet("ovh_cloud_s3_bucket.bucket", "checksum"),
				),
			},
		},
	})
}

func TestAccCloudS3Bucket_objectLock(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	region := os.Getenv("OVH_CLOUD_PROJECT_STORAGE_REGION_TEST")

	bucketName := acctest.RandomWithPrefix(testAccCloudS3BucketNamePrefix)

	// The API rejects object lock unless versioning is ENABLED.
	config := fmt.Sprintf(`
resource "ovh_cloud_s3_bucket" "bucket" {
  service_name = "%s"
  name         = "%s"
  region       = "%s"

  versioning = {
    status = "ENABLED"
  }

  object_lock = {
    mode           = "GOVERNANCE"
    retention_days = 30
  }
}
`, serviceName, bucketName, region)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCloudS3Bucket(t)
			testAccCheckCloudProjectExists(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_s3_bucket.bucket", "object_lock.mode", "GOVERNANCE"),
					resource.TestCheckResourceAttr("ovh_cloud_s3_bucket.bucket", "object_lock.retention_days", "30"),
					resource.TestCheckResourceAttr("ovh_cloud_s3_bucket.bucket", "current_state.object_lock.retention_days", "30"),
					resource.TestCheckResourceAttr("ovh_cloud_s3_bucket.bucket", "versioning.status", "ENABLED"),
					resource.TestCheckResourceAttrSet("ovh_cloud_s3_bucket.bucket", "id"),
				),
			},
			// Test import
			{
				ResourceName:      "ovh_cloud_s3_bucket.bucket",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccCloudS3BucketImportStateIdFunc("ovh_cloud_s3_bucket.bucket"),
			},
		},
	})
}
