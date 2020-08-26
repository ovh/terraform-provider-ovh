package ovh

import (
	"fmt"
	"log"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

var testAccCloudNetworkPrivateConfig_attachVrack = `
resource "ovh_vrack_cloudproject" "attach" {
  vrack_id   = "%s"
  project_id = "%s"
}

data "ovh_cloud_regions" "regions" {
  project_id = ovh_vrack_cloudproject.attach.project_id

  has_services_up = ["network"]
}
`

var testAccCloudNetworkPrivateConfig_noAttachVrack = `
data "ovh_cloud_regions" "regions" {
  project_id = "%s"

  has_services_up = ["network"]
}
`

var testAccCloudNetworkPrivateConfig_basic = `
%s

resource "ovh_cloud_network_private" "network" {
  project_id = data.ovh_cloud_regions.regions.project_id
  vlan_id    = 0
  name       = "terraform_testacc_private_net"
  regions    = tolist(data.ovh_cloud_regions.regions.names)
}
`

func testAccCloudNetworkPrivateConfig() string {
	attachVrack := fmt.Sprintf(
		testAccCloudNetworkPrivateConfig_attachVrack,
		os.Getenv("OVH_VRACK"),
		os.Getenv("OVH_PUBLIC_CLOUD"),
	)
	noAttachVrack := fmt.Sprintf(
		testAccCloudNetworkPrivateConfig_noAttachVrack,
		os.Getenv("OVH_PUBLIC_CLOUD"),
	)

	if os.Getenv("OVH_ATTACH_VRACK") == "0" {
		return fmt.Sprintf(
			testAccCloudNetworkPrivateConfig_basic,
			noAttachVrack,
		)
	}

	return fmt.Sprintf(
		testAccCloudNetworkPrivateConfig_basic,
		attachVrack,
	)
}

func init() {
	resource.AddTestSweepers("ovh_cloud_network_private", &resource.Sweeper{
		Name: "ovh_cloud_network_private",
		F:    testSweepCloudNetworkPrivate,
	})
}

func testSweepCloudNetworkPrivate(region string) error {
	client, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting client: %s", err)
	}

	vrack := os.Getenv("OVH_VRACK")
	if vrack == "" {
		log.Print("[DEBUG] OVH_VRACK is not set. No cloud_network_private to sweep")
		return nil
	}

	projectId := os.Getenv("OVH_PUBLIC_CLOUD")
	if projectId == "" {
		log.Print("[DEBUG] OVH_PUBLIC_CLOUD is not set. No cloud_network_private to sweep")
		return nil
	}

	networks := []CloudNetworkPrivateResponse{}
	err = client.Get(fmt.Sprintf("/cloud/project/%s/network/private", projectId), &networks)
	if err != nil {
		return fmt.Errorf("error listing private networks for project %q:\n\t %q", projectId, err)
	}

	for _, n := range networks {
		if !strings.HasPrefix(n.Name, test_prefix) {
			continue
		}

		log.Printf("[DEBUG] found dangling network & subnets for project: %s, id: %s", projectId, n.Id)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			subnetIds := []string{}
			err = client.Get(fmt.Sprintf("/cloud/project/%s/network/private/%s/subnet", projectId, n.Id), &subnetIds)
			if err != nil {
				return resource.RetryableError(fmt.Errorf("error listing private network subnets for project %q:\n\t %q", projectId, err))
			}

			for _, s := range subnetIds {
				if err := client.Delete(fmt.Sprintf("/cloud/project/%s/network/private/%s/subnet/%s", projectId, n.Id, s), nil); err != nil {
					return resource.RetryableError(err)
				}
			}

			if err := client.Delete(fmt.Sprintf("/cloud/project/%s/network/private/%s", projectId, n.Id), nil); err != nil {
				return resource.RetryableError(err)
			}

			// Successful cascade delete
			log.Printf("[DEBUG] successful cascade delete of network & subnets for project: %s, id: %s", projectId, n.Id)
			return nil
		})

		if err != nil {
			return err
		}

	}

	return nil
}

func TestAccCloudNetworkPrivate_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccCheckCloudNetworkPrivatePreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudNetworkPrivateConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("ovh_cloud_network_private.network", "project_id"),
					resource.TestCheckResourceAttrSet("ovh_cloud_network_private.network", "id"),
					resource.TestCheckResourceAttr("ovh_cloud_network_private.network", "vlan_id", "0"),
				),
			},
		},
	})
}

func testAccCheckCloudNetworkPrivatePreCheck(t *testing.T) {
	testAccPreCheckCloud(t)
	testAccCheckCloudExists(t)
	testAccPreCheckVRack(t)
}
