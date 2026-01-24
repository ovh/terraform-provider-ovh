package ovh

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccCloudProjectStorageReplicationJob_basic(t *testing.T) {
	sourceBucketName := fmt.Sprintf("test-src-%d", acctest.RandIntRange(1000, 9999))
	replicaBucketName := fmt.Sprintf("test-rep-%d", acctest.RandIntRange(1000, 9999))
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	sourceRegion := "DE"
	replicaRegion := "GRA"

	configBuckets := fmt.Sprintf(`
	resource "ovh_cloud_project_storage" "replica" {
		service_name = "%s"
		region_name  = "%s"
		name         = "%s"
		versioning = {
			status = "enabled"
		}
	}

	resource "ovh_cloud_project_storage" "source" {
		service_name = "%s"
		region_name  = "%s"
		name         = "%s"
		versioning = {
			status = "enabled"
		}
		replication = {
			rules = [
				{
					id       = "replication-rule-1"
					priority = 1
					status   = "enabled"
					destination = {
						name                           = ovh_cloud_project_storage.replica.name
						region                         = ovh_cloud_project_storage.replica.region
						remove_on_main_bucket_deletion = false
					}
					delete_marker_replication = "disabled"
				}
			]
		}
	}
	`, serviceName, replicaRegion, replicaBucketName, serviceName, sourceRegion, sourceBucketName)

	configWithJob := configBuckets + fmt.Sprintf(`
	resource "ovh_cloud_project_storage_replication_job" "job" {
		service_name   = ovh_cloud_project_storage.source.service_name
		region_name    = ovh_cloud_project_storage.source.region_name
		container_name = ovh_cloud_project_storage.source.name

		depends_on = [ovh_cloud_project_storage.source]
	}
	`)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckCloud(t); testAccCheckCloudProjectExists(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: configBuckets,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_project_storage.source", "service_name", serviceName),
					resource.TestCheckResourceAttr("ovh_cloud_project_storage.replica", "service_name", serviceName),
					// Wait for replication configuration to propagate
					func(s *terraform.State) error {
						time.Sleep(15 * time.Second)
						return nil
					},
				),
			},
			{
				Config: configWithJob,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_project_storage_replication_job.job", "service_name", serviceName),
					resource.TestCheckResourceAttr("ovh_cloud_project_storage_replication_job.job", "region_name", sourceRegion),
					resource.TestCheckResourceAttr("ovh_cloud_project_storage_replication_job.job", "container_name", sourceBucketName),
					resource.TestCheckResourceAttrSet("ovh_cloud_project_storage_replication_job.job", "id"),
				),
			},
		},
	})
}
