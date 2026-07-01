package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccEmailDomainAccountsDataSource_basic(t *testing.T) {
	domain := os.Getenv("OVH_EMAIL_DOMAIN_TEST")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckCredentials(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					data "ovh_email_domain_accounts" "accounts" {
						domain = %q
					}`, domain),
				Check: resource.TestCheckResourceAttrSet("data.ovh_email_domain_accounts.accounts", "accounts.#"),
			},
		},
	})
}
