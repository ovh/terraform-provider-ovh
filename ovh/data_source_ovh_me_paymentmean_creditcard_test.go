package ovh

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccMePaymentmeanCreditcardDataSource_basic(t *testing.T) {
	// ovh credit card payment mean is not mandatory
	// this datasource is tested only if env var `OVH_TEST_CREDITCARD`
	// is set to "1"
	v := os.Getenv("OVH_TEST_CREDITCARD")
	if v == "1" {
		resource.Test(t, resource.TestCase{
			PreCheck:  func() { testAccPreCheckMePaymentMean(t) },
			Providers: testAccProviders,
			Steps: []resource.TestStep{
				{
					Config: testAccMePaymentmeanCreditcardDatasourceConfig,
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr(
							"data.ovh_me_paymentmean_creditcard.cc", "default", "true"),
					),
				},
			},
		})
	}
}

const testAccMePaymentmeanCreditcardDatasourceConfig = `
data "ovh_me_paymentmean_creditcard" "cc" {
 use_default = true
}
`
