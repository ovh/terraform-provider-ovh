package ovh

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"regexp"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func init() {
	resource.AddTestSweepers("ovh_domain_zone_dynhost_login", &resource.Sweeper{
		Name: "ovh_domain_zone_dynhost_login",
		F:    testSweepDomainZoneDynhostLogin,
	})
}

func testSweepDomainZoneDynhostLogin(region string) error {
	client, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting client: %s", err)
	}

	zoneName := os.Getenv("OVH_ZONE_TEST")
	if zoneName == "" {
		log.Print("[DEBUG] OVH_ZONE_TEST is not set. No zone dynhost logins to sweep")
		return nil
	}

	logins := make([]string, 0)
	if err := client.Get(fmt.Sprintf("/domain/zone/%s/dynhost/login", url.PathEscape(zoneName)), &logins); err != nil {
		return fmt.Errorf("Error calling /domain/zone/%s/dynhost/login:\n\t %q", zoneName, err)
	}

	if len(logins) == 0 {
		log.Print("[DEBUG] No record to sweep")
		return nil
	}

	for _, l := range logins {
		login := &DomainZoneDynhostLoginModel{}

		if err := client.Get(fmt.Sprintf("/domain/zone/%s/dynhost/login/%s", url.PathEscape(zoneName), url.PathEscape(l)), &login); err != nil {
			return fmt.Errorf("Error calling /domain/zone/%s/dynhost/login/%s:\n\t %q", zoneName, l, err)
		}

		log.Printf("[DEBUG] login found %v", login)

		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			log.Printf("[INFO] Deleting login %v", login)
			if err := client.Delete(fmt.Sprintf("/domain/zone/%s/dynhost/login/%s", url.PathEscape(zoneName), url.PathEscape(l)), nil); err != nil {
				return resource.RetryableError(err)
			}
			// Successful delete
			return nil
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func TestAccDomainZoneDynHostLogin_Basic(t *testing.T) {
	zone := os.Getenv("OVH_ZONE_TEST")
	subdomain := acctest.RandomWithPrefix(test_prefix)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckDomain(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// provider should send an error if no zone is given
			{
				Config:      testAccOvhDomainZoneDynhostLoginNoZoneNameConfig(subdomain, "suffix", "aStrongPassword"),
				ExpectError: regexp.MustCompile(`The argument "zone_name" is required, but no definition was found.`),
			},
			// provider should send an error if zone does not exist
			{
				Config:      testAccOvhDomainZoneDynhostLoginConfig("non-existing-domain.com", subdomain, "suffix", "aStrongPassword"),
				ExpectError: regexp.MustCompile(`Client::NotFound`),
			},
			// provider should send an error if password is not strong password (i.e. < 12char)
			{
				Config:      testAccOvhDomainZoneDynhostLoginConfig(zone, subdomain, "suffix", "weakPwd"),
				ExpectError: regexp.MustCompile(`'Length must be between 12 and 512.'`),
			},

			// resource creation
			{
				Config: testAccOvhDomainZoneDynhostLoginConfig(zone, subdomain, "suffix", "aStrongPassword"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_domain_zone_dynhost_login.this", "login", zone+"-suffix"),
					resource.TestCheckResourceAttr("ovh_domain_zone_dynhost_login.this", "login_suffix", "suffix"),
					resource.TestCheckResourceAttr("ovh_domain_zone_dynhost_login.this", "password", "aStrongPassword"),
					resource.TestCheckResourceAttr("ovh_domain_zone_dynhost_login.this", "sub_domain", subdomain),
					resource.TestCheckResourceAttr("ovh_domain_zone_dynhost_login.this", "zone", zone),
				),
			},
			// resource password update
			{
				Config: testAccOvhDomainZoneDynhostLoginConfig(zone, subdomain, "suffix", "anotherStrongPassword"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_domain_zone_dynhost_login.this", "login", zone+"-suffix"),
					resource.TestCheckResourceAttr("ovh_domain_zone_dynhost_login.this", "login_suffix", "suffix"),
					resource.TestCheckResourceAttr("ovh_domain_zone_dynhost_login.this", "password", "anotherStrongPassword"),
					resource.TestCheckResourceAttr("ovh_domain_zone_dynhost_login.this", "sub_domain", subdomain),
					resource.TestCheckResourceAttr("ovh_domain_zone_dynhost_login.this", "zone", zone),
				),
			},
			// resource suffix update
			{
				Config: testAccOvhDomainZoneDynhostLoginConfig(zone, subdomain, "another_suffix", "anotherStrongPassword"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_domain_zone_dynhost_login.this", "login", zone+"-suffix"),
					resource.TestCheckResourceAttr("ovh_domain_zone_dynhost_login.this", "login_suffix", "another_suffix"),
					resource.TestCheckResourceAttr("ovh_domain_zone_dynhost_login.this", "password", "anotherStrongPassword"),
					resource.TestCheckResourceAttr("ovh_domain_zone_dynhost_login.this", "sub_domain", subdomain),
					resource.TestCheckResourceAttr("ovh_domain_zone_dynhost_login.this", "zone", zone),
				),
			},
		},
	})
}

func testAccOvhDomainZoneDynhostLoginNoZoneNameConfig(subDomain, loginSuffix, password string) string {
	return fmt.Sprintf(`	
resource "ovh_domain_zone_dynhost_login" "this" {
  sub_domain   = %q
  login_suffix = %q
  password     = %q
}`, subDomain, loginSuffix, password)
}

func testAccOvhDomainZoneDynhostLoginConfig(zoneName, subDomain, loginSuffix, password string) string {
	return fmt.Sprintf(`	
resource "ovh_domain_zone_dynhost_login" "this" {
  zone_name    = %q
  sub_domain   = %q
  login_suffix = %q
  password     = %q
}`, zoneName, subDomain, loginSuffix, password)
}
