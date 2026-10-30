package ovh

import (
	"fmt"
	"net/url"
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/ovh/go-ovh/ovh"
)

func testAccPreCheckCloudStorageObjectBucket(t *testing.T) {
	testAccPreCheckCloud(t)
	checkEnvOrSkip(t, "OVH_CLOUD_PROJECT_STORAGE_REGION_TEST")
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

func testAccCheckCloudStorageObjectBucketDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ovh_cloud_storage_object_bucket" {
			continue
		}

		endpoint := "/v2/publicCloud/project/" + url.PathEscape(rs.Primary.Attributes["service_name"]) + "/storage/object/bucket/" + url.PathEscape(rs.Primary.Attributes["id"])
		err := testAccOVHClient.Get(endpoint, &CloudStorageObjectBucketAPIResponse{})
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

func TestAccCloudStorageObjectBucket_basic(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	region := strings.ToUpper(os.Getenv("OVH_CLOUD_PROJECT_STORAGE_REGION_TEST"))

	bucketName := acctest.RandomWithPrefix(test_prefix)

	config := fmt.Sprintf(`
resource "ovh_cloud_storage_object_bucket" "bucket" {
  service_name = "%s"
  name         = "%s"
  region       = "%s"
}
`, serviceName, bucketName, region)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCloudStorageObjectBucket(t)
			testAccCheckCloudProjectExists(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudStorageObjectBucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_storage_object_bucket.bucket", "service_name", serviceName),
					resource.TestCheckResourceAttr("ovh_cloud_storage_object_bucket.bucket", "name", bucketName),
					resource.TestCheckResourceAttr("ovh_cloud_storage_object_bucket.bucket", "region", region),
					resource.TestCheckResourceAttr("ovh_cloud_storage_object_bucket.bucket", "resource_status", "READY"),
					resource.TestCheckResourceAttrSet("ovh_cloud_storage_object_bucket.bucket", "id"),
					resource.TestCheckResourceAttrSet("ovh_cloud_storage_object_bucket.bucket", "checksum"),
					resource.TestCheckResourceAttrSet("ovh_cloud_storage_object_bucket.bucket", "created_at"),
					resource.TestCheckResourceAttrSet("ovh_cloud_storage_object_bucket.bucket", "updated_at"),
					resource.TestCheckResourceAttr("ovh_cloud_storage_object_bucket.bucket", "current_state.name", bucketName),
					resource.TestCheckResourceAttr("ovh_cloud_storage_object_bucket.bucket", "current_state.location.region", region),
					resource.TestCheckResourceAttrSet("ovh_cloud_storage_object_bucket.bucket", "current_state.virtual_host"),
					resource.TestCheckNoResourceAttr("ovh_cloud_storage_object_bucket.bucket", "object_lock.%"),
					resource.TestCheckResourceAttr("ovh_cloud_storage_object_bucket.bucket", "current_state.object_lock.status", "DISABLED"),
				),
			},
			// The always-present current_state.object_lock {status = DISABLED} must not re-plan an unset object_lock.
			{
				Config: config,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
			{
				ResourceName:      "ovh_cloud_storage_object_bucket.bucket",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccCloudStorageObjectBucketImportStateIdFunc("ovh_cloud_storage_object_bucket.bucket"),
			},
		},
	})
}

func TestAccCloudStorageObjectBucket_update(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	region := strings.ToUpper(os.Getenv("OVH_CLOUD_PROJECT_STORAGE_REGION_TEST"))
	// Optional: owner_user_id must reference an existing user, so it is only exercised when provided;
	// immutable, so it is set in every step.
	ownerUserId := os.Getenv("OVH_CLOUD_PROJECT_S3_OWNER_USER_ID_TEST")

	bucketName := acctest.RandomWithPrefix(test_prefix)

	ownerLine := ""
	if ownerUserId != "" {
		ownerLine = fmt.Sprintf("owner_user_id = %q", ownerUserId)
	}

	config := fmt.Sprintf(`
resource "ovh_cloud_storage_object_bucket" "bucket" {
  service_name = "%s"
  name         = "%s"
  region       = "%s"
  %s

  versioning = {
    status = "ENABLED"
  }

  tags = {
    env = "test"
  }
}
`, serviceName, bucketName, region, ownerLine)

	updatedConfig := fmt.Sprintf(`
resource "ovh_cloud_storage_object_bucket" "bucket" {
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
resource "ovh_cloud_storage_object_bucket" "bucket" {
  service_name = "%s"
  name         = "%s"
  region       = "%s"
  %s

  versioning = {
    status = "SUSPENDED"
  }

  tags = {
    env = "test"
  }
}
`, serviceName, bucketName, region, ownerLine)

	untaggedConfig := fmt.Sprintf(`
resource "ovh_cloud_storage_object_bucket" "bucket" {
  service_name = "%s"
  name         = "%s"
  region       = "%s"
  %s
}
`, serviceName, bucketName, region, ownerLine)

	updatedChecks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr("ovh_cloud_storage_object_bucket.bucket", "versioning.status", "SUSPENDED"),
		resource.TestCheckResourceAttr("ovh_cloud_storage_object_bucket.bucket", "encryption.algorithm", "AES256"),
		resource.TestCheckResourceAttr("ovh_cloud_storage_object_bucket.bucket", "tags.env", "test"),
		resource.TestCheckResourceAttr("ovh_cloud_storage_object_bucket.bucket", "tags.owner", "terraform"),
		resource.TestCheckResourceAttr("ovh_cloud_storage_object_bucket.bucket", "resource_status", "READY"),
		resource.TestCheckResourceAttrSet("ovh_cloud_storage_object_bucket.bucket", "checksum"),
	}
	if ownerUserId != "" {
		updatedChecks = append(updatedChecks, resource.TestCheckResourceAttr("ovh_cloud_storage_object_bucket.bucket", "owner_user_id", ownerUserId))
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCloudStorageObjectBucket(t)
			testAccCheckCloudProjectExists(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudStorageObjectBucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_storage_object_bucket.bucket", "versioning.status", "ENABLED"),
					resource.TestCheckResourceAttr("ovh_cloud_storage_object_bucket.bucket", "tags.env", "test"),
					resource.TestCheckResourceAttr("ovh_cloud_storage_object_bucket.bucket", "current_state.versioning.status", "ENABLED"),
				),
			},
			{
				Config: updatedConfig,
				Check:  resource.ComposeTestCheckFunc(updatedChecks...),
			},
			{
				Config: prunedConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckNoResourceAttr("ovh_cloud_storage_object_bucket.bucket", "encryption.algorithm"),
					resource.TestCheckResourceAttr("ovh_cloud_storage_object_bucket.bucket", "tags.%", "1"),
					resource.TestCheckResourceAttr("ovh_cloud_storage_object_bucket.bucket", "tags.env", "test"),
					resource.TestCheckNoResourceAttr("ovh_cloud_storage_object_bucket.bucket", "tags.owner"),
					resource.TestCheckNoResourceAttr("ovh_cloud_storage_object_bucket.bucket", "current_state.tags.owner"),
					resource.TestCheckResourceAttr("ovh_cloud_storage_object_bucket.bucket", "resource_status", "READY"),
				),
			},
			{
				Config: untaggedConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckNoResourceAttr("ovh_cloud_storage_object_bucket.bucket", "tags.%"),
					resource.TestCheckNoResourceAttr("ovh_cloud_storage_object_bucket.bucket", "tags.env"),
					resource.TestCheckNoResourceAttr("ovh_cloud_storage_object_bucket.bucket", "current_state.tags.env"),
					resource.TestCheckResourceAttr("ovh_cloud_storage_object_bucket.bucket", "resource_status", "READY"),
				),
			},
		},
	})
}

func TestAccCloudStorageObjectBucket_objectLock(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	region := strings.ToUpper(os.Getenv("OVH_CLOUD_PROJECT_STORAGE_REGION_TEST"))

	bucketName := acctest.RandomWithPrefix(test_prefix)

	config := fmt.Sprintf(`
resource "ovh_cloud_storage_object_bucket" "bucket" {
  service_name = "%s"
  name         = "%s"
  region       = "%s"

  versioning = {
    status = "ENABLED"
  }

  object_lock = {
    status         = "ENABLED"
    mode           = "GOVERNANCE"
    retention_days = 30
  }
}
`, serviceName, bucketName, region)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCloudStorageObjectBucket(t)
			testAccCheckCloudProjectExists(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudStorageObjectBucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_storage_object_bucket.bucket", "object_lock.status", "ENABLED"),
					resource.TestCheckResourceAttr("ovh_cloud_storage_object_bucket.bucket", "object_lock.mode", "GOVERNANCE"),
					resource.TestCheckResourceAttr("ovh_cloud_storage_object_bucket.bucket", "object_lock.retention_days", "30"),
					resource.TestCheckResourceAttr("ovh_cloud_storage_object_bucket.bucket", "current_state.object_lock.status", "ENABLED"),
					resource.TestCheckResourceAttr("ovh_cloud_storage_object_bucket.bucket", "current_state.object_lock.mode", "GOVERNANCE"),
					resource.TestCheckResourceAttr("ovh_cloud_storage_object_bucket.bucket", "current_state.object_lock.retention_days", "30"),
					resource.TestCheckResourceAttr("ovh_cloud_storage_object_bucket.bucket", "versioning.status", "ENABLED"),
					resource.TestCheckResourceAttrSet("ovh_cloud_storage_object_bucket.bucket", "id"),
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
				ResourceName:      "ovh_cloud_storage_object_bucket.bucket",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccCloudStorageObjectBucketImportStateIdFunc("ovh_cloud_storage_object_bucket.bucket"),
			},
		},
	})
}

func TestAccCloudStorageObjectBucket_objectLockArmedWithoutRule(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	region := strings.ToUpper(os.Getenv("OVH_CLOUD_PROJECT_STORAGE_REGION_TEST"))

	bucketName := acctest.RandomWithPrefix(test_prefix)

	config := fmt.Sprintf(`
resource "ovh_cloud_storage_object_bucket" "bucket" {
  service_name = "%s"
  name         = "%s"
  region       = "%s"

  versioning = {
    status = "ENABLED"
  }

  object_lock = {
    status = "ENABLED"
  }
}
`, serviceName, bucketName, region)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCloudStorageObjectBucket(t)
			testAccCheckCloudProjectExists(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudStorageObjectBucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_storage_object_bucket.bucket", "object_lock.status", "ENABLED"),
					resource.TestCheckNoResourceAttr("ovh_cloud_storage_object_bucket.bucket", "object_lock.mode"),
					resource.TestCheckNoResourceAttr("ovh_cloud_storage_object_bucket.bucket", "object_lock.retention_days"),
					resource.TestCheckResourceAttr("ovh_cloud_storage_object_bucket.bucket", "current_state.object_lock.status", "ENABLED"),
					resource.TestCheckNoResourceAttr("ovh_cloud_storage_object_bucket.bucket", "current_state.object_lock.mode"),
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
		},
	})
}

func TestAccCloudStorageObjectBucket_objectLockValidation(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	region := strings.ToUpper(os.Getenv("OVH_CLOUD_PROJECT_STORAGE_REGION_TEST"))

	bucketName := acctest.RandomWithPrefix(test_prefix)

	configFor := func(versioning, objectLock string) string {
		return fmt.Sprintf(`
resource "ovh_cloud_storage_object_bucket" "bucket" {
  service_name = "%s"
  name         = "%s"
  region       = "%s"
  %s
  object_lock = {
    %s
  }
}
`, serviceName, bucketName, region, versioning, objectLock)
	}
	versioningEnabled := `versioning = { status = "ENABLED" }`

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCloudStorageObjectBucket(t)
			testAccCheckCloudProjectExists(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudStorageObjectBucketDestroy,
		Steps: []resource.TestStep{
			{
				Config:      configFor(versioningEnabled, `mode = "GOVERNANCE"`),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`"status" is required`),
			},
			{
				Config:      configFor(versioningEnabled, `status = "ARMED"`),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`must be one of`),
			},
			{
				Config:      configFor("", `status = "ENABLED"`),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`object_lock requires versioning`),
			},
			{
				Config:      configFor(versioningEnabled, "status = \"ENABLED\"\n    mode = \"GOVERNANCE\""),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`must be set together`),
			},
			{
				Config:      configFor(versioningEnabled, "status = \"ENABLED\"\n    retention_days = 30"),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`must be set together`),
			},
			{
				Config:      configFor(versioningEnabled, "status = \"DISABLED\"\n    mode = \"GOVERNANCE\"\n    retention_days = 30"),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`require status "ENABLED"`),
			},
		},
	})
}
