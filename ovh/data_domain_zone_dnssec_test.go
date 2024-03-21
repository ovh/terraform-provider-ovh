package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const testAccCheckOvhDomainZoneDnssecConfig_basic = `
data "ovh_domain_zone_dnssec" "sec" {
	zone_name = "%s"
}`

func TestAccDomainZoneDnssecDataSource_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckDomain(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckOvhDomainZoneDnssecConfig_basic, os.Getenv("OVH_ZONE_TEST")),
				Check:  resource.TestCheckResourceAttr("data.ovh_domain_zone_dnssec.sec", "status", "disabled"),
			},
		},
	})
}
