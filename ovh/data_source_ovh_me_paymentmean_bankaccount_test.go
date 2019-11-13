package ovh

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccMePaymentmeanBankaccountDataSource_basic(t *testing.T) {
	// ovh bank account payment mean is not mandatory
	// this datasource is tested only if env var `OVH_TEST_BANKACCOUNT`
	// is set to "1"
	v := os.Getenv("OVH_TEST_BANKACCOUNT")
	if v == "1" {
		resource.Test(t, resource.TestCase{
			PreCheck:  func() { testAccPreCheckMePaymentMean(t) },
			Providers: testAccProviders,
			Steps: []resource.TestStep{
				{
					Config: testAccMePaymentmeanBankaccountDatasourceConfig,
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr(
							"data.ovh_me_paymentmean_bankaccount.ba", "default", "true"),
					),
				},
			},
		})
	}
}

const testAccMePaymentmeanBankaccountDatasourceConfig = `
data "ovh_me_paymentmean_bankaccount" "ba" {
 use_default = true
}
`
