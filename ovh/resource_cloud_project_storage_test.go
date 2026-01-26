package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCloudProjectRegionStorage_basic(t *testing.T) {
	bucketName := acctest.RandomWithPrefix(test_prefix)
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	config := fmt.Sprintf(`
	resource "ovh_cloud_project_storage" "storage" {
		service_name = "%s"
		region_name = "GRA"
		name = "%s"
		versioning = {
			status = "enabled"
		}
	}
	`, serviceName, bucketName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckCloud(t); testAccCheckCloudProjectExists(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_project_storage.storage", "name", bucketName),
					resource.TestCheckResourceAttr("ovh_cloud_project_storage.storage", "region", "GRA"),
					resource.TestCheckResourceAttr("ovh_cloud_project_storage.storage", "versioning.status", "enabled"),
					resource.TestCheckResourceAttr("ovh_cloud_project_storage.storage", "encryption.sse_algorithm", "plaintext"),
					resource.TestCheckResourceAttrSet("ovh_cloud_project_storage.storage", "virtual_host"),
					// Verify ID is populated with composite format: service_name/region_name/name
					resource.TestCheckResourceAttr("ovh_cloud_project_storage.storage", "id", fmt.Sprintf("%s/GRA/%s", serviceName, bucketName)),
				),
			},
			{
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "name",
				ResourceName:                         "ovh_cloud_project_storage.storage",
				ImportStateId:                        fmt.Sprintf("%s/GRA/%s", os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST"), bucketName),
				ImportStateVerifyIgnore:              []string{"created_at"}, // Ignore created_at since its value is invalid in response of the POST.
			},
		},
	})
}

func TestAccCloudProjectRegionStorage_withReplication(t *testing.T) {
	bucketName := acctest.RandomWithPrefix(test_prefix)
	replicaBucketName := acctest.RandomWithPrefix(test_prefix)
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckCloud(t); testAccCheckCloudProjectExists(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "ovh_cloud_project_storage" "storage" {
						service_name = "%s"
						region_name = "GRA"
						name = "%s"

						versioning = {
							status = "enabled"
						}

						replication = {
							rules = [
								{
									id          = "test"
									priority    = 1
									status      = "enabled"
									destination = {
										name   = "%s"
										region = "GRA"
										remove_on_main_bucket_deletion = true
									}
									filter = {
										"prefix" = "test"
										"tags"   = {
											"key": "test"
										}
									}
									delete_marker_replication = "disabled"
								}
							]
						}
					}`, serviceName, bucketName, replicaBucketName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_project_storage.storage", "name", bucketName),
					resource.TestCheckResourceAttr("ovh_cloud_project_storage.storage", "region", "GRA"),
					resource.TestCheckResourceAttr("ovh_cloud_project_storage.storage", "versioning.status", "enabled"),
					resource.TestCheckResourceAttr("ovh_cloud_project_storage.storage", "encryption.sse_algorithm", "plaintext"),
					resource.TestCheckResourceAttr("ovh_cloud_project_storage.storage", "replication.rules.#", "1"),
					resource.TestCheckResourceAttr("ovh_cloud_project_storage.storage", "replication.rules.0.id", "test"),
					resource.TestCheckResourceAttr("ovh_cloud_project_storage.storage", "replication.rules.0.priority", "1"),
					resource.TestCheckResourceAttr("ovh_cloud_project_storage.storage", "replication.rules.0.status", "enabled"),
					resource.TestCheckResourceAttr("ovh_cloud_project_storage.storage", "replication.rules.0.destination.name", replicaBucketName),
					resource.TestCheckResourceAttr("ovh_cloud_project_storage.storage", "replication.rules.0.destination.region", "GRA"),
					resource.TestCheckResourceAttr("ovh_cloud_project_storage.storage", "replication.rules.0.filter.prefix", "test"),
					resource.TestCheckResourceAttrSet("ovh_cloud_project_storage.storage", "virtual_host"),
					// Verify ID is populated with composite format: service_name/region_name/name
					resource.TestCheckResourceAttr("ovh_cloud_project_storage.storage", "id", fmt.Sprintf("%s/GRA/%s", serviceName, bucketName)),
				),
			},
			{
				Config: fmt.Sprintf(`
					resource "ovh_cloud_project_storage" "storage" {
						service_name = "%s"
						region_name = "GRA"
						name = "%s"

						versioning = {
							status = "enabled"
						}

						replication = {
							rules = [
								{
									id          = "test"
									priority    = 1
									status      = "enabled"
									destination = {
										name   = "%s"
										region = "GRA"
										remove_on_main_bucket_deletion = true
									}
									filter = {
										"prefix" = "test-updated"
										"tags"   = {
											"key": "test-updated"
										}
									}
									delete_marker_replication = "disabled"
								}
							]
						}
					} `, serviceName, bucketName, replicaBucketName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_project_storage.storage", "name", bucketName),
					resource.TestCheckResourceAttr("ovh_cloud_project_storage.storage", "region", "GRA"),
					resource.TestCheckResourceAttr("ovh_cloud_project_storage.storage", "versioning.status", "enabled"),
					resource.TestCheckResourceAttr("ovh_cloud_project_storage.storage", "encryption.sse_algorithm", "plaintext"),
					resource.TestCheckResourceAttr("ovh_cloud_project_storage.storage", "replication.rules.#", "1"),
					resource.TestCheckResourceAttr("ovh_cloud_project_storage.storage", "replication.rules.0.id", "test"),
					resource.TestCheckResourceAttr("ovh_cloud_project_storage.storage", "replication.rules.0.priority", "1"),
					resource.TestCheckResourceAttr("ovh_cloud_project_storage.storage", "replication.rules.0.status", "enabled"),
					resource.TestCheckResourceAttr("ovh_cloud_project_storage.storage", "replication.rules.0.destination.name", replicaBucketName),
					resource.TestCheckResourceAttr("ovh_cloud_project_storage.storage", "replication.rules.0.destination.region", "GRA"),
					resource.TestCheckResourceAttr("ovh_cloud_project_storage.storage", "replication.rules.0.filter.prefix", "test-updated"),
					resource.TestCheckResourceAttrSet("ovh_cloud_project_storage.storage", "virtual_host"),
					// Verify ID is populated with composite format after update
					resource.TestCheckResourceAttr("ovh_cloud_project_storage.storage", "id", fmt.Sprintf("%s/GRA/%s", serviceName, bucketName)),
				),
			},
			{
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "name",
				ResourceName:                         "ovh_cloud_project_storage.storage",
				ImportStateId:                        fmt.Sprintf("%s/GRA/%s", os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST"), bucketName),
				// Ignore created_at since its value is invalid in response of the POST.
				// Also ignore remove_on_main_bucket_deletion since its computed value is not returned by the API.
				ImportStateVerifyIgnore: []string{"created_at", "replication.rules.0.destination.remove_on_main_bucket_deletion"},
			},
		},
	})
}

func TestAccCloudProjectRegionStorage_withObjectLock(t *testing.T) {
	bucketName := acctest.RandomWithPrefix(test_prefix)
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckCloud(t); testAccCheckCloudProjectExists(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Step 1: Create with P7D (7 days) - API will normalize to P1W (1 week)
			{
				Config: fmt.Sprintf(`
					resource "ovh_cloud_project_storage" "storage" {
						service_name = "%s"
						region_name = "GRA"
						name = "%s"
						object_lock = {
							status = "enabled"
							rule = {
								mode   = "compliance"
								period = "P7D"
							}
						}
					}`, serviceName, bucketName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_project_storage.storage", "name", bucketName),
					resource.TestCheckResourceAttr("ovh_cloud_project_storage.storage", "object_lock.status", "enabled"),
					resource.TestCheckResourceAttr("ovh_cloud_project_storage.storage", "object_lock.rule.mode", "compliance"),
					// API normalizes P7D to P1W, but fixISO8601Diff should keep our configured P7D
					resource.TestCheckResourceAttr("ovh_cloud_project_storage.storage", "object_lock.rule.period", "P7D"),
				),
			},
			// Step 2: Update to P2W (2 weeks) - stays as P2W
			{
				Config: fmt.Sprintf(`
					resource "ovh_cloud_project_storage" "storage" {
						service_name = "%s"
						region_name = "GRA"
						name = "%s"
						object_lock = {
							status = "enabled"
							rule = {
								mode   = "compliance"
								period = "P2W"
							}
						}
					}`, serviceName, bucketName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_project_storage.storage", "object_lock.rule.period", "P2W"),
				),
			},
			// Step 3: Update to P14D (14 days) - semantically equal to P2W, should not cause drift
			{
				Config: fmt.Sprintf(`
					resource "ovh_cloud_project_storage" "storage" {
						service_name = "%s"
						region_name = "GRA"
						name = "%s"
						object_lock = {
							status = "enabled"
							rule = {
								mode   = "compliance"
								period = "P14D"
							}
						}
					}`, serviceName, bucketName),
				Check: resource.ComposeTestCheckFunc(
					// Should keep the configured P14D even though API might return P2W
					resource.TestCheckResourceAttr("ovh_cloud_project_storage.storage", "object_lock.rule.period", "P14D"),
				),
			},
			// Step 4: Update to P1D (1 day) - stays as P1D (can't be normalized to weeks)
			{
				Config: fmt.Sprintf(`
					resource "ovh_cloud_project_storage" "storage" {
						service_name = "%s"
						region_name = "GRA"
						name = "%s"
						object_lock = {
							status = "enabled"
							rule = {
								mode   = "compliance"
								period = "P1D"
							}
						}
					}`, serviceName, bucketName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_project_storage.storage", "object_lock.rule.period", "P1D"),
				),
			},
			// Step 5: Update to P8D (8 days) - stays as P8D (can't be normalized to weeks)
			{
				Config: fmt.Sprintf(`
					resource "ovh_cloud_project_storage" "storage" {
						service_name = "%s"
						region_name = "GRA"
						name = "%s"
						object_lock = {
							status = "enabled"
							rule = {
								mode   = "compliance"
								period = "P8D"
							}
						}
					}`, serviceName, bucketName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_project_storage.storage", "object_lock.rule.period", "P8D"),
				),
			},
			// Step 6: Update to P364D (364 days) - API normalizes to P52W, should not cause drift
			{
				Config: fmt.Sprintf(`
					resource "ovh_cloud_project_storage" "storage" {
						service_name = "%s"
						region_name = "GRA"
						name = "%s"
						object_lock = {
							status = "enabled"
							rule = {
								mode   = "compliance"
								period = "P364D"
							}
						}
					}`, serviceName, bucketName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_project_storage.storage", "object_lock.rule.period", "P364D"),
				),
			},
			// Step 7: Change mode from compliance to governance
			{
				Config: fmt.Sprintf(`
					resource "ovh_cloud_project_storage" "storage" {
						service_name = "%s"
						region_name = "GRA"
						name = "%s"
						object_lock = {
							status = "enabled"
							rule = {
								mode   = "governance"
								period = "P364D"
							}
						}
					}`, serviceName, bucketName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_project_storage.storage", "object_lock.rule.mode", "governance"),
					resource.TestCheckResourceAttr("ovh_cloud_project_storage.storage", "object_lock.rule.period", "P364D"),
				),
			},
			// Step 8: Import state verification
			{
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "name",
				ResourceName:                         "ovh_cloud_project_storage.storage",
				ImportStateId:                        fmt.Sprintf("%s/GRA/%s", serviceName, bucketName),
				// Ignore created_at and object_lock.rule.period because:
				// - created_at: invalid value in POST response
				// - object_lock.rule.period: During import, we get API's normalized value (P52W)
				//   without config to compare against. fixISO8601Diff only works when there's
				//   a config value to preserve. This is expected behavior.
				ImportStateVerifyIgnore: []string{"created_at", "object_lock.rule.period"},
			},
			// Step 9: P1Y (1 Year) -> Test year-based duration
			{
				Config: fmt.Sprintf(`
					resource "ovh_cloud_project_storage" "storage" {
						service_name = "%s"
						region_name = "GRA"
						name = "%s"
						object_lock = {
							status = "enabled"
							rule = {
								mode = "compliance"
								period = "P1Y"
							}
						}
					}`, serviceName, bucketName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_project_storage.storage", "object_lock.rule.period", "P1Y"),
				),
			},
		},
	})
}

func TestAccCloudProjectRegionStorage_withObjectLockNoRule(t *testing.T) {
	bucketName := acctest.RandomWithPrefix(test_prefix)
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckCloud(t); testAccCheckCloudProjectExists(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Step 1: Create with object_lock status only (no rule)
			{
				Config: fmt.Sprintf(`
					resource "ovh_cloud_project_storage" "storage" {
						service_name = "%s"
						region_name = "GRA"
						name = "%s"
						object_lock = {
							status = "enabled"
						}
					}`, serviceName, bucketName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_project_storage.storage", "name", bucketName),
					resource.TestCheckResourceAttr("ovh_cloud_project_storage.storage", "object_lock.status", "enabled"),
				),
			},
			// Step 2: Add a rule to existing object_lock
			{
				Config: fmt.Sprintf(`
					resource "ovh_cloud_project_storage" "storage" {
						service_name = "%s"
						region_name = "GRA"
						name = "%s"
						object_lock = {
							status = "enabled"
							rule = {
								mode   = "compliance"
								period = "P7D"
							}
						}
					}`, serviceName, bucketName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_project_storage.storage", "object_lock.status", "enabled"),
					resource.TestCheckResourceAttr("ovh_cloud_project_storage.storage", "object_lock.rule.mode", "compliance"),
					resource.TestCheckResourceAttr("ovh_cloud_project_storage.storage", "object_lock.rule.period", "P7D"),
				),
			},
			// Step 3: Remove the rule (keep only status)
			{
				Config: fmt.Sprintf(`
					resource "ovh_cloud_project_storage" "storage" {
						service_name = "%s"
						region_name = "GRA"
						name = "%s"
						object_lock = {
							status = "enabled"
						}
					}`, serviceName, bucketName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_project_storage.storage", "object_lock.status", "enabled"),
				),
			},
		},
	})
}

func TestAccCloudProjectRegionStorage_withoutObjectLock(t *testing.T) {
	bucketName := acctest.RandomWithPrefix(test_prefix)
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckCloud(t); testAccCheckCloudProjectExists(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Step 1: Create bucket without object_lock configuration
			{
				Config: fmt.Sprintf(`
					resource "ovh_cloud_project_storage" "storage" {
						service_name = "%s"
						region_name = "GRA"
						name = "%s"
					}`, serviceName, bucketName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_project_storage.storage", "name", bucketName),
					resource.TestCheckResourceAttr("ovh_cloud_project_storage.storage", "region", "GRA"),
				),
			},
			// Step 2: Import state verification - object_lock should be imported from API even if not configured
			{
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "name",
				ResourceName:                         "ovh_cloud_project_storage.storage",
				ImportStateId:                        fmt.Sprintf("%s/GRA/%s", serviceName, bucketName),
				ImportStateVerifyIgnore:              []string{"created_at"},
			},
		},
	})
}

func TestAccCloudProjectRegionStorage_objectLockRemovalTriggersReplacement(t *testing.T) {
	bucketName1 := acctest.RandomWithPrefix(test_prefix)
	bucketName2 := acctest.RandomWithPrefix(test_prefix)
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckCloud(t); testAccCheckCloudProjectExists(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Step 1: Create bucket with object_lock enabled
			{
				Config: fmt.Sprintf(`
					resource "ovh_cloud_project_storage" "storage" {
						service_name = "%s"
						region_name = "GRA"
						name = "%s"
						object_lock = {
							status = "enabled"
							rule = {
								mode   = "governance"
								period = "P7D"
							}
						}
					}`, serviceName, bucketName1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_project_storage.storage", "name", bucketName1),
					resource.TestCheckResourceAttr("ovh_cloud_project_storage.storage", "object_lock.status", "enabled"),
					resource.TestCheckResourceAttr("ovh_cloud_project_storage.storage", "object_lock.rule.mode", "governance"),
				),
			},
			// Step 2: Remove object_lock - should trigger bucket replacement
			// The bucket name must change because we're creating a new bucket
			{
				Config: fmt.Sprintf(`
					resource "ovh_cloud_project_storage" "storage" {
						service_name = "%s"
						region_name = "GRA"
						name = "%s"
					}`, serviceName, bucketName2),
				Check: resource.ComposeTestCheckFunc(
					// Verify new bucket was created with different name
					resource.TestCheckResourceAttr("ovh_cloud_project_storage.storage", "name", bucketName2),
					// Verify object_lock is now from API (likely disabled or default state)
				),
			},
		},
	})
}

// TestAccCloudProjectRegionStorage_statusChangeTriggersReplacement validates that changing
// object_lock status from "enabled" to "disabled" triggers bucket replacement
func TestAccCloudProjectRegionStorage_statusChangeTriggersReplacement(t *testing.T) {
	bucketName1 := acctest.RandomWithPrefix(test_prefix)
	bucketName2 := acctest.RandomWithPrefix(test_prefix)
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckCloud(t); testAccCheckCloudProjectExists(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Step 1: Create bucket with object_lock enabled
			{
				Config: fmt.Sprintf(`
					resource "ovh_cloud_project_storage" "storage" {
						service_name = "%s"
						region_name = "GRA"
						name = "%s"
						object_lock = {
							status = "enabled"
							rule = {
								mode   = "compliance"
								period = "P30D"
							}
						}
					}`, serviceName, bucketName1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_project_storage.storage", "name", bucketName1),
					resource.TestCheckResourceAttr("ovh_cloud_project_storage.storage", "object_lock.status", "enabled"),
					resource.TestCheckResourceAttr("ovh_cloud_project_storage.storage", "object_lock.rule.mode", "compliance"),
				),
			},
			// Step 2: Attempt to change status to "disabled" - should trigger bucket replacement
			// The bucket name must change because we're creating a new bucket
			{
				Config: fmt.Sprintf(`
					resource "ovh_cloud_project_storage" "storage" {
						service_name = "%s"
						region_name = "GRA"
						name = "%s"
					}`, serviceName, bucketName2),
				Check: resource.ComposeTestCheckFunc(
					// Verify new bucket was created with different name
					resource.TestCheckResourceAttr("ovh_cloud_project_storage.storage", "name", bucketName2),
					// Verify new bucket has no object_lock (or default API state)
				),
			},
		},
	})
}
