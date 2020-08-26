package ovh

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const (
	testAccDedicatedCephACLConfig = `
data "ovh_dedicated_ceph" "ceph"{
  service_name = "%s"
}

resource "ovh_dedicated_ceph_acl" "acl" {
  service_name = data.ovh_dedicated_ceph.ceph.id
  network = "%s"
  netmask = "%s"
}
`
	testAccDedicatedCephACLNetmask = "255.255.255.255"
	testAccDedicatedCephACLNetwork = "109.190.254.59"
)

func init() {
	resource.AddTestSweepers("ovh_dedicated_ceph_acl", &resource.Sweeper{
		Name: "ovh_dedicated_ceph_acl",
		F:    testSweepDedicatedCephACL,
	})
}

func testSweepDedicatedCephACL(region string) error {
	client, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting client: %s", err)
	}
	serviceName := os.Getenv("OVH_DEDICATED_CEPH")
	if serviceName == "" {
		log.Print("[DEBUG] No OVH_DEDICATED_CEPH envvar specified. nothing to sweep")
		return nil
	}

	var acls []DedicatedCephACL
	endpoint := fmt.Sprintf("/dedicated/ceph/%s/acl", url.PathEscape(serviceName))

	if err := client.Get(endpoint, acls); err != nil {
		return fmt.Errorf("Error calling GET %s:\n\t%q", endpoint, err)
	}

	if len(acls) == 0 {
		log.Printf("[DEBUG] No ACL to sweep on dedicated ceph %s", serviceName)
		return nil
	}

	for _, acl := range acls {
		if acl.Netmask != testAccDedicatedCephACLNetmask || acl.Network != testAccDedicatedCephACLNetwork {
			continue
		}

		log.Printf(
			"[INFO] Deleting ACL %d %v/%v on service %v",
			acl.Id,
			acl.Network,
			acl.Netmask,
			serviceName,
		)
		endpoint = fmt.Sprintf("/dedicated/ceph/%s/acl/%d", url.PathEscape(serviceName), acl.Id)
		var taskId string
		if err := client.Delete(endpoint, &taskId); err != nil {
			return fmt.Errorf("Error calling DELETE %s:\n\t%q", endpoint, err)
		}

		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			endpoint = fmt.Sprintf("/dedicated/ceph/%s/task/%s", serviceName, taskId)
			var stateResp []DedicatedCephTask
			if err := client.Get(endpoint, &stateResp); err != nil {
				return resource.RetryableError(err)
			}
			state := stateResp[0].State
			if state != "DONE" {
				return resource.RetryableError(fmt.Errorf("Deletion Task still pending: %v", state))
			}

			// Successful delete
			return nil
		})

		if err != nil {
			return fmt.Errorf("Error waiting for CEPH ACL deletion:\n\t %q", err)
		}
	}
	return nil
}

func TestAccDedicatedCephACLCreate(t *testing.T) {
	config := fmt.Sprintf(
		testAccDedicatedCephACLConfig,
		os.Getenv("OVH_DEDICATED_CEPH"),
		testAccDedicatedCephACLNetwork,
		testAccDedicatedCephACLNetmask,
	)
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckDedicatedCeph(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_dedicated_ceph_acl.acl", "network", testAccDedicatedCephACLNetwork),
					resource.TestCheckResourceAttr("ovh_dedicated_ceph_acl.acl", "netmask", testAccDedicatedCephACLNetmask),
				),
			},
		},
	})
}
