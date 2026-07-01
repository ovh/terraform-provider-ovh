package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccEmailDomainAccountDataSource_basic(t *testing.T) {
	domain := os.Getenv("OVH_EMAIL_DOMAIN_TEST")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckCredentials(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					data "ovh_email_domain_accounts" "accounts" {
						domain = %q
					}

					data "ovh_email_domain_account" "account" {
						domain       = data.ovh_email_domain_accounts.accounts.domain
						account_name = tolist(data.ovh_email_domain_accounts.accounts.accounts)[0]
					}`, domain),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.ovh_email_domain_account.account", "email"),
					resource.TestCheckResourceAttrSet("data.ovh_email_domain_account.account", "size"),
				),
			},
		},
	})
}
