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
				),
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
	// Optional: owner_user_id must reference an existing user, so it is only exercised when provided.
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

  versioning = {
    status = "ENABLED"
  }

  tags = {
    env = "test"
  }
}
`, serviceName, bucketName, region)

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

  versioning = {
    status = "SUSPENDED"
  }

  tags = {
    env = "test"
  }
}
`, serviceName, bucketName, region)

	untaggedConfig := fmt.Sprintf(`
resource "ovh_cloud_storage_object_bucket" "bucket" {
  service_name = "%s"
  name         = "%s"
  region       = "%s"
}
`, serviceName, bucketName, region)

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
					resource.TestCheckNoResourceAttr("ovh_cloud_storage_object_bucket.bucket", "owner_user_id"),
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
					resource.TestCheckResourceAttr("ovh_cloud_storage_object_bucket.bucket", "object_lock.mode", "GOVERNANCE"),
					resource.TestCheckResourceAttr("ovh_cloud_storage_object_bucket.bucket", "object_lock.retention_days", "30"),
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

// testAccCloudStorageObjectBucketForceDestroyConfig returns a bucket and an S3 user. The object
// step uploads an object with the AWS provider, once the S3 credential exists. Once forgotten
// with a removed block, the object stays in the bucket when it is destroyed.
func testAccCloudStorageObjectBucketForceDestroyConfig(serviceName, region, bucketName, bucket, objectStep string) string {
	provider := `
provider "aws" {
  region                      = lower("%[2]s")
  access_key                  = ovh_cloud_project_user_s3_credential.credential.access_key_id
  secret_key                  = ovh_cloud_project_user_s3_credential.credential.secret_access_key
  skip_credentials_validation = true
  skip_metadata_api_check     = true
  skip_region_validation      = true
  skip_requesting_account_id  = true
  endpoints {
    s3 = "https://s3.${lower("%[2]s")}.io.cloud.ovh.net"
  }
}

`
	object := provider + `
resource "aws_s3_object" "object" {
  bucket     = "%[3]s"
  key        = "dir/object.txt"
  content    = "content"
  depends_on = [ovh_cloud_storage_object_bucket.bucket, time_sleep.policy]
}
`
	switch objectStep {
	case "none":
		object = ""
	case "removed":
		object = provider + `
removed {
  from = aws_s3_object.object
  lifecycle {
    destroy = false
  }
}
`
	}

	if object != "" {
		object = fmt.Sprintf(object, serviceName, region, bucketName)
	}

	return fmt.Sprintf(`
resource "ovh_cloud_project_user" "user" {
  service_name = "%[1]s"
  description  = "%[3]s"
  role_name    = "objectstore_operator"
}

resource "ovh_cloud_project_user_s3_credential" "credential" {
  service_name = "%[1]s"
  user_id      = ovh_cloud_project_user.user.id
}

resource "ovh_cloud_project_user_s3_policy" "policy" {
  service_name = "%[1]s"
  user_id      = ovh_cloud_project_user.user.id
  policy = jsonencode({
    Statement = [{
      Sid      = "all"
      Effect   = "Allow"
      Action   = ["s3:*"]
      Resource = ["arn:aws:s3:::%[3]s", "arn:aws:s3:::%[3]s/*"]
    }]
  })
}

# The S3 policy takes a few seconds to apply
resource "time_sleep" "policy" {
  create_duration = "30s"
  depends_on      = [ovh_cloud_project_user_s3_policy.policy]
}
%[4]s
%[5]s
`, serviceName, region, bucketName, bucket, object)
}

func TestAccCloudStorageObjectBucket_forceDestroy(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	region := strings.ToUpper(os.Getenv("OVH_CLOUD_PROJECT_STORAGE_REGION_TEST"))
	bucketName := acctest.RandomWithPrefix(test_prefix)

	bucket := func(forceDestroy string) string {
		return fmt.Sprintf(`
resource "ovh_cloud_storage_object_bucket" "bucket" {
  service_name  = "%s"
  name          = "%s"
  region        = "%s"
  owner_user_id = ovh_cloud_project_user.user.id
  %s
}
`, serviceName, bucketName, region, forceDestroy)
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCloudStorageObjectBucket(t)
			testAccCheckCloudProjectExists(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"aws":  {VersionConstraint: "5.82.2", Source: "registry.terraform.io/hashicorp/aws"},
			"time": {VersionConstraint: "0.13.1", Source: "registry.terraform.io/hashicorp/time"},
		},
		CheckDestroy: testAccCheckCloudStorageObjectBucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudStorageObjectBucketForceDestroyConfig(serviceName, region, bucketName, bucket(""), "none"),
				Check:  resource.TestCheckResourceAttr("ovh_cloud_storage_object_bucket.bucket", "force_destroy", "false"),
			},
			{
				Config: testAccCloudStorageObjectBucketForceDestroyConfig(serviceName, region, bucketName, bucket(""), "managed"),
			},
			{
				// Destroying the bucket with its object fails, without sending the deletion
				Config:      testAccCloudStorageObjectBucketForceDestroyConfig(serviceName, region, bucketName, "", "removed"),
				ExpectError: regexp.MustCompile("is not empty"),
			},
			{
				// The bucket is still there, and force_destroy is applied without calling the API
				Config: testAccCloudStorageObjectBucketForceDestroyConfig(serviceName, region, bucketName, bucket("force_destroy = true"), "removed"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_storage_object_bucket.bucket", "force_destroy", "true"),
					resource.TestCheckResourceAttr("ovh_cloud_storage_object_bucket.bucket", "resource_status", "READY"),
				),
			},
			// The final destroy empties the bucket, then deletes it
		},
	})
}
