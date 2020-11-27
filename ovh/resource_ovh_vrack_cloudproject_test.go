package ovh

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/ovh/go-ovh/ovh"
)

var testAccVrackCloudProjectConfig = fmt.Sprintf(`
resource "ovh_vrack_cloudproject" "vcp" {
  service_name = "%s"
  project_id = "%s"
}
`, os.Getenv("OVH_VRACK_SERVICE_TEST"), os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST"))

var testAccVrackCloudProjectDeprecatedConfig = fmt.Sprintf(`
resource "ovh_vrack_cloudproject" "attach" {
  vrack_id   = "%s"
  project_id = "%s"
}
`, os.Getenv("OVH_VRACK_SERVICE_TEST"), os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST"))

func init() {
	resource.AddTestSweepers("ovh_vrack_cloudproject", &resource.Sweeper{
		Name:         "ovh_vrack_cloudproject",
		Dependencies: []string{"ovh_cloud_network_private"},
		F:            testSweepVrackCloudProject,
	})
}

func testSweepVrackCloudProject(region string) error {
	client, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting client: %s", err)
	}

	vrackId := os.Getenv("OVH_VRACK_SERVICE_TEST")
	if vrackId == "" {
		log.Print("[DEBUG] OVH_VRACK_SERVICE_TEST is not set. No vrack_cloud_project to sweep")
		return nil
	}

	projectId := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	if projectId == "" {
		log.Print("[DEBUG] OVH_CLOUD_PROJECT_SERVICE_TEST is not set. No vrack_cloud_project to sweep")
		return nil
	}

	endpoint := fmt.Sprintf("/vrack/%s/cloudProject/%s",
		url.PathEscape(vrackId),
		url.PathEscape(projectId),
	)

	vcp := &VrackCloudProject{}

	if err := client.Get(endpoint, vcp); err != nil {
		if err.(*ovh.APIError).Code == 404 {
			return nil
		}
		return err
	}

	task := &VrackTask{}

	if err := client.Delete(endpoint, task); err != nil {
		return fmt.Errorf("Error calling DELETE %s with %s/%s:\n\t %q", endpoint, vrackId, projectId, err)
	}

	if err := waitForVrackTask(task, client); err != nil {
		return fmt.Errorf("Error waiting for vrack (%s) to detach cloud project (%s): %s", vrackId, projectId, err)
	}

	return nil
}

func TestAccVrackCloudProject_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccCheckVrackCloudProjectPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVrackCloudProjectConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_vrack_cloudproject.vcp", "service_name", os.Getenv("OVH_VRACK_SERVICE_TEST")),
					resource.TestCheckResourceAttr("ovh_vrack_cloudproject.vcp", "project_id", os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")),
				),
			},
		},
	})
}

func TestAccVrackCloudProjectDeprecated_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccCheckVrackCloudProjectPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVrackCloudProjectDeprecatedConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_vrack_cloudproject.attach", "service_name", os.Getenv("OVH_VRACK_SERVICE_TEST")),
					resource.TestCheckResourceAttr("ovh_vrack_cloudproject.attach", "project_id", os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")),
				),
			},
		},
	})
}

func testAccCheckVrackCloudProjectPreCheck(t *testing.T) {
	testAccPreCheckVRack(t)
	testAccCheckVRackExists(t)
	testAccCheckCloudExists(t)
}
