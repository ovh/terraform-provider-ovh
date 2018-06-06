package ovh

import (
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccOvhDomainZoneRedirection_Basic(t *testing.T) {
	var redirection OvhDomainZoneRedirection
	zone := os.Getenv("OVH_ZONE")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOvhDomainZoneRedirectionDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: fmt.Sprintf(testAccCheckOvhDomainZoneRedirectionConfig_basic, zone),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOvhDomainZoneRedirectionExists("ovh_domain_zone_redirection.foobar", &redirection),
					resource.TestCheckResourceAttr(
						"ovh_domain_zone_redirection.foobar", "subdomain", "terraform"),
					resource.TestCheckResourceAttr(
						"ovh_domain_zone_redirection.foobar", "zone", zone),
					resource.TestCheckResourceAttr(
						"ovh_domain_zone_redirection.foobar", "target", "https://terraform.net"),
					resource.TestCheckResourceAttr(
						"ovh_domain_zone_redirection.foobar", "type", "visible"),
				),
			},
		},
	})
}

func TestAccOvhDomainZoneRedirection_Updated(t *testing.T) {
	redirection := OvhDomainZoneRedirection{}
	zone := os.Getenv("OVH_ZONE")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOvhDomainZoneRedirectionDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: fmt.Sprintf(testAccCheckOvhDomainZoneRedirectionConfig_basic, zone),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOvhDomainZoneRedirectionExists("ovh_domain_zone_redirection.foobar", &redirection),
					resource.TestCheckResourceAttr(
						"ovh_domain_zone_redirection.foobar", "subdomain", "terraform"),
					resource.TestCheckResourceAttr(
						"ovh_domain_zone_redirection.foobar", "zone", zone),
					resource.TestCheckResourceAttr(
						"ovh_domain_zone_redirection.foobar", "target", "https://terraform.net"),
				),
			},
			resource.TestStep{
				Config: fmt.Sprintf(testAccCheckOvhDomainZoneRedirectionConfig_new_value_1, zone),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOvhDomainZoneRedirectionExists("ovh_domain_zone_redirection.foobar", &redirection),
					resource.TestCheckResourceAttr(
						"ovh_domain_zone_redirection.foobar", "subdomain", "terraform"),
					resource.TestCheckResourceAttr(
						"ovh_domain_zone_redirection.foobar", "zone", zone),
					resource.TestCheckResourceAttr(
						"ovh_domain_zone_redirection.foobar", "target", "https://terraform.io"),
				),
			},
			resource.TestStep{
				Config: fmt.Sprintf(testAccCheckOvhDomainZoneRedirectionConfig_new_value_2, zone),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOvhDomainZoneRedirectionExists("ovh_domain_zone_redirection.foobar", &redirection),
					resource.TestCheckResourceAttr(
						"ovh_domain_zone_redirection.foobar", "subdomain", "terraform2"),
					resource.TestCheckResourceAttr(
						"ovh_domain_zone_redirection.foobar", "zone", zone),
					resource.TestCheckResourceAttr(
						"ovh_domain_zone_redirection.foobar", "target", "https://terraform.io"),
				),
			},
			resource.TestStep{
				Config: fmt.Sprintf(testAccCheckOvhDomainZoneRedirectionConfig_new_value_3, zone),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOvhDomainZoneRedirectionExists("ovh_domain_zone_redirection.foobar", &redirection),
					resource.TestCheckResourceAttr(
						"ovh_domain_zone_redirection.foobar", "subdomain", "terraform3"),
					resource.TestCheckResourceAttr(
						"ovh_domain_zone_redirection.foobar", "zone", zone),
					resource.TestCheckResourceAttr(
						"ovh_domain_zone_redirection.foobar", "target", "https://terraform.com"),
				),
			},
		},
	})
}

func testAccCheckOvhDomainZoneRedirectionDestroy(s *terraform.State) error {
	provider := testAccProvider.Meta().(*Config)
	zone := os.Getenv("OVH_ZONE")

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ovh_domain_zone_redirection" {
			continue
		}

		resultRedirection := OvhDomainZoneRedirection{}
		err := provider.OVHClient.Get(
			fmt.Sprintf("/domain/zone/%s/redirection/%s", zone, rs.Primary.ID),
			&resultRedirection,
		)

		if err == nil {
			return fmt.Errorf("Redirection still exists")
		}
	}

	return nil
}

func testAccCheckOvhDomainZoneRedirectionExists(n string, redirection *OvhDomainZoneRedirection) resource.TestCheckFunc {
	zone := os.Getenv("OVH_ZONE")
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Redirection ID is set")
		}

		provider := testAccProvider.Meta().(*Config)

		err := provider.OVHClient.Get(
			fmt.Sprintf("/domain/zone/%s/redirection/%s", zone, rs.Primary.ID),
			redirection,
		)

		if err != nil {
			return err
		}

		if strconv.Itoa(redirection.Id) != rs.Primary.ID {
			return fmt.Errorf("Redirection not found")
		}

		return nil
	}
}

const testAccCheckOvhDomainZoneRedirectionConfig_basic = `
resource "ovh_domain_zone_redirection" "foobar" {
	zone = "%s"
	subdomain = "terraform"
	target = "https://terraform.net"
	type = "visible"
}`

const testAccCheckOvhDomainZoneRedirectionConfig_new_value_1 = `
resource "ovh_domain_zone_redirection" "foobar" {
	zone = "%s"
	subdomain = "terraform"
	target = "https://terraform.io"
	type = "visible"
}
`

const testAccCheckOvhDomainZoneRedirectionConfig_new_value_2 = `
resource "ovh_domain_zone_redirection" "foobar" {
	zone = "%s"
	subdomain = "terraform2"
	target = "https://terraform.io"
	type = "visible"
}
`

const testAccCheckOvhDomainZoneRedirectionConfig_new_value_3 = `
resource "ovh_domain_zone_redirection" "foobar" {
	zone = "%s"
	subdomain = "terraform3"
	target = "https://terraform.com"
	type = "visible"
}`
