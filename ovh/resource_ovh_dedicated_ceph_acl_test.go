package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
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
)

func TestAccDedicatedCephACLCreate(t *testing.T) {
	config := fmt.Sprintf(
		testAccDedicatedCephACLConfig,
		os.Getenv("OVH_DEDICATED_CEPH"),
		"109.190.254.59",
		"255.255.255.255",
	)
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckDedicatedCeph(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_dedicated_ceph_acl.acl", "network", "109.190.254.59"),
					resource.TestCheckResourceAttr("ovh_dedicated_ceph_acl.acl", "netmask", "255.255.255.255"),
				),
			},
		},
	})
}
