package ovh

import (
	"fmt"
	"log"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

var testAccPublicCloudPrivateNetworkConfig = fmt.Sprintf(`
resource "ovh_vrack_cloudproject" "attach" {
  vrack_id   = "%s"
  project_id = "%s"
}

data "ovh_cloud_regions" "regions" {
  project_id = ovh_vrack_cloudproject.attach.project_id
}

resource "ovh_cloud_network_private" "network" {
  project_id = ovh_vrack_cloudproject.attach.project_id
  vlan_id    = 0
  name       = "terraform_testacc_private_net"
  regions    = tolist(data.ovh_cloud_regions.regions.names)
}

`, os.Getenv("OVH_VRACK"), os.Getenv("OVH_PUBLIC_CLOUD"))

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
		return fmt.Errorf("OVH_VRACK must be set")
	}

	projectId := os.Getenv("OVH_PUBLIC_CLOUD")
	if projectId == "" {
		return fmt.Errorf("OVH_PUBLIC_CLOUD must be set")
	}

	networkIds := []string{}
	err = client.Get(fmt.Sprintf("/cloud/project/%s/network/private", projectId), &networkIds)
	if err != nil {
		return fmt.Errorf("error listing private networks for project %q:\n\t %q", projectId, err)
	}

	for _, n := range networkIds {
		r := &PublicCloudPrivateNetworkResponse{}
		err = client.Get(fmt.Sprintf("/cloud/project/%s/network/private/%s", projectId, n), r)
		if err != nil {
			return fmt.Errorf("error getting private network %q for project %q:\n\t %q", n, projectId, err)
		}

		if !strings.HasPrefix(r.Name, test_prefix) {
			continue
		}

		log.Printf("[DEBUG] found dangling network & subnets for project: %s, id: %s", projectId, n)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			subnetIds := []string{}
			err = client.Get(fmt.Sprintf("/cloud/project/%s/network/private/%s/subnet", projectId, n), &subnetIds)
			if err != nil {
				return resource.RetryableError(fmt.Errorf("error listing private network subnets for project %q:\n\t %q", projectId, err))
			}

			for _, s := range subnetIds {
				if err := client.Delete(fmt.Sprintf("/cloud/project/%s/network/private/%s/subnet/%s", projectId, n, s), nil); err != nil {
					return resource.RetryableError(err)
				}
			}

			if err := client.Delete(fmt.Sprintf("/cloud/project/%s/network/private/%s", projectId, n), nil); err != nil {
				return resource.RetryableError(err)
			}

			// Successful cascade delete
			log.Printf("[DEBUG] successful cascade delete of network & subnets for project: %s, id: %s", projectId, n)
			return nil
		})

		if err != nil {
			return err
		}

	}

	return nil
}

func TestAccPublicCloudPrivateNetwork_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccCheckPublicCloudPrivateNetworkPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPublicCloudPrivateNetworkDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPublicCloudPrivateNetworkConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVRackPublicCloudAttachmentExists("ovh_vrack_cloudproject.attach", t),
					testAccCheckPublicCloudPrivateNetworkExists("ovh_cloud_network_private.network", t),
				),
			},
		},
	})
}

func testAccCheckPublicCloudPrivateNetworkPreCheck(t *testing.T) {
	testAccPreCheckPublicCloud(t)
	testAccCheckPublicCloudExists(t)
}

func testAccCheckPublicCloudPrivateNetworkExists(n string, t *testing.T) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := testAccProvider.Meta().(*Config)

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		if rs.Primary.Attributes["project_id"] == "" {
			return fmt.Errorf("No Project ID is set")
		}

		return publicCloudPrivateNetworkExists(rs.Primary.Attributes["project_id"], rs.Primary.ID, config.OVHClient)
	}
}

func testAccCheckPublicCloudPrivateNetworkDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ovh_cloud_network_private" {
			continue
		}

		err := publicCloudPrivateNetworkExists(rs.Primary.Attributes["project_id"], rs.Primary.ID, config.OVHClient)
		if err == nil {
			return fmt.Errorf("VRack > Public Cloud Private Network still exists")
		}

	}
	return nil
}
