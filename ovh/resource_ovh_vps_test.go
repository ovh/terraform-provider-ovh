package ovh

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccOVHVps_importBasic(t *testing.T) {
	resourceName := "ovh_vps.bar"
	vps_id := os.Getenv("OVH_VPS_ID")
	vps_type := os.Getenv("OVH_VPS_TYPE")
	vps_displayname := os.Getenv("OVH_VPS_DISPLAYNAME")
	checkFn := func(s []*terraform.InstanceState) error {
		if len(s) == 1 {
			id := s[0].ID
			attr := s[0].Attributes
			if id == vps_id && attr["displayname"] == vps_displayname && attr["type"] == vps_type {
				return nil
			}
		}
		return fmt.Errorf("Bad content: %#v", s[0])
	}
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOVHDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config:      testAccOVHVpsConfig,
				ExpectError: regexp.MustCompile(`use import`),
			},

			resource.TestStep{
				Config:           fmt.Sprintf(testAccOVHVpsConfig, vps_type, vps_displayname),
				ResourceName:     resourceName,
				ImportStateId:    vps_id,
				ImportState:      true,
				ImportStateCheck: checkFn,
			},
		},
	})
}

func testAccCheckOVHDestroy(s *terraform.State) error {
	return nil
}

const testAccOVHVpsConfig = `
resource ovh_vps bar {
  type = "%s"
  displayname= "%s"
}
`
