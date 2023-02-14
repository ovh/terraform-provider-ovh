package ovh

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const WORKFLOW_BACKUP_TEST_CONF = `
resource "ovh_cloud_project_workflow_backup" "my_backup"{
	service_name		= "%s"
	region_name				= "%s"
	cron				= "50 4 * * *"
	instance_id			= "%s"
	max_execution_count	= "0"
	name				= "%s"
	rotation			= "7"
}
`
const WORKFLOW_BACKUP_TEST_CONF_UPDATED = `
resource "ovh_cloud_project_workflow_backup" "my_backup"{
	service_name		= "%s"
	region_name				= "%s"
	cron				= "50 4 * * *"
	instance_id			= "%s"
	max_execution_count	= "0"
	name				= "%s"
	rotation			= "5"
}
`

const CLOUD_PROJECT_TEST_ENV_VAR = "OVH_CLOUD_PROJECT_SERVICE_TEST"
const WORKFLOW_BACKUP_TEST_REGION_ENV_VAR = "OVH_CLOUD_PROJECT_WORKFLOW_BACKUP_REGION_TEST"
const WORKFLOW_BACKUP_TEST_INSTANCE_ID_ENV_VAR = "OVH_CLOUD_PROJECT_WORKFLOW_BACKUP_INSTANCE_ID_TEST"
const WORKFLOW_BACKUP_RESOURCE = "ovh_cloud_project_workflow_backup"
const WORKFLOW_BACKUP_RESOURCE_NAME = "ovh_cloud_project_workflow_backup.my_backup"

func init() {
	resource.AddTestSweepers(WORKFLOW_BACKUP_RESOURCE, &resource.Sweeper{
		Name: WORKFLOW_BACKUP_RESOURCE,
		F:    testSweepCloudProjectWorkflowBackup,
	})
}

func testSweepCloudProjectWorkflowBackup(region string) error {
	client, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting client: %s", err)
	}

	region = os.Getenv(WORKFLOW_BACKUP_TEST_REGION_ENV_VAR)
	if region == "" {
		log.Printf("[DEBUG] No region in env variable => no sweeping")
		return nil
	}
	wfToSweep := make([]CloudProjectWorkflowBackupResponse, 0)
	var endpoint = fmt.Sprintf(OVH_CLOUD_PROJECT_WORKFLOW_BACKUP_ENDPOINT, os.Getenv(CLOUD_PROJECT_TEST_ENV_VAR),
		region)

	if err := client.Get(endpoint, &wfToSweep); err != nil {
		return fmt.Errorf("Error calling GET %s:\n\t %q", endpoint, err)
	}

	for _, wf := range wfToSweep {
		var deleteEndpoint = endpoint + "/%s"
		client.Delete(fmt.Sprintf(deleteEndpoint, wf.Id), nil)
		if err != nil {
			return fmt.Errorf("Error while deleting workflow with id %s through endpoint %s", wf.Id, deleteEndpoint)
		}
		log.Printf("[DEBUG] workflow %s is deleted", wf.Id)
	}
	return nil
}

func TestAccCloudProjectWorkflowBackup(t *testing.T) {
	name := acctest.RandomWithPrefix(test_prefix)
	var serviceName = os.Getenv(CLOUD_PROJECT_TEST_ENV_VAR)
	var instanceId = os.Getenv(WORKFLOW_BACKUP_TEST_INSTANCE_ID_ENV_VAR)
	var config = fmt.Sprintf(WORKFLOW_BACKUP_TEST_CONF,
		serviceName,
		os.Getenv(WORKFLOW_BACKUP_TEST_REGION_ENV_VAR),
		instanceId,
		name,
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCloud(t)
			testAccCheckCloudProjectExists(t)
			testAccPreCheckWorkflowBackup(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(WORKFLOW_BACKUP_RESOURCE_NAME, OVH_CLOUD_PROJECT_WORKFLOW_BACKUP_SERVICE, serviceName),
					resource.TestCheckResourceAttr(WORKFLOW_BACKUP_RESOURCE_NAME, "name", name),
					resource.TestCheckResourceAttr(WORKFLOW_BACKUP_RESOURCE_NAME, "instance_id", instanceId),
				),
			},
			{
				Config: fmt.Sprintf(WORKFLOW_BACKUP_TEST_CONF_UPDATED,
					serviceName,
					os.Getenv(WORKFLOW_BACKUP_TEST_REGION_ENV_VAR),
					instanceId,
					name,
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(WORKFLOW_BACKUP_RESOURCE_NAME, OVH_CLOUD_PROJECT_WORKFLOW_BACKUP_SERVICE, serviceName),
					resource.TestCheckResourceAttr(WORKFLOW_BACKUP_RESOURCE_NAME, "name", name),
					resource.TestCheckResourceAttr(WORKFLOW_BACKUP_RESOURCE_NAME, "instance_id", instanceId),
				),
			},
			{
				ResourceName:      WORKFLOW_BACKUP_RESOURCE_NAME,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
