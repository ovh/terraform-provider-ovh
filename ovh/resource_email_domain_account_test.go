package ovh

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"regexp"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func init() {
	resource.AddTestSweepers("ovh_email_domain_account", &resource.Sweeper{
		Name: "ovh_email_domain_account",
		F:    testSweepEmailDomainAccount,
	})
}

func testSweepEmailDomainAccount(region string) error {
	client, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting client: %s", err)
	}

	domain := os.Getenv("OVH_EMAIL_DOMAIN_TEST")
	if domain == "" {
		log.Print("[DEBUG] OVH_EMAIL_DOMAIN_TEST is not set. No email domain accounts to sweep")
		return nil
	}

	accountNames := make([]string, 0)
	if err := client.Get(fmt.Sprintf("/email/domain/%s/account", url.PathEscape(domain)), &accountNames); err != nil {
		return fmt.Errorf("Error calling /email/domain/%s/account:\n\t %q", domain, err)
	}

	if len(accountNames) == 0 {
		log.Print("[DEBUG] No email domain accounts to sweep")
		return nil
	}

	for _, accountName := range accountNames {
		log.Printf("[DEBUG] Deleting email domain account %s@%s", accountName, domain)

		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			if err := client.Delete(
				fmt.Sprintf("/email/domain/%s/account/%s", url.PathEscape(domain), url.PathEscape(accountName)),
				nil,
			); err != nil {
				return resource.RetryableError(err)
			}
			return nil
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func TestAccEmailDomainAccount_Basic(t *testing.T) {
	domain := os.Getenv("OVH_EMAIL_DOMAIN_TEST")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckEmailDomain(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Error when domain is missing
			{
				Config:      testAccEmailDomainAccountConfig_noDomain("testaccount", "P@ssw0rd1234!"),
				ExpectError: regexp.MustCompile(`The argument "domain" is required`),
			},
			// Error when account_name is missing
			{
				Config:      testAccEmailDomainAccountConfig_noAccountName(domain, "P@ssw0rd1234!"),
				ExpectError: regexp.MustCompile(`The argument "account_name" is required`),
			},
			// Successful creation
			{
				Config: testAccEmailDomainAccountConfig(domain, "testaccount", "P@ssw0rd1234!", "Test account", 5000000000),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_email_domain_account.test", "domain", domain),
					resource.TestCheckResourceAttr("ovh_email_domain_account.test", "account_name", "testaccount"),
					resource.TestCheckResourceAttr("ovh_email_domain_account.test", "password", "P@ssw0rd1234!"),
					resource.TestCheckResourceAttr("ovh_email_domain_account.test", "description", "Test account"),
					resource.TestCheckResourceAttr("ovh_email_domain_account.test", "size", "5000000000"),
					resource.TestCheckResourceAttr("ovh_email_domain_account.test", "email", "testaccount@"+domain),
					resource.TestCheckResourceAttr("ovh_email_domain_account.test", "is_blocked", "false"),
				),
			},
			// Update description
			{
				Config: testAccEmailDomainAccountConfig(domain, "testaccount", "P@ssw0rd1234!", "Updated description", 5000000000),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_email_domain_account.test", "description", "Updated description"),
					resource.TestCheckResourceAttr("ovh_email_domain_account.test", "account_name", "testaccount"),
					// The id used to be lost on update
					resource.TestCheckResourceAttr("ovh_email_domain_account.test", "id", domain+"/testaccount"),
				),
			},
			// Update password
			{
				Config: testAccEmailDomainAccountConfig(domain, "testaccount", "N3wP@ssw0rd5678!", "Updated description", 5000000000),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_email_domain_account.test", "password", "N3wP@ssw0rd5678!"),
					resource.TestCheckResourceAttr("ovh_email_domain_account.test", "account_name", "testaccount"),
				),
			},
			// Import: the password is never returned by the API
			{
				ResourceName:            "ovh_email_domain_account.test",
				ImportState:             true,
				ImportStateId:           domain + "/testaccount",
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password"},
			},
		},
	})
}

// An account managed without a password, as after an import, keeps the one it
// has: removing the password from the configuration changes nothing upstream,
// and the account can still be read, updated and imported.
func TestAccEmailDomainAccount_withoutPassword(t *testing.T) {
	domain := os.Getenv("OVH_EMAIL_DOMAIN_TEST")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckEmailDomain(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Creating an account still needs a password
			{
				Config:      testAccEmailDomainAccountConfig_noPassword(domain, "testnopassword", "No password"),
				ExpectError: regexp.MustCompile(`A password is required to create an email account`),
			},
			{
				Config: testAccEmailDomainAccountConfig(domain, "testnopassword", "P@ssw0rd1234!", "No password", 5000000000),
				Check:  resource.TestCheckResourceAttr("ovh_email_domain_account.test", "password", "P@ssw0rd1234!"),
			},
			// Leaving the password out keeps the current one
			{
				Config: testAccEmailDomainAccountConfig_noPassword(domain, "testnopassword", "Still no password"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckNoResourceAttr("ovh_email_domain_account.test", "password"),
					resource.TestCheckResourceAttr("ovh_email_domain_account.test", "description", "Still no password"),
				),
			},
			{
				ResourceName:      "ovh_email_domain_account.test",
				ImportState:       true,
				ImportStateId:     domain + "/testnopassword",
				ImportStateVerify: true,
			},
			// Setting it again changes it
			{
				Config: testAccEmailDomainAccountConfig(domain, "testnopassword", "N3wP@ssw0rd5678!", "Still no password", 5000000000),
				Check:  resource.TestCheckResourceAttr("ovh_email_domain_account.test", "password", "N3wP@ssw0rd5678!"),
			},
		},
	})
}

func testAccEmailDomainAccountConfig_noPassword(domain, accountName, description string) string {
	return fmt.Sprintf(`
resource "ovh_email_domain_account" "test" {
  domain       = %q
  account_name = %q
  description  = %q
}`, domain, accountName, description)
}

func testAccEmailDomainAccountConfig_noDomain(accountName, password string) string {
	return fmt.Sprintf(`
resource "ovh_email_domain_account" "test" {
  account_name = %q
  password     = %q
}`, accountName, password)
}

func testAccEmailDomainAccountConfig_noAccountName(domain, password string) string {
	return fmt.Sprintf(`
resource "ovh_email_domain_account" "test" {
  domain   = %q
  password = %q
}`, domain, password)
}

func testAccEmailDomainAccountConfig(domain, accountName, password, description string, size int64) string {
	return fmt.Sprintf(`
resource "ovh_email_domain_account" "test" {
  domain       = %q
  account_name = %q
  password     = %q
  description  = %q
  size         = %d
}`, domain, accountName, password, description, size)
}
