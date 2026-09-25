package ovh

import (
	"fmt"
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/ovh/go-ovh/ovh"
)

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

func testAccCheckCloudS3BucketDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ovh_cloud_s3_bucket" {
			continue
		}

		endpoint := "/v2/publicCloud/project/" + url.PathEscape(rs.Primary.Attributes["service_name"]) + "/storage/object/bucket/" + url.PathEscape(rs.Primary.Attributes["id"])
		err := testAccOVHClient.Get(endpoint, &CloudS3BucketAPIResponse{})
		if err == nil {
			return fmt.Errorf("bucket %s still exists", rs.Primary.Attributes["id"])
		}
		if errOvh, ok := err.(*ovh.APIError); ok && errOvh.Code == 404 {
			continue
		}
		return fmt.Errorf("error checking bucket %s was destroyed: %w", rs.Primary.Attributes["id"], err)
	}
	return nil
}

func TestAccCloudS3Bucket_basic(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	region := strings.ToUpper(os.Getenv("OVH_CLOUD_PROJECT_STORAGE_REGION_TEST"))

	bucketName := acctest.RandomWithPrefix(test_prefix)

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
		CheckDestroy:             testAccCheckCloudS3BucketDestroy,
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
	region := strings.ToUpper(os.Getenv("OVH_CLOUD_PROJECT_STORAGE_REGION_TEST"))
	// Optional: owner_user_id must reference an existing user, so it is only exercised when provided.
	ownerUserId := os.Getenv("OVH_CLOUD_PROJECT_S3_OWNER_USER_ID_TEST")

	bucketName := acctest.RandomWithPrefix(test_prefix)

	ownerLine := ""
	if ownerUserId != "" {
		ownerLine = fmt.Sprintf("owner_user_id = %q", ownerUserId)
	}

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
  %s

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
`, serviceName, bucketName, region, ownerLine)

	prunedConfig := fmt.Sprintf(`
resource "ovh_cloud_s3_bucket" "bucket" {
  service_name = "%s"
  name         = "%s"
  region       = "%s"

  versioning = {
    status = "SUSPENDED"
  }

  tags = {
    env = "test"
  }
}
`, serviceName, bucketName, region)

	updatedChecks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr("ovh_cloud_s3_bucket.bucket", "versioning.status", "SUSPENDED"),
		resource.TestCheckResourceAttr("ovh_cloud_s3_bucket.bucket", "encryption.algorithm", "AES256"),
		resource.TestCheckResourceAttr("ovh_cloud_s3_bucket.bucket", "tags.env", "test"),
		resource.TestCheckResourceAttr("ovh_cloud_s3_bucket.bucket", "tags.owner", "terraform"),
		resource.TestCheckResourceAttr("ovh_cloud_s3_bucket.bucket", "resource_status", "READY"),
		resource.TestCheckResourceAttrSet("ovh_cloud_s3_bucket.bucket", "checksum"),
	}
	if ownerUserId != "" {
		updatedChecks = append(updatedChecks, resource.TestCheckResourceAttr("ovh_cloud_s3_bucket.bucket", "owner_user_id", ownerUserId))
	}

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
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_s3_bucket.bucket", "versioning.status", "ENABLED"),
					resource.TestCheckResourceAttr("ovh_cloud_s3_bucket.bucket", "tags.env", "test"),
					resource.TestCheckResourceAttr("ovh_cloud_s3_bucket.bucket", "current_state.versioning.status", "ENABLED"),
				),
			},
			{
				Config: updatedConfig,
				Check:  resource.ComposeTestCheckFunc(updatedChecks...),
			},
			{
				Config: prunedConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckNoResourceAttr("ovh_cloud_s3_bucket.bucket", "owner_user_id"),
					resource.TestCheckNoResourceAttr("ovh_cloud_s3_bucket.bucket", "encryption.algorithm"),
					resource.TestCheckResourceAttr("ovh_cloud_s3_bucket.bucket", "tags.%", "1"),
					resource.TestCheckResourceAttr("ovh_cloud_s3_bucket.bucket", "tags.env", "test"),
					resource.TestCheckNoResourceAttr("ovh_cloud_s3_bucket.bucket", "tags.owner"),
					resource.TestCheckNoResourceAttr("ovh_cloud_s3_bucket.bucket", "current_state.tags.owner"),
					resource.TestCheckResourceAttr("ovh_cloud_s3_bucket.bucket", "resource_status", "READY"),
				),
			},
		},
	})
}

func TestAccCloudS3Bucket_objectLock(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	region := strings.ToUpper(os.Getenv("OVH_CLOUD_PROJECT_STORAGE_REGION_TEST"))

	bucketName := acctest.RandomWithPrefix(test_prefix)

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
		CheckDestroy:             testAccCheckCloudS3BucketDestroy,
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
			// Guards against a perpetually-unknown object_lock child re-planning an update or replace.
			{
				Config: config,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
			{
				ResourceName:      "ovh_cloud_s3_bucket.bucket",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccCloudS3BucketImportStateIdFunc("ovh_cloud_s3_bucket.bucket"),
			},
		},
	})
}
