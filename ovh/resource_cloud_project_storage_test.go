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
