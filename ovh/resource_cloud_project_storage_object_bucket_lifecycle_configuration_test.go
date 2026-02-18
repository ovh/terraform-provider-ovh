package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCloudProjectStorageLifecycleConfiguration_basic(t *testing.T) {
	bucketName := acctest.RandomWithPrefix(test_prefix)
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	region := "GRA"

	configBucket := fmt.Sprintf(`
	resource "ovh_cloud_project_storage" "bucket" {
		service_name = "%s"
		region_name  = "%s"
		name         = "%s"
	}
	`, serviceName, region, bucketName)

	configLifecycle := configBucket + fmt.Sprintf(`
	resource "ovh_cloud_project_storage_object_bucket_lifecycle_configuration" "lifecycle" {
		service_name   = ovh_cloud_project_storage.bucket.service_name
		region_name    = ovh_cloud_project_storage.bucket.region_name
		container_name = ovh_cloud_project_storage.bucket.name

		rules = [
			{
				id     = "expiration-rule"
				status = "enabled"
				expiration = {
					days = 30
				}
			}
		]
	}
	`)

	configLifecycleUpdated := configBucket + fmt.Sprintf(`
	resource "ovh_cloud_project_storage_object_bucket_lifecycle_configuration" "lifecycle" {
		service_name   = ovh_cloud_project_storage.bucket.service_name
		region_name    = ovh_cloud_project_storage.bucket.region_name
		container_name = ovh_cloud_project_storage.bucket.name

		rules = [
			{
				id     = "expiration-rule"
				status = "enabled"
				expiration = {
					days = 90
				}
			}
		]
	}
	`)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckCloud(t); testAccCheckCloudProjectExists(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create
			{
				Config: configLifecycle,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_project_storage_object_bucket_lifecycle_configuration.lifecycle", "service_name", serviceName),
					resource.TestCheckResourceAttr("ovh_cloud_project_storage_object_bucket_lifecycle_configuration.lifecycle", "region_name", region),
					resource.TestCheckResourceAttr("ovh_cloud_project_storage_object_bucket_lifecycle_configuration.lifecycle", "container_name", bucketName),
					resource.TestCheckResourceAttr("ovh_cloud_project_storage_object_bucket_lifecycle_configuration.lifecycle", "rules.#", "1"),
					resource.TestCheckResourceAttr("ovh_cloud_project_storage_object_bucket_lifecycle_configuration.lifecycle", "rules.0.id", "expiration-rule"),
					resource.TestCheckResourceAttr("ovh_cloud_project_storage_object_bucket_lifecycle_configuration.lifecycle", "rules.0.status", "enabled"),
					resource.TestCheckResourceAttr("ovh_cloud_project_storage_object_bucket_lifecycle_configuration.lifecycle", "rules.0.expiration.days", "30"),
					resource.TestCheckResourceAttrSet("ovh_cloud_project_storage_object_bucket_lifecycle_configuration.lifecycle", "id"),
				),
			},
			// Update (change expiration days)
			{
				Config: configLifecycleUpdated,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_project_storage_object_bucket_lifecycle_configuration.lifecycle", "rules.0.expiration.days", "90"),
				),
			},
			// ImportState
			{
				ImportState:       true,
				ImportStateVerify: true,
				ResourceName:      "ovh_cloud_project_storage_object_bucket_lifecycle_configuration.lifecycle",
				ImportStateId:     fmt.Sprintf("%s/%s/%s", serviceName, region, bucketName),
			},
		},
	})
}

func TestAccCloudProjectStorageLifecycleConfiguration_multipleRules(t *testing.T) {
	bucketName := acctest.RandomWithPrefix(test_prefix)
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	region := "GRA"

	config := fmt.Sprintf(`
	resource "ovh_cloud_project_storage" "bucket" {
		service_name = "%s"
		region_name  = "%s"
		name         = "%s"
		versioning = {
			status = "enabled"
		}
	}

	resource "ovh_cloud_project_storage_object_bucket_lifecycle_configuration" "lifecycle" {
		service_name   = ovh_cloud_project_storage.bucket.service_name
		region_name    = ovh_cloud_project_storage.bucket.region_name
		container_name = ovh_cloud_project_storage.bucket.name

		rules = [
			{
				id     = "expiration-rule"
				status = "enabled"
				expiration = {
					days = 365
				}
			},
			{
				id     = "abort-multipart-rule"
				status = "enabled"
				abort_incomplete_multipart_upload = {
					days_after_initiation = 7
				}
			},
			{
				id     = "transition-rule"
				status = "enabled"
				transitions = [
					{
						days          = 90
						storage_class = "STANDARD_IA"
					}
				]
			},
			{
				id     = "noncurrent-expiration-rule"
				status = "enabled"
				noncurrent_version_expiration = {
					noncurrent_days           = 30
					newer_noncurrent_versions = 3
				}
			}
		]
	}
	`, serviceName, region, bucketName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckCloud(t); testAccCheckCloudProjectExists(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_project_storage_object_bucket_lifecycle_configuration.lifecycle", "rules.#", "4"),
					resource.TestCheckResourceAttr("ovh_cloud_project_storage_object_bucket_lifecycle_configuration.lifecycle", "rules.0.id", "expiration-rule"),
					resource.TestCheckResourceAttr("ovh_cloud_project_storage_object_bucket_lifecycle_configuration.lifecycle", "rules.0.expiration.days", "365"),
					resource.TestCheckResourceAttr("ovh_cloud_project_storage_object_bucket_lifecycle_configuration.lifecycle", "rules.1.id", "abort-multipart-rule"),
					resource.TestCheckResourceAttr("ovh_cloud_project_storage_object_bucket_lifecycle_configuration.lifecycle", "rules.1.abort_incomplete_multipart_upload.days_after_initiation", "7"),
					resource.TestCheckResourceAttr("ovh_cloud_project_storage_object_bucket_lifecycle_configuration.lifecycle", "rules.2.id", "transition-rule"),
					resource.TestCheckResourceAttr("ovh_cloud_project_storage_object_bucket_lifecycle_configuration.lifecycle", "rules.2.transitions.0.days", "90"),
					resource.TestCheckResourceAttr("ovh_cloud_project_storage_object_bucket_lifecycle_configuration.lifecycle", "rules.2.transitions.0.storage_class", "STANDARD_IA"),
					resource.TestCheckResourceAttr("ovh_cloud_project_storage_object_bucket_lifecycle_configuration.lifecycle", "rules.3.id", "noncurrent-expiration-rule"),
					resource.TestCheckResourceAttr("ovh_cloud_project_storage_object_bucket_lifecycle_configuration.lifecycle", "rules.3.noncurrent_version_expiration.noncurrent_days", "30"),
					resource.TestCheckResourceAttr("ovh_cloud_project_storage_object_bucket_lifecycle_configuration.lifecycle", "rules.3.noncurrent_version_expiration.newer_noncurrent_versions", "3"),
				),
			},
		},
	})
}

func TestAccCloudProjectStorageLifecycleConfiguration_withFilter(t *testing.T) {
	bucketName := acctest.RandomWithPrefix(test_prefix)
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	region := "GRA"

	config := fmt.Sprintf(`
	resource "ovh_cloud_project_storage" "bucket" {
		service_name = "%s"
		region_name  = "%s"
		name         = "%s"
	}

	resource "ovh_cloud_project_storage_object_bucket_lifecycle_configuration" "lifecycle" {
		service_name   = ovh_cloud_project_storage.bucket.service_name
		region_name    = ovh_cloud_project_storage.bucket.region_name
		container_name = ovh_cloud_project_storage.bucket.name

		rules = [
			{
				id     = "filtered-rule"
				status = "enabled"
				filter = {
					prefix                   = "logs/"
					object_size_greater_than = 1048576
				}
				expiration = {
					days = 60
				}
			}
		]
	}
	`, serviceName, region, bucketName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckCloud(t); testAccCheckCloudProjectExists(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_project_storage_object_bucket_lifecycle_configuration.lifecycle", "rules.#", "1"),
					resource.TestCheckResourceAttr("ovh_cloud_project_storage_object_bucket_lifecycle_configuration.lifecycle", "rules.0.filter.prefix", "logs/"),
					resource.TestCheckResourceAttr("ovh_cloud_project_storage_object_bucket_lifecycle_configuration.lifecycle", "rules.0.filter.object_size_greater_than", "1048576"),
					resource.TestCheckResourceAttr("ovh_cloud_project_storage_object_bucket_lifecycle_configuration.lifecycle", "rules.0.expiration.days", "60"),
				),
			},
		},
	})
}

// TestAccCloudProjectStorageLifecycleConfiguration_noncurrentVersionTransitions tests the
// noncurrent_version_transitions block which was entirely absent from previous tests.
func TestAccCloudProjectStorageLifecycleConfiguration_noncurrentVersionTransitions(t *testing.T) {
	bucketName := acctest.RandomWithPrefix(test_prefix)
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	region := "GRA"

	config := fmt.Sprintf(`
	resource "ovh_cloud_project_storage" "bucket" {
		service_name = "%s"
		region_name  = "%s"
		name         = "%s"
		versioning = {
			status = "enabled"
		}
	}

	resource "ovh_cloud_project_storage_object_bucket_lifecycle_configuration" "lifecycle" {
		service_name   = ovh_cloud_project_storage.bucket.service_name
		region_name    = ovh_cloud_project_storage.bucket.region_name
		container_name = ovh_cloud_project_storage.bucket.name

		rules = [
			{
				id     = "noncurrent-transition-rule"
				status = "enabled"
				noncurrent_version_transitions = [
					{
						noncurrent_days           = 30
						newer_noncurrent_versions = 2
						storage_class             = "STANDARD_IA"
					}
				]
			}
		]
	}
	`, serviceName, region, bucketName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckCloud(t); testAccCheckCloudProjectExists(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_project_storage_object_bucket_lifecycle_configuration.lifecycle", "rules.#", "1"),
					resource.TestCheckResourceAttr("ovh_cloud_project_storage_object_bucket_lifecycle_configuration.lifecycle", "rules.0.id", "noncurrent-transition-rule"),
					resource.TestCheckResourceAttr("ovh_cloud_project_storage_object_bucket_lifecycle_configuration.lifecycle", "rules.0.noncurrent_version_transitions.#", "1"),
					resource.TestCheckResourceAttr("ovh_cloud_project_storage_object_bucket_lifecycle_configuration.lifecycle", "rules.0.noncurrent_version_transitions.0.noncurrent_days", "30"),
					resource.TestCheckResourceAttr("ovh_cloud_project_storage_object_bucket_lifecycle_configuration.lifecycle", "rules.0.noncurrent_version_transitions.0.newer_noncurrent_versions", "2"),
					resource.TestCheckResourceAttr("ovh_cloud_project_storage_object_bucket_lifecycle_configuration.lifecycle", "rules.0.noncurrent_version_transitions.0.storage_class", "STANDARD_IA"),
				),
			},
		},
	})
}

// TestAccCloudProjectStorageLifecycleConfiguration_expirationByDate tests expiration.date and
// transitions[].date (ISO 8601 date-based triggers instead of day counts).
func TestAccCloudProjectStorageLifecycleConfiguration_expirationByDate(t *testing.T) {
	bucketName := acctest.RandomWithPrefix(test_prefix)
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	region := "GRA"

	config := fmt.Sprintf(`
	resource "ovh_cloud_project_storage" "bucket" {
		service_name = "%s"
		region_name  = "%s"
		name         = "%s"
	}

	resource "ovh_cloud_project_storage_object_bucket_lifecycle_configuration" "lifecycle" {
		service_name   = ovh_cloud_project_storage.bucket.service_name
		region_name    = ovh_cloud_project_storage.bucket.region_name
		container_name = ovh_cloud_project_storage.bucket.name

		rules = [
			{
				id     = "date-expiration-rule"
				status = "enabled"
				expiration = {
					date = "2030-01-01"
				}
			},
			{
				id     = "date-transition-rule"
				status = "enabled"
				transitions = [
					{
						date          = "2027-06-01"
						storage_class = "STANDARD_IA"
					}
				]
			}
		]
	}
	`, serviceName, region, bucketName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckCloud(t); testAccCheckCloudProjectExists(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_project_storage_object_bucket_lifecycle_configuration.lifecycle", "rules.#", "2"),
					resource.TestCheckResourceAttr("ovh_cloud_project_storage_object_bucket_lifecycle_configuration.lifecycle", "rules.0.id", "date-expiration-rule"),
					resource.TestCheckResourceAttr("ovh_cloud_project_storage_object_bucket_lifecycle_configuration.lifecycle", "rules.0.expiration.date", "2030-01-01"),
					resource.TestCheckResourceAttr("ovh_cloud_project_storage_object_bucket_lifecycle_configuration.lifecycle", "rules.1.id", "date-transition-rule"),
					resource.TestCheckResourceAttr("ovh_cloud_project_storage_object_bucket_lifecycle_configuration.lifecycle", "rules.1.transitions.0.date", "2027-06-01"),
					resource.TestCheckResourceAttr("ovh_cloud_project_storage_object_bucket_lifecycle_configuration.lifecycle", "rules.1.transitions.0.storage_class", "STANDARD_IA"),
				),
			},
		},
	})
}

// TestAccCloudProjectStorageLifecycleConfiguration_expiredObjectDeleteMarker tests
// expiration.expired_object_delete_marker which removes delete markers with no noncurrent versions.
func TestAccCloudProjectStorageLifecycleConfiguration_expiredObjectDeleteMarker(t *testing.T) {
	bucketName := acctest.RandomWithPrefix(test_prefix)
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	region := "GRA"

	config := fmt.Sprintf(`
	resource "ovh_cloud_project_storage" "bucket" {
		service_name = "%s"
		region_name  = "%s"
		name         = "%s"
		versioning = {
			status = "enabled"
		}
	}

	resource "ovh_cloud_project_storage_object_bucket_lifecycle_configuration" "lifecycle" {
		service_name   = ovh_cloud_project_storage.bucket.service_name
		region_name    = ovh_cloud_project_storage.bucket.region_name
		container_name = ovh_cloud_project_storage.bucket.name

		rules = [
			{
				id     = "delete-marker-rule"
				status = "enabled"
				expiration = {
					expired_object_delete_marker = true
				}
			}
		]
	}
	`, serviceName, region, bucketName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckCloud(t); testAccCheckCloudProjectExists(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_project_storage_object_bucket_lifecycle_configuration.lifecycle", "rules.#", "1"),
					resource.TestCheckResourceAttr("ovh_cloud_project_storage_object_bucket_lifecycle_configuration.lifecycle", "rules.0.id", "delete-marker-rule"),
					resource.TestCheckResourceAttr("ovh_cloud_project_storage_object_bucket_lifecycle_configuration.lifecycle", "rules.0.expiration.expired_object_delete_marker", "true"),
				),
			},
		},
	})
}

// TestAccCloudProjectStorageLifecycleConfiguration_fullFilter tests filter fields not yet covered:
// object_size_less_than and tags.
func TestAccCloudProjectStorageLifecycleConfiguration_fullFilter(t *testing.T) {
	bucketName := acctest.RandomWithPrefix(test_prefix)
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	region := "GRA"

	config := fmt.Sprintf(`
	resource "ovh_cloud_project_storage" "bucket" {
		service_name = "%s"
		region_name  = "%s"
		name         = "%s"
	}

	resource "ovh_cloud_project_storage_object_bucket_lifecycle_configuration" "lifecycle" {
		service_name   = ovh_cloud_project_storage.bucket.service_name
		region_name    = ovh_cloud_project_storage.bucket.region_name
		container_name = ovh_cloud_project_storage.bucket.name

		rules = [
			{
				id     = "size-filter-rule"
				status = "enabled"
				filter = {
					object_size_less_than = 524288
				}
				expiration = {
					days = 180
				}
			},
			{
				id     = "tag-filter-rule"
				status = "enabled"
				filter = {
					tags = {
						env  = "staging"
						team = "infra"
					}
				}
				expiration = {
					days = 90
				}
			}
		]
	}
	`, serviceName, region, bucketName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckCloud(t); testAccCheckCloudProjectExists(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_project_storage_object_bucket_lifecycle_configuration.lifecycle", "rules.#", "2"),
					resource.TestCheckResourceAttr("ovh_cloud_project_storage_object_bucket_lifecycle_configuration.lifecycle", "rules.0.filter.object_size_less_than", "524288"),
					resource.TestCheckResourceAttr("ovh_cloud_project_storage_object_bucket_lifecycle_configuration.lifecycle", "rules.1.filter.tags.env", "staging"),
					resource.TestCheckResourceAttr("ovh_cloud_project_storage_object_bucket_lifecycle_configuration.lifecycle", "rules.1.filter.tags.team", "infra"),
				),
			},
		},
	})
}

// TestAccCloudProjectStorageLifecycleConfiguration_disabledRule verifies that a rule with
// status = "disabled" is correctly stored and returned.
func TestAccCloudProjectStorageLifecycleConfiguration_disabledRule(t *testing.T) {
	bucketName := acctest.RandomWithPrefix(test_prefix)
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	region := "GRA"

	config := fmt.Sprintf(`
	resource "ovh_cloud_project_storage" "bucket" {
		service_name = "%s"
		region_name  = "%s"
		name         = "%s"
	}

	resource "ovh_cloud_project_storage_object_bucket_lifecycle_configuration" "lifecycle" {
		service_name   = ovh_cloud_project_storage.bucket.service_name
		region_name    = ovh_cloud_project_storage.bucket.region_name
		container_name = ovh_cloud_project_storage.bucket.name

		rules = [
			{
				id     = "disabled-rule"
				status = "disabled"
				expiration = {
					days = 30
				}
			}
		]
	}
	`, serviceName, region, bucketName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckCloud(t); testAccCheckCloudProjectExists(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_project_storage_object_bucket_lifecycle_configuration.lifecycle", "rules.#", "1"),
					resource.TestCheckResourceAttr("ovh_cloud_project_storage_object_bucket_lifecycle_configuration.lifecycle", "rules.0.id", "disabled-rule"),
					resource.TestCheckResourceAttr("ovh_cloud_project_storage_object_bucket_lifecycle_configuration.lifecycle", "rules.0.status", "disabled"),
				),
			},
		},
	})
}

// TestAccCloudProjectStorageLifecycleConfiguration_multipleTransitions tests multiple rules each
// with a transition targeting STANDARD_IA on different prefixes.
func TestAccCloudProjectStorageLifecycleConfiguration_multipleTransitions(t *testing.T) {
	bucketName := acctest.RandomWithPrefix(test_prefix)
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	region := "GRA"

	config := fmt.Sprintf(`
	resource "ovh_cloud_project_storage" "bucket" {
		service_name = "%s"
		region_name  = "%s"
		name         = "%s"
	}

	resource "ovh_cloud_project_storage_object_bucket_lifecycle_configuration" "lifecycle" {
		service_name   = ovh_cloud_project_storage.bucket.service_name
		region_name    = ovh_cloud_project_storage.bucket.region_name
		container_name = ovh_cloud_project_storage.bucket.name

		rules = [
			{
				id     = "rule-archive-prefix"
				status = "enabled"
				filter = { prefix = "archive/" }
				transitions = [
					{
						days          = 30
						storage_class = "STANDARD_IA"
					}
				]
			},
			{
				id     = "rule-logs-prefix"
				status = "enabled"
				filter = { prefix = "logs/" }
				transitions = [
					{
						days          = 60
						storage_class = "STANDARD_IA"
					}
				]
			}
		]
	}
	`, serviceName, region, bucketName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckCloud(t); testAccCheckCloudProjectExists(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_project_storage_object_bucket_lifecycle_configuration.lifecycle", "rules.#", "2"),
					resource.TestCheckResourceAttr("ovh_cloud_project_storage_object_bucket_lifecycle_configuration.lifecycle", "rules.0.id", "rule-archive-prefix"),
					resource.TestCheckResourceAttr("ovh_cloud_project_storage_object_bucket_lifecycle_configuration.lifecycle", "rules.0.transitions.#", "1"),
					resource.TestCheckResourceAttr("ovh_cloud_project_storage_object_bucket_lifecycle_configuration.lifecycle", "rules.0.transitions.0.days", "30"),
					resource.TestCheckResourceAttr("ovh_cloud_project_storage_object_bucket_lifecycle_configuration.lifecycle", "rules.0.transitions.0.storage_class", "STANDARD_IA"),
					resource.TestCheckResourceAttr("ovh_cloud_project_storage_object_bucket_lifecycle_configuration.lifecycle", "rules.1.id", "rule-logs-prefix"),
					resource.TestCheckResourceAttr("ovh_cloud_project_storage_object_bucket_lifecycle_configuration.lifecycle", "rules.1.transitions.#", "1"),
					resource.TestCheckResourceAttr("ovh_cloud_project_storage_object_bucket_lifecycle_configuration.lifecycle", "rules.1.transitions.0.days", "60"),
					resource.TestCheckResourceAttr("ovh_cloud_project_storage_object_bucket_lifecycle_configuration.lifecycle", "rules.1.transitions.0.storage_class", "STANDARD_IA"),
				),
			},
		},
	})
}

// TestAccCloudProjectStorageLifecycleConfiguration_noServiceName validates that when service_name
// is not set in the resource configuration, the provider falls back to the OVH_CLOUD_PROJECT_SERVICE
// environment variable.
func TestAccCloudProjectStorageLifecycleConfiguration_noServiceName(t *testing.T) {
	bucketName := acctest.RandomWithPrefix(test_prefix)
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	region := "GRA"

	config := fmt.Sprintf(`
	resource "ovh_cloud_project_storage" "bucket" {
		service_name = "%s"
		region_name  = "%s"
		name         = "%s"
	}

	resource "ovh_cloud_project_storage_object_bucket_lifecycle_configuration" "lifecycle" {
		region_name    = ovh_cloud_project_storage.bucket.region_name
		container_name = ovh_cloud_project_storage.bucket.name

		rules = [
			{
				id     = "expiration-rule"
				status = "enabled"
				expiration = {
					days = 30
				}
			}
		]
	}
	`, serviceName, region, bucketName)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCloud(t)
			testAccCheckCloudProjectExists(t)
			checkEnvOrSkip(t, "OVH_CLOUD_PROJECT_SERVICE")
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_project_storage_object_bucket_lifecycle_configuration.lifecycle", "service_name", os.Getenv("OVH_CLOUD_PROJECT_SERVICE")),
					resource.TestCheckResourceAttr("ovh_cloud_project_storage_object_bucket_lifecycle_configuration.lifecycle", "region_name", region),
					resource.TestCheckResourceAttr("ovh_cloud_project_storage_object_bucket_lifecycle_configuration.lifecycle", "container_name", bucketName),
					resource.TestCheckResourceAttr("ovh_cloud_project_storage_object_bucket_lifecycle_configuration.lifecycle", "rules.#", "1"),
					resource.TestCheckResourceAttr("ovh_cloud_project_storage_object_bucket_lifecycle_configuration.lifecycle", "rules.0.id", "expiration-rule"),
					resource.TestCheckResourceAttr("ovh_cloud_project_storage_object_bucket_lifecycle_configuration.lifecycle", "rules.0.status", "enabled"),
					resource.TestCheckResourceAttr("ovh_cloud_project_storage_object_bucket_lifecycle_configuration.lifecycle", "rules.0.expiration.days", "30"),
					resource.TestCheckResourceAttrSet("ovh_cloud_project_storage_object_bucket_lifecycle_configuration.lifecycle", "id"),
				),
			},
		},
	})
}
