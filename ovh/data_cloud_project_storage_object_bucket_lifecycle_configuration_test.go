package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCloudProjectStorageLifecycleConfigurationDataSource_basic(t *testing.T) {
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
				id     = "expiration-rule"
				status = "enabled"
				expiration = {
					days = 30
				}
			},
			{
				id     = "abort-multipart-rule"
				status = "enabled"
				abort_incomplete_multipart_upload = {
					days_after_initiation = 7
				}
			}
		]
	}

	data "ovh_cloud_project_storage_object_bucket_lifecycle_configuration" "lifecycle" {
		service_name   = ovh_cloud_project_storage_object_bucket_lifecycle_configuration.lifecycle.service_name
		region_name    = ovh_cloud_project_storage_object_bucket_lifecycle_configuration.lifecycle.region_name
		container_name = ovh_cloud_project_storage_object_bucket_lifecycle_configuration.lifecycle.container_name
	}
	`, serviceName, region, bucketName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckCloud(t); testAccCheckCloudProjectExists(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.ovh_cloud_project_storage_object_bucket_lifecycle_configuration.lifecycle", "service_name", serviceName),
					resource.TestCheckResourceAttr("data.ovh_cloud_project_storage_object_bucket_lifecycle_configuration.lifecycle", "region_name", region),
					resource.TestCheckResourceAttr("data.ovh_cloud_project_storage_object_bucket_lifecycle_configuration.lifecycle", "container_name", bucketName),
					resource.TestCheckResourceAttrSet("data.ovh_cloud_project_storage_object_bucket_lifecycle_configuration.lifecycle", "id"),
					resource.TestCheckResourceAttr("data.ovh_cloud_project_storage_object_bucket_lifecycle_configuration.lifecycle", "rules.#", "2"),
					resource.TestCheckResourceAttr("data.ovh_cloud_project_storage_object_bucket_lifecycle_configuration.lifecycle", "rules.0.id", "expiration-rule"),
					resource.TestCheckResourceAttr("data.ovh_cloud_project_storage_object_bucket_lifecycle_configuration.lifecycle", "rules.0.status", "enabled"),
					resource.TestCheckResourceAttr("data.ovh_cloud_project_storage_object_bucket_lifecycle_configuration.lifecycle", "rules.0.expiration.days", "30"),
					resource.TestCheckResourceAttr("data.ovh_cloud_project_storage_object_bucket_lifecycle_configuration.lifecycle", "rules.1.id", "abort-multipart-rule"),
					resource.TestCheckResourceAttr("data.ovh_cloud_project_storage_object_bucket_lifecycle_configuration.lifecycle", "rules.1.abort_incomplete_multipart_upload.days_after_initiation", "7"),
				),
			},
		},
	})
}

// TestAccCloudProjectStorageLifecycleConfigurationDataSource_noServiceName validates that when
// service_name is not set in the data source configuration, the provider falls back to the
// OVH_CLOUD_PROJECT_SERVICE environment variable.
func TestAccCloudProjectStorageLifecycleConfigurationDataSource_noServiceName(t *testing.T) {
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
				id     = "expiration-rule"
				status = "enabled"
				expiration = {
					days = 30
				}
			}
		]
	}

	data "ovh_cloud_project_storage_object_bucket_lifecycle_configuration" "lifecycle" {
		region_name    = ovh_cloud_project_storage_object_bucket_lifecycle_configuration.lifecycle.region_name
		container_name = ovh_cloud_project_storage_object_bucket_lifecycle_configuration.lifecycle.container_name
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
					resource.TestCheckResourceAttr("data.ovh_cloud_project_storage_object_bucket_lifecycle_configuration.lifecycle", "service_name", os.Getenv("OVH_CLOUD_PROJECT_SERVICE")),
					resource.TestCheckResourceAttr("data.ovh_cloud_project_storage_object_bucket_lifecycle_configuration.lifecycle", "region_name", region),
					resource.TestCheckResourceAttr("data.ovh_cloud_project_storage_object_bucket_lifecycle_configuration.lifecycle", "container_name", bucketName),
					resource.TestCheckResourceAttrSet("data.ovh_cloud_project_storage_object_bucket_lifecycle_configuration.lifecycle", "id"),
					resource.TestCheckResourceAttr("data.ovh_cloud_project_storage_object_bucket_lifecycle_configuration.lifecycle", "rules.#", "1"),
					resource.TestCheckResourceAttr("data.ovh_cloud_project_storage_object_bucket_lifecycle_configuration.lifecycle", "rules.0.id", "expiration-rule"),
					resource.TestCheckResourceAttr("data.ovh_cloud_project_storage_object_bucket_lifecycle_configuration.lifecycle", "rules.0.status", "enabled"),
					resource.TestCheckResourceAttr("data.ovh_cloud_project_storage_object_bucket_lifecycle_configuration.lifecycle", "rules.0.expiration.days", "30"),
				),
			},
		},
	})
}
