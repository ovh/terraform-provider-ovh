package ovh

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func init() {
	resource.AddTestSweepers("ovh_email_domain_dkim", &resource.Sweeper{
		Name: "ovh_email_domain_dkim",
		F:    testSweepEmailDomainDkim,
	})
}

func testSweepEmailDomainDkim(region string) error {
	client, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting client: %s", err)
	}

	domain := os.Getenv("OVH_EMAIL_DOMAIN_TEST")
	if domain == "" {
		log.Print("[DEBUG] OVH_EMAIL_DOMAIN_TEST is not set. No email domain DKIM to sweep")
		return nil
	}

	var dkim struct {
		Status string `json:"status"`
	}
	if err := client.Get(fmt.Sprintf("/email/domain/%s/dkim", url.PathEscape(domain)), &dkim); err != nil {
		return fmt.Errorf("Error calling /email/domain/%s/dkim:\n\t %q", domain, err)
	}

	if dkim.Status == dkimStatusDisabled {
		log.Print("[DEBUG] DKIM is already disabled, nothing to sweep")
		return nil
	}

	log.Printf("[DEBUG] Disabling DKIM on email domain %s", domain)
	if err := client.Put(fmt.Sprintf("/email/domain/%s/dkim/disable", url.PathEscape(domain)), nil, nil); err != nil {
		return fmt.Errorf("Error calling /email/domain/%s/dkim/disable:\n\t %q", domain, err)
	}

	return nil
}

// testAccPreCheckEmailDomainDkim skips rather than failing obscurely when the
// test domain is not in a fit state.
//
// Creating on a domain that already has DKIM on gets a 409, and quietly
// disabling somebody's DKIM to make room is too destructive for a precheck.
// A domain whose offer does not support DKIM, a redirect-only one for instance,
// accepts the enable call with a 200, allocates no selector and stays
// "disabled" forever, so the test would only fail two minutes later on the
// activation timeout.
func testAccPreCheckEmailDomainDkim(t *testing.T) {
	testAccPreCheckEmailDomain(t)

	domain := os.Getenv("OVH_EMAIL_DOMAIN_TEST")

	var dkim struct {
		Status string `json:"status"`
	}
	endpoint := fmt.Sprintf("/email/domain/%s/dkim", url.PathEscape(domain))
	if err := testAccOVHClient.Get(endpoint, &dkim); err != nil {
		t.Fatalf("error calling %s: %v", endpoint, err)
	}

	if dkim.Status != dkimStatusDisabled {
		t.Skipf("DKIM is %q on %s, expected %q: disable it first, or run the ovh_email_domain_dkim sweeper",
			dkim.Status, domain, dkimStatusDisabled)
	}

	var service struct {
		Offer string `json:"offer"`
	}
	endpoint = fmt.Sprintf("/email/domain/%s", url.PathEscape(domain))
	if err := testAccOVHClient.Get(endpoint, &service); err != nil {
		t.Fatalf("error calling %s: %v", endpoint, err)
	}

	if !strings.Contains(strings.ToUpper(service.Offer), "MXPLAN") {
		t.Skipf("%s is on the %q offer, which does not support DKIM: set OVH_EMAIL_DOMAIN_TEST to an MX Plan domain",
			domain, service.Offer)
	}
}

func TestAccEmailDomainDkim_Basic(t *testing.T) {
	domain := os.Getenv("OVH_EMAIL_DOMAIN_TEST")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckEmailDomainDkim(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Error when domain is missing
			{
				Config:      testAccEmailDomainDkimConfig_noDomain(),
				ExpectError: regexp.MustCompile(`The argument "domain" is required`),
			},
			// Successful activation
			{
				Config: testAccEmailDomainDkimConfig(domain),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_email_domain_dkim.test", "domain", domain),
					resource.TestCheckResourceAttr("ovh_email_domain_dkim.test", "id", domain),
					// The settled status depends on whether the selector records are
					// already published, so only assert that DKIM is no longer off.
					resource.TestCheckResourceAttrWith("ovh_email_domain_dkim.test", "status", func(value string) error {
						if value == dkimStatusDisabled || value == dkimStatusModifying {
							return fmt.Errorf("expected a settled, enabled status, got %q", value)
						}
						return nil
					}),
					resource.TestCheckResourceAttrSet("ovh_email_domain_dkim.test", "autoconfig"),
					resource.TestCheckResourceAttrSet("ovh_email_domain_dkim.test", "selectors.#"),
					resource.TestCheckResourceAttrSet("ovh_email_domain_dkim.test", "selectors.0.selector_name"),
					resource.TestCheckResourceAttrSet("ovh_email_domain_dkim.test", "selectors.0.cname"),
					resource.TestCheckResourceAttrSet("ovh_email_domain_dkim.test", "selectors.0.status"),
				),
			},
			// Import
			{
				ResourceName:      "ovh_email_domain_dkim.test",
				ImportState:       true,
				ImportStateId:     domain,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccEmailDomainDkimConfig_noDomain() string {
	return `
resource "ovh_email_domain_dkim" "test" {
}`
}

func testAccEmailDomainDkimConfig(domain string) string {
	return fmt.Sprintf(`
resource "ovh_email_domain_dkim" "test" {
  domain = %q
}`, domain)
}
