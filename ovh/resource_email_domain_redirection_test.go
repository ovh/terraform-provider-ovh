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
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// Redirections are user data, so unlike a single toggle the sweeper must not
// remove everything it finds. Only entries created by these tests, recognised
// by this prefix, are eligible.
const testAccEmailDomainRedirectionPrefix = "tf-acc-test-redirection"

func init() {
	resource.AddTestSweepers("ovh_email_domain_redirection", &resource.Sweeper{
		Name: "ovh_email_domain_redirection",
		F:    testSweepEmailDomainRedirection,
	})
}

func testSweepEmailDomainRedirection(region string) error {
	client, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting client: %s", err)
	}

	domain := os.Getenv("OVH_EMAIL_DOMAIN_TEST")
	if domain == "" {
		log.Print("[DEBUG] OVH_EMAIL_DOMAIN_TEST is not set. No email domain redirections to sweep")
		return nil
	}

	var ids []string
	if err := client.Get(fmt.Sprintf("/email/domain/%s/redirection", url.PathEscape(domain)), &ids); err != nil {
		return fmt.Errorf("Error calling /email/domain/%s/redirection:\n\t %q", domain, err)
	}

	for _, id := range ids {
		var redirection struct {
			From string `json:"from"`
		}
		endpoint := fmt.Sprintf("/email/domain/%s/redirection/%s", url.PathEscape(domain), url.PathEscape(id))
		if err := client.Get(endpoint, &redirection); err != nil {
			return fmt.Errorf("Error calling %s:\n\t %q", endpoint, err)
		}

		if !strings.HasPrefix(redirection.From, testAccEmailDomainRedirectionPrefix) {
			continue
		}

		log.Printf("[DEBUG] Deleting email domain redirection %s", redirection.From)
		if err := client.Delete(endpoint, nil); err != nil {
			return fmt.Errorf("Error calling delete on %s:\n\t %q", endpoint, err)
		}
	}

	return nil
}

func TestAccEmailDomainRedirection_Basic(t *testing.T) {
	domain := os.Getenv("OVH_EMAIL_DOMAIN_TEST")
	from := fmt.Sprintf("%s@%s", testAccEmailDomainRedirectionPrefix, domain)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckEmailDomain(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Error when from is missing
			{
				Config:      testAccEmailDomainRedirectionConfig_noFrom(domain),
				ExpectError: regexp.MustCompile(`The argument "from" is required`),
			},
			// Successful creation
			{
				Config: testAccEmailDomainRedirectionConfig(domain, from, "first-target@example.com", false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_email_domain_redirection.test", "domain", domain),
					resource.TestCheckResourceAttr("ovh_email_domain_redirection.test", "from", from),
					resource.TestCheckResourceAttr("ovh_email_domain_redirection.test", "to", "first-target@example.com"),
					resource.TestCheckResourceAttr("ovh_email_domain_redirection.test", "local_copy", "false"),
					resource.TestCheckResourceAttrSet("ovh_email_domain_redirection.test", "id"),
				),
			},
			// Changing the target goes through changeRedirection rather than
			// replacing the resource, and reassigns the id.
			{
				Config: testAccEmailDomainRedirectionConfig(domain, from, "second-target@example.com", false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_email_domain_redirection.test", "from", from),
					resource.TestCheckResourceAttr("ovh_email_domain_redirection.test", "to", "second-target@example.com"),
					resource.TestCheckResourceAttrSet("ovh_email_domain_redirection.test", "id"),
				),
			},
			// Import. local_copy is never returned by the API, so it cannot be
			// verified against imported state.
			{
				ResourceName:            "ovh_email_domain_redirection.test",
				ImportState:             true,
				ImportStateIdFunc:       testAccEmailDomainRedirectionImportID("ovh_email_domain_redirection.test"),
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"local_copy"},
			},
		},
	})
}

func testAccEmailDomainRedirectionImportID(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("resource %s not found", resourceName)
		}

		return fmt.Sprintf("%s/%s", rs.Primary.Attributes["domain"], rs.Primary.Attributes["id"]), nil
	}
}

func testAccEmailDomainRedirectionConfig_noFrom(domain string) string {
	return fmt.Sprintf(`
resource "ovh_email_domain_redirection" "test" {
  domain = %q
  to     = "target@example.com"
}`, domain)
}

func testAccEmailDomainRedirectionConfig(domain, from, to string, localCopy bool) string {
	return fmt.Sprintf(`
resource "ovh_email_domain_redirection" "test" {
  domain     = %q
  from       = %q
  to         = %q
  local_copy = %t
}`, domain, from, to, localCopy)
}
