package ovh

import (
	"fmt"
	"net/url"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

const testAccDomainZoneDnssecResourceConfig = `
resource "ovh_domain_zone_dnssec" "sec" {
	zone_name = "%s"
}`

func testAccDomainZoneDnssecStatusMatchesAPI(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource not found: %s", resourceName)
		}

		zone := rs.Primary.Attributes["zone_name"]
		var r struct {
			Status string `json:"status"`
		}
		endpoint := fmt.Sprintf("/domain/zone/%s/dnssec", url.PathEscape(zone))
		if err := testAccOVHClient.Get(endpoint, &r); err != nil {
			return fmt.Errorf("calling Get %s: %w", endpoint, err)
		}

		if got := rs.Primary.Attributes["status"]; got != r.Status {
			return fmt.Errorf("state status %q does not match API status %q", got, r.Status)
		}
		return nil
	}
}

func TestAccDomainZoneDnssecResource_Basic(t *testing.T) {
	config := fmt.Sprintf(testAccDomainZoneDnssecResourceConfig, os.Getenv("OVH_ZONE_TEST"))

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckDomain(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("ovh_domain_zone_dnssec.sec", "status"),
					testAccDomainZoneDnssecStatusMatchesAPI("ovh_domain_zone_dnssec.sec"),
				),
			},
			{
				// second apply forces a refresh (runs Read); status must still match the API
				Config: config,
				Check:  testAccDomainZoneDnssecStatusMatchesAPI("ovh_domain_zone_dnssec.sec"),
			},
		},
	})
}
