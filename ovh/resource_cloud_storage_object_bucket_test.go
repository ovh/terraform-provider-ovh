package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

const testAccResourceCloudStorageObjectBucketNamePrefix = "tf-test-bucket-v2-"

func testAccPreCheckCloudStorageObjectBucket(t *testing.T) {
	testAccPreCheckCredentials(t)
	if os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST") == "" {
		t.Skip("OVH_CLOUD_PROJECT_SERVICE_TEST must be set for acceptance tests")
	}
	if os.Getenv("OVH_CLOUD_PROJECT_REGION_TEST") == "" {
		t.Skip("OVH_CLOUD_PROJECT_REGION_TEST must be set for acceptance tests")
	}
}

func testAccCloudStorageObjectBucketImportStateIdFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}
		return fmt.Sprintf("%s/%s", rs.Primary.Attributes["service_name"], rs.Primary.Attributes["id"]), nil
	}
}

func TestAccCloudStorageObjectBucket_basic(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	region := os.Getenv("OVH_CLOUD_PROJECT_REGION_TEST")

	bucketName := acctest.RandomWithPrefix(testAccResourceCloudStorageObjectBucketNamePrefix)

	config := fmt.Sprintf(`
resource "ovh_cloud_storage_object_bucket" "test" {
  service_name = "%s"
  name         = "%s"
  region       = "%s"
}
`, serviceName, bucketName, region)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCloudStorageObjectBucket(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_storage_object_bucket.test", "service_name", serviceName),
					resource.TestCheckResourceAttr("ovh_cloud_storage_object_bucket.test", "name", bucketName),
					resource.TestCheckResourceAttr("ovh_cloud_storage_object_bucket.test", "region", region),
					resource.TestCheckResourceAttrSet("ovh_cloud_storage_object_bucket.test", "id"),
					resource.TestCheckResourceAttrSet("ovh_cloud_storage_object_bucket.test", "checksum"),
					resource.TestCheckResourceAttrSet("ovh_cloud_storage_object_bucket.test", "created_at"),
					resource.TestCheckResourceAttr("ovh_cloud_storage_object_bucket.test", "resource_status", "READY"),
					resource.TestCheckResourceAttrSet("ovh_cloud_storage_object_bucket.test", "current_state.name"),
					resource.TestCheckResourceAttrSet("ovh_cloud_storage_object_bucket.test", "current_state.location.region"),
				),
			},
			// Test import
			{
				ResourceName:      "ovh_cloud_storage_object_bucket.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccCloudStorageObjectBucketImportStateIdFunc("ovh_cloud_storage_object_bucket.test"),
			},
		},
	})
}

func TestAccCloudStorageObjectBucket_update(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	region := os.Getenv("OVH_CLOUD_PROJECT_REGION_TEST")

	bucketName := acctest.RandomWithPrefix(testAccResourceCloudStorageObjectBucketNamePrefix)

	config := fmt.Sprintf(`
resource "ovh_cloud_storage_object_bucket" "test" {
  service_name = "%s"
  name         = "%s"
  region       = "%s"
}
`, serviceName, bucketName, region)

	updatedConfig := fmt.Sprintf(`
resource "ovh_cloud_storage_object_bucket" "test" {
  service_name = "%s"
  name         = "%s"
  region       = "%s"
  versioning   = {
    status = "ENABLED"
  }
}
`, serviceName, bucketName, region)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCloudStorageObjectBucket(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_storage_object_bucket.test", "name", bucketName),
					resource.TestCheckResourceAttrSet("ovh_cloud_storage_object_bucket.test", "id"),
					resource.TestCheckResourceAttrSet("ovh_cloud_storage_object_bucket.test", "checksum"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_storage_object_bucket.test", "name", bucketName),
					resource.TestCheckResourceAttr("ovh_cloud_storage_object_bucket.test", "versioning.status", "ENABLED"),
					resource.TestCheckResourceAttrSet("ovh_cloud_storage_object_bucket.test", "id"),
					resource.TestCheckResourceAttrSet("ovh_cloud_storage_object_bucket.test", "checksum"),
					resource.TestCheckResourceAttr("ovh_cloud_storage_object_bucket.test", "resource_status", "READY"),
				),
			},
		},
	})
}
