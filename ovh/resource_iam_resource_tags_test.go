package ovh

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func init() {
	resource.AddTestSweepers("ovh_iam_resource_tags", &resource.Sweeper{
		Name: "ovh_iam_resource_tags",
		F:    testSweepIamResourceTags,
	})
}

func testSweepIamResourceTags(region string) error {
	client, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting client: %s", err)
	}

	// Get the test cloud project service name
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	if serviceName == "" {
		log.Print("[DEBUG] OVH_CLOUD_PROJECT_SERVICE_TEST not set, skipping IAM resource tags sweep")
		return nil
	}

	// Get the cloud project to obtain its URN
	var project CloudProject
	endpoint := fmt.Sprintf("/cloud/project/%s", url.PathEscape(serviceName))
	if err := client.Get(endpoint, &project); err != nil {
		return fmt.Errorf("Error calling GET %s:\n\t %q", endpoint, err)
	}

	if project.URN == "" {
		log.Print("[DEBUG] Cloud project URN is empty, skipping IAM resource tags sweep")
		return nil
	}

	// Get resource details including tags
	var resourceDetails IamResourceDetails
	resourceEndpoint := fmt.Sprintf("/v2/iam/resource/%s", url.PathEscape(project.URN))
	if err := client.Get(resourceEndpoint, &resourceDetails); err != nil {
		log.Printf("[DEBUG] Error getting resource details for %s: %v", project.URN, err)
		return nil
	}

	if len(resourceDetails.Tags) == 0 {
		log.Print("[DEBUG] No tags to sweep on cloud project")
		return nil
	}

	// Delete tags that start with test prefix
	for key := range resourceDetails.Tags {
		// Only delete tags with keys or values that contain the test prefix
		if !strings.Contains(key, test_prefix) {
			continue
		}

		log.Printf("[DEBUG] Found test tag to sweep: %s on resource %s", key, project.URN)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			log.Printf("[INFO] Deleting tag %s from resource %s", key, project.URN)
			deleteEndpoint := fmt.Sprintf("/v2/iam/resource/%s/tag/%s", url.PathEscape(project.URN), url.PathEscape(key))
			if err := client.Delete(deleteEndpoint, nil); err != nil {
				return resource.RetryableError(err)
			}

			// Successful delete
			return nil
		})

		if err != nil {
			return err
		}
	}

	return nil
}

// TestAccResourceIamResource_basic tests basic create, read, and delete operations
func TestAccResourceIamResourceTags_basic(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	resourceName := "ovh_iam_resource_tags.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckCloud(t); testAccCheckCloudProjectExists(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceIamResourceTagsConfig_basic(serviceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "urn"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(resourceName, fmt.Sprintf("tags.%s-environment", test_prefix), "test"),
					resource.TestCheckResourceAttr(resourceName, fmt.Sprintf("tags.%s-managed_by", test_prefix), "terraform"),
				),
			},
		},
	})
}

// TestAccResourceIamResource_update tests updating tags (add, modify, remove)
func TestAccResourceIamResourceTags_update(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	resourceName := "ovh_iam_resource_tags.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckCloud(t); testAccCheckCloudProjectExists(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				// Step 1: Create with initial tags
				Config: testAccResourceIamResourceTagsConfig_basic(serviceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "urn"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(resourceName, fmt.Sprintf("tags.%s-environment", test_prefix), "test"),
					resource.TestCheckResourceAttr(resourceName, fmt.Sprintf("tags.%s-managed_by", test_prefix), "terraform"),
				),
			},
			{
				// Step 2: Update tags - modify existing, add new, remove one
				Config: testAccResourceIamResourceTagsConfig_update(serviceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "urn"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "3"),
					resource.TestCheckResourceAttr(resourceName, fmt.Sprintf("tags.%s-environment", test_prefix), "production"),
					resource.TestCheckResourceAttr(resourceName, fmt.Sprintf("tags.%s-managed_by", test_prefix), "terraform"),
					resource.TestCheckResourceAttr(resourceName, fmt.Sprintf("tags.%s-team", test_prefix), "platform"),
				),
			},
			{
				// Step 3: Remove all tags except one
				Config: testAccResourceIamResourceTagsConfig_minimal(serviceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "urn"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, fmt.Sprintf("tags.%s-managed_by", test_prefix), "terraform"),
				),
			},
		},
	})
}

// TestAccResourceIamResource_import tests resource import functionality
func TestAccResourceIamResourceTags_import(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	resourceName := "ovh_iam_resource_tags.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckCloud(t); testAccCheckCloudProjectExists(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceIamResourceTagsConfig_basic(serviceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "urn"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(resourceName, fmt.Sprintf("tags.%s-environment", test_prefix), "test"),
					resource.TestCheckResourceAttr(resourceName, fmt.Sprintf("tags.%s-managed_by", test_prefix), "terraform"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// TestAccResourceIamResource_complexTags tests tags with special characters
func TestAccResourceIamResourceTags_complexTags(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	resourceName := "ovh_iam_resource_tags.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckCloud(t); testAccCheckCloudProjectExists(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceIamResourceTagsConfig_complex(serviceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "urn"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "5"),
					resource.TestCheckResourceAttr(resourceName, fmt.Sprintf("tags.%s-app:name", test_prefix), "my-app"),
					resource.TestCheckResourceAttr(resourceName, fmt.Sprintf("tags.%s-cost_center", test_prefix), "12345"),
					resource.TestCheckResourceAttr(resourceName, fmt.Sprintf("tags.%s-owner@email", test_prefix), "team@example.com"),
					resource.TestCheckResourceAttr(resourceName, fmt.Sprintf("tags.%s-version_1.0", test_prefix), "stable"),
					resource.TestCheckResourceAttr(resourceName, fmt.Sprintf("tags.%s-env+region", test_prefix), "prod+eu-west"),
				),
			},
		},
	})
}

// TestAccResourceIamResource_emptyTags tests resource with no tags
func TestAccResourceIamResourceTags_emptyTags(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	resourceName := "ovh_iam_resource_tags.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckCloud(t); testAccCheckCloudProjectExists(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceIamResourceTagsConfig_empty(serviceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "urn"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "0"),
				),
			},
			{
				// Add tags after starting with empty
				Config: testAccResourceIamResourceTagsConfig_basic(serviceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "urn"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(resourceName, fmt.Sprintf("tags.%s-environment", test_prefix), "test"),
					resource.TestCheckResourceAttr(resourceName, fmt.Sprintf("tags.%s-managed_by", test_prefix), "terraform"),
				),
			},
		},
	})
}

// TestAccResourceIamResource_manyTags tests resource with many tags
func TestAccResourceIamResourceTags_manyTags(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	resourceName := "ovh_iam_resource_tags.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckCloud(t); testAccCheckCloudProjectExists(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceIamResourceTagsConfig_many(serviceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "urn"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "10"),
					resource.TestCheckResourceAttr(resourceName, fmt.Sprintf("tags.%s-tag1", test_prefix), "value1"),
					resource.TestCheckResourceAttr(resourceName, fmt.Sprintf("tags.%s-tag2", test_prefix), "value2"),
					resource.TestCheckResourceAttr(resourceName, fmt.Sprintf("tags.%s-tag3", test_prefix), "value3"),
					resource.TestCheckResourceAttr(resourceName, fmt.Sprintf("tags.%s-tag4", test_prefix), "value4"),
					resource.TestCheckResourceAttr(resourceName, fmt.Sprintf("tags.%s-tag5", test_prefix), "value5"),
					resource.TestCheckResourceAttr(resourceName, fmt.Sprintf("tags.%s-tag6", test_prefix), "value6"),
					resource.TestCheckResourceAttr(resourceName, fmt.Sprintf("tags.%s-tag7", test_prefix), "value7"),
					resource.TestCheckResourceAttr(resourceName, fmt.Sprintf("tags.%s-tag8", test_prefix), "value8"),
					resource.TestCheckResourceAttr(resourceName, fmt.Sprintf("tags.%s-tag9", test_prefix), "value9"),
					resource.TestCheckResourceAttr(resourceName, fmt.Sprintf("tags.%s-tag10", test_prefix), "value10"),
				),
			},
		},
	})
}

// TestAccResourceIamResource_resourceWithOvhPrefixedTags tests resource with ovh: prefixed tags
func TestAccResourceIamResourceTags_resourceWithOvhPrefixedTags(t *testing.T) {
	ovhPrefixedResourceURN := os.Getenv("OVH_IAM_RESOURCE_OVH_PREFIXED_IP_URN_TEST")
	resourceName := "ovh_iam_resource_tags.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckIamResourceTagsOvhPrefixed(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceIamResourceTagsConfig_resourceWithOvhPrefixedTags(ovhPrefixedResourceURN),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "urn"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "5"),
					resource.TestCheckResourceAttr(resourceName, fmt.Sprintf("tags.%s-tag1", test_prefix), "value1"),
					resource.TestCheckResourceAttr(resourceName, "tags.ovh:isAdditionalIp", "true"),
					resource.TestCheckResourceAttrSet(resourceName, "tags.ovh:routedTo"),
					resource.TestCheckResourceAttrSet(resourceName, "tags.ovh:type"),
					resource.TestCheckResourceAttrSet(resourceName, "tags.ovh:version"),
				),
			},
		},
	})
}

// TestAccResourceIamResource_invalidTagKey tests validation of invalid tag keys
func TestAccResourceIamResourceTags_invalidTagKey(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckCloud(t); testAccCheckCloudProjectExists(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccResourceIamResourceTagsConfig_invalidKey(serviceName),
				ExpectError: regexp.MustCompile("Invalid tag key"),
			},
		},
	})
}

// TestAccResourceIamResource_invalidTagValue tests validation of invalid tag values
func TestAccResourceIamResourceTags_invalidTagValue(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckCloud(t); testAccCheckCloudProjectExists(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccResourceIamResourceTagsConfig_invalidValue(serviceName),
				ExpectError: regexp.MustCompile("Invalid tag value"),
			},
		},
	})
}

// TestAccResourceIamResource_tooLongTagKey tests validation of tag keys that are too long
func TestAccResourceIamResourceTags_tooLongTagKey(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckCloud(t); testAccCheckCloudProjectExists(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccResourceIamResourceTagsConfig_tooLongKey(serviceName),
				ExpectError: regexp.MustCompile("Invalid tag key"),
			},
		},
	})
}

// TestAccResourceIamResource_tooLongTagValue tests validation of tag values that are too long
func TestAccResourceIamResourceTags_tooLongTagValue(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckCloud(t); testAccCheckCloudProjectExists(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccResourceIamResourceTagsConfig_tooLongValue(serviceName),
				ExpectError: regexp.MustCompile("Invalid tag value"),
			},
		},
	})
}

// TestAccResourceIamResource_ovhPrefixTagKey tests validation of tag keys with ovh: prefix
func TestAccResourceIamResourceTags_ovhPrefixTagKey(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckCloud(t); testAccCheckCloudProjectExists(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccResourceIamResourceTagsConfig_ovhPrefix(serviceName),
				ExpectError: regexp.MustCompile("cannot start with 'ovh:' prefix"),
			},
		},
	})
}

// Configuration helper functions

// testAccResourceIamResourceTagsConfig_basic returns a basic configuration with 2 tags
func testAccResourceIamResourceTagsConfig_basic(serviceName string) string {
	return fmt.Sprintf(`
data "ovh_cloud_project" "project" {
  service_name = "%s"
}

resource "ovh_iam_resource_tags" "test" {
  urn = data.ovh_cloud_project.project.iam.urn

  tags = {
    "%s-environment" = "test"
    "%s-managed_by"  = "terraform"
  }
}
`, serviceName, test_prefix, test_prefix)
}

// testAccResourceIamResourceTagsConfig_update returns a configuration with updated tags
func testAccResourceIamResourceTagsConfig_update(serviceName string) string {
	return fmt.Sprintf(`
data "ovh_cloud_project" "project" {
  service_name = "%s"
}

resource "ovh_iam_resource_tags" "test" {
  urn = data.ovh_cloud_project.project.iam.urn

  tags = {
    "%s-environment" = "production"
    "%s-managed_by"  = "terraform"
    "%s-team"        = "platform"
  }
}
`, serviceName, test_prefix, test_prefix, test_prefix)
}

// testAccResourceIamResourceTagsConfig_minimal returns a configuration with minimal tags
func testAccResourceIamResourceTagsConfig_minimal(serviceName string) string {
	return fmt.Sprintf(`
data "ovh_cloud_project" "project" {
  service_name = "%s"
}

resource "ovh_iam_resource_tags" "test" {
  urn = data.ovh_cloud_project.project.iam.urn

  tags = {
    "%s-managed_by" = "terraform"
  }
}
`, serviceName, test_prefix)
}

// testAccResourceIamResourceTagsConfig_complex returns a configuration with complex tag keys
func testAccResourceIamResourceTagsConfig_complex(serviceName string) string {
	return fmt.Sprintf(`
data "ovh_cloud_project" "project" {
  service_name = "%s"
}

resource "ovh_iam_resource_tags" "test" {
  urn = data.ovh_cloud_project.project.iam.urn

  tags = {
    "%s-app:name"      = "my-app"
    "%s-cost_center"   = "12345"
    "%s-owner@email"   = "team@example.com"
    "%s-version_1.0"   = "stable"
    "%s-env+region"    = "prod+eu-west"
  }
}
`, serviceName, test_prefix, test_prefix, test_prefix, test_prefix, test_prefix)
}

// testAccResourceIamResourceTagsConfig_empty returns a configuration with no tags
func testAccResourceIamResourceTagsConfig_empty(serviceName string) string {
	return fmt.Sprintf(`
data "ovh_cloud_project" "project" {
  service_name = "%s"
}

resource "ovh_iam_resource_tags" "test" {
  urn = data.ovh_cloud_project.project.iam.urn

  tags = {}
}
`, serviceName)
}

// testAccResourceIamResourceTagsConfig_many returns a configuration with many tags
func testAccResourceIamResourceTagsConfig_many(serviceName string) string {
	return fmt.Sprintf(`
data "ovh_cloud_project" "project" {
  service_name = "%s"
}

resource "ovh_iam_resource_tags" "test" {
  urn = data.ovh_cloud_project.project.iam.urn

  tags = {
    "%s-tag1"  = "value1"
    "%s-tag2"  = "value2"
    "%s-tag3"  = "value3"
    "%s-tag4"  = "value4"
    "%s-tag5"  = "value5"
    "%s-tag6"  = "value6"
    "%s-tag7"  = "value7"
    "%s-tag8"  = "value8"
    "%s-tag9"  = "value9"
    "%s-tag10" = "value10"
  }
}
`, serviceName, test_prefix, test_prefix, test_prefix, test_prefix, test_prefix, test_prefix, test_prefix, test_prefix, test_prefix, test_prefix)
}

// testAccResourceIamResourceTagsConfig_resourceWithOvhPrefixedTags returns a configuration with a resource having ovh: prefixed tags
func testAccResourceIamResourceTagsConfig_resourceWithOvhPrefixedTags(ovhPrefixedResourceURN string) string {
	return fmt.Sprintf(`
resource "ovh_iam_resource_tags" "test" {
  urn = "%s"

  tags = {
    "%s-tag1"  = "value1"
  }
}
`, ovhPrefixedResourceURN, test_prefix)
}

// testAccResourceIamResourceTagsConfig_invalidKey returns a configuration with an invalid tag key (contains invalid character)
func testAccResourceIamResourceTagsConfig_invalidKey(serviceName string) string {
	return fmt.Sprintf(`
data "ovh_cloud_project" "project" {
  service_name = "%s"
}

resource "ovh_iam_resource_tags" "test" {
  urn = data.ovh_cloud_project.project.iam.urn

  tags = {
    "invalid key!" = "value"
  }
}
`, serviceName)
}

// testAccResourceIamResourceTagsConfig_invalidValue returns a configuration with an invalid tag value (contains invalid character)
func testAccResourceIamResourceTagsConfig_invalidValue(serviceName string) string {
	return fmt.Sprintf(`
data "ovh_cloud_project" "project" {
  service_name = "%s"
}

resource "ovh_iam_resource_tags" "test" {
  urn = data.ovh_cloud_project.project.iam.urn

  tags = {
    "%s-key" = "invalid value!"
  }
}
`, serviceName, test_prefix)
}

// testAccResourceIamResourceTagsConfig_tooLongKey returns a configuration with a tag key that exceeds 128 characters
func testAccResourceIamResourceTagsConfig_tooLongKey(serviceName string) string {
	// Create a key that is 129 characters long (exceeds the 128 character limit)
	longKey := strings.Repeat("a", 129)
	return fmt.Sprintf(`
data "ovh_cloud_project" "project" {
  service_name = "%s"
}

resource "ovh_iam_resource_tags" "test" {
  urn = data.ovh_cloud_project.project.iam.urn

  tags = {
    "%s" = "value"
  }
}
`, serviceName, longKey)
}

// testAccResourceIamResourceTagsConfig_tooLongValue returns a configuration with a tag value that exceeds 256 characters
func testAccResourceIamResourceTagsConfig_tooLongValue(serviceName string) string {
	// Create a value that is 257 characters long (exceeds the 256 character limit)
	longValue := strings.Repeat("a", 257)
	return fmt.Sprintf(`
data "ovh_cloud_project" "project" {
  service_name = "%s"
}

resource "ovh_iam_resource_tags" "test" {
  urn = data.ovh_cloud_project.project.iam.urn

  tags = {
    "%s-key" = "%s"
  }
}
`, serviceName, test_prefix, longValue)
}

// testAccResourceIamResourceTagsConfig_ovhPrefix returns a configuration with a tag key that has ovh: prefix
func testAccResourceIamResourceTagsConfig_ovhPrefix(serviceName string) string {
	return fmt.Sprintf(`
data "ovh_cloud_project" "project" {
  service_name = "%s"
}

resource "ovh_iam_resource_tags" "test" {
  urn = data.ovh_cloud_project.project.iam.urn

  tags = {
    "ovh:managed" = "true"
  }
}
`, serviceName)
}
