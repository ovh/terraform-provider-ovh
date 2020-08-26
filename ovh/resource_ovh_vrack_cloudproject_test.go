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
  vrack_id = "%s"
  project_id = "%s"
}
`, os.Getenv("OVH_VRACK"), os.Getenv("OVH_PUBLIC_CLOUD"))

var testAccVrackCloudProjectConfig_legacy = fmt.Sprintf(`
resource "ovh_vrack_cloudproject" "attach" {
  vrack_id = "%s"
  project_id = "%s"
}
`, os.Getenv("OVH_VRACK"), os.Getenv("OVH_PUBLIC_CLOUD"))

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

	vrackId := os.Getenv("OVH_VRACK")
	if vrackId == "" {
		log.Print("[DEBUG] OVH_VRACK is not set. No vrack_cloud_project to sweep")
		return nil
	}

	projectId := os.Getenv("OVH_PUBLIC_CLOUD")
	if projectId == "" {
		log.Print("[DEBUG] OVH_PUBLIC_CLOUD is not set. No vrack_cloud_project to sweep")
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
					resource.TestCheckResourceAttr("ovh_vrack_cloudproject.vcp", "vrack_id", os.Getenv("OVH_VRACK")),
					resource.TestCheckResourceAttr("ovh_vrack_cloudproject.vcp", "project_id", os.Getenv("OVH_PUBLIC_CLOUD")),
				),
			},
		},
	})
}

func TestAccVrackCloudProject_legacy(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccCheckVrackCloudProjectPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVrackCloudProjectConfig_legacy,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_vrack_cloudproject.attach", "vrack_id", os.Getenv("OVH_VRACK")),
					resource.TestCheckResourceAttr("ovh_vrack_cloudproject.attach", "project_id", os.Getenv("OVH_PUBLIC_CLOUD")),
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
