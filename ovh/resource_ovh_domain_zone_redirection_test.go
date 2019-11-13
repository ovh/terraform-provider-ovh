package ovh

import (
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"log"
	"strings"
	"time"
)

func init() {
	resource.AddTestSweepers("ovh_domain_zone_redirection", &resource.Sweeper{
		Name: "ovh_domain_zone_redirection",
		F:    testSweepDomainZoneRedirection,
	})
}

func testSweepDomainZoneRedirection(region string) error {
	client, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting client: %s", err)
	}

	zoneName := os.Getenv("OVH_ZONE")
	if zoneName == "" {
		return fmt.Errorf("OVH_ZONE must be set")
	}

	dz := &DomainZone{}

	if err := client.Get(fmt.Sprintf("/domain/zone/%s", zoneName), &dz); err != nil {
		return fmt.Errorf("Error calling /domain/zone/%s:\n\t %q", zoneName, err)
	}

	redirections := make([]int64, 0)
	if err := client.Get(fmt.Sprintf("/domain/zone/%s/redirection", zoneName), &redirections); err != nil {
		return fmt.Errorf("Error calling /domain/zone/%s:\n\t %q", zoneName, err)
	}

	if len(redirections) == 0 {
		log.Print("[DEBUG] No redirection to sweep")
		return nil
	}

	for _, rec := range redirections {
		redirection := &OvhDomainZoneRedirection{}

		if err := client.Get(fmt.Sprintf("/domain/zone/%s/redirection/%v", zoneName, rec), &redirection); err != nil {
			return fmt.Errorf("Error calling /domain/zone/%s/redirection/%v:\n\t %q", zoneName, rec, err)
		}

		if !strings.HasPrefix(redirection.SubDomain, test_prefix) {
			continue
		}

		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			if err := client.Delete(fmt.Sprintf("/domain/zone/%s/redirection/%v", zoneName, rec), nil); err != nil {
				return resource.RetryableError(err)
			}
			// Successful delete
			return nil
		})
		if err != nil {
			return err
		}
	}

	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		err := client.Post(
			fmt.Sprintf("/domain/zone/%s/refresh", zoneName),
			nil,
			nil,
		)

		if err != nil {
			return resource.RetryableError(fmt.Errorf("Error refresh OVH Zone: %s", err))
		}
		// Successful refresh
		return nil
	})

	return nil
}

func TestAccOvhDomainZoneRedirection_Basic(t *testing.T) {
	var redirection OvhDomainZoneRedirection
	zone := os.Getenv("OVH_ZONE")
	subdomain := acctest.RandomWithPrefix(test_prefix)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckDomain(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOvhDomainZoneRedirectionDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckOvhDomainZoneRedirectionConfig_basic, zone, subdomain),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOvhDomainZoneRedirectionExists("ovh_domain_zone_redirection.foobar", &redirection),
					resource.TestCheckResourceAttr(
						"ovh_domain_zone_redirection.foobar", "subdomain", subdomain),
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
	subdomain := acctest.RandomWithPrefix(test_prefix)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckDomain(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOvhDomainZoneRedirectionDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckOvhDomainZoneRedirectionConfig_basic, zone, subdomain),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOvhDomainZoneRedirectionExists("ovh_domain_zone_redirection.foobar", &redirection),
					resource.TestCheckResourceAttr(
						"ovh_domain_zone_redirection.foobar", "subdomain", subdomain),
					resource.TestCheckResourceAttr(
						"ovh_domain_zone_redirection.foobar", "zone", zone),
					resource.TestCheckResourceAttr(
						"ovh_domain_zone_redirection.foobar", "target", "https://terraform.net"),
				),
			},
			{
				Config: fmt.Sprintf(testAccCheckOvhDomainZoneRedirectionConfig_new_value_1, zone, subdomain),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOvhDomainZoneRedirectionExists("ovh_domain_zone_redirection.foobar", &redirection),
					resource.TestCheckResourceAttr(
						"ovh_domain_zone_redirection.foobar", "subdomain", subdomain),
					resource.TestCheckResourceAttr(
						"ovh_domain_zone_redirection.foobar", "zone", zone),
					resource.TestCheckResourceAttr(
						"ovh_domain_zone_redirection.foobar", "target", "https://terraform.io"),
				),
			},
			{
				Config: fmt.Sprintf(testAccCheckOvhDomainZoneRedirectionConfig_new_value_2, zone, subdomain),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOvhDomainZoneRedirectionExists("ovh_domain_zone_redirection.foobar", &redirection),
					resource.TestCheckResourceAttr(
						"ovh_domain_zone_redirection.foobar", "subdomain", fmt.Sprintf("%s2", subdomain)),
					resource.TestCheckResourceAttr(
						"ovh_domain_zone_redirection.foobar", "zone", zone),
					resource.TestCheckResourceAttr(
						"ovh_domain_zone_redirection.foobar", "target", "https://terraform.io"),
				),
			},
			{
				Config: fmt.Sprintf(testAccCheckOvhDomainZoneRedirectionConfig_new_value_3, zone, subdomain),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOvhDomainZoneRedirectionExists("ovh_domain_zone_redirection.foobar", &redirection),
					resource.TestCheckResourceAttr(
						"ovh_domain_zone_redirection.foobar", "subdomain", fmt.Sprintf("%s3", subdomain)),
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
	subdomain = "%s"
	target = "https://terraform.net"
	type = "visible"
}`

const testAccCheckOvhDomainZoneRedirectionConfig_new_value_1 = `
resource "ovh_domain_zone_redirection" "foobar" {
	zone = "%s"
	subdomain = "%s"
	target = "https://terraform.io"
	type = "visible"
}
`

const testAccCheckOvhDomainZoneRedirectionConfig_new_value_2 = `
resource "ovh_domain_zone_redirection" "foobar" {
	zone = "%s"
	subdomain = "%s2"
	target = "https://terraform.io"
	type = "visible"
}
`

const testAccCheckOvhDomainZoneRedirectionConfig_new_value_3 = `
resource "ovh_domain_zone_redirection" "foobar" {
	zone = "%s"
	subdomain = "%s3"
	target = "https://terraform.com"
	type = "visible"
}`
