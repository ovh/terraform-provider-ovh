package ovh

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccMeDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCredentials(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMeDatasourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(
						"data.ovh_me.my_account",
						"country",
					),
					resource.TestCheckResourceAttrSet(
						"data.ovh_me.my_account",
						"currency.0.code",
					),
					resource.TestCheckResourceAttrSet(
						"data.ovh_me.my_account",
						"currency.0.symbol",
					),
					resource.TestCheckResourceAttrSet(
						"data.ovh_me.my_account",
						"email",
					),
					resource.TestCheckResourceAttrSet(
						"data.ovh_me.my_account",
						"legalform",
					),
					resource.TestCheckResourceAttrSet(
						"data.ovh_me.my_account",
						"nichandle",
					),
					resource.TestCheckResourceAttrSet(
						"data.ovh_me.my_account",
						"ovh_company",
					),
					resource.TestCheckResourceAttrSet(
						"data.ovh_me.my_account",
						"ovh_subsidiary",
					),
					resource.TestCheckResourceAttrSet(
						"data.ovh_me.my_account",
						"state",
					),
				),
			},
		},
	})
}

const testAccMeDatasourceConfig = `
data "ovh_me" "my_account" {}
`
