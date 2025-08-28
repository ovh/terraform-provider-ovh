package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Step 1 config: create secret version 1 and read latest with data.
// NOTE: resource and data source must be prefixed with provider name 'ovh_'.
const testAccOkmsSecretDataSourceConfigV1 = `
resource "ovh_okms_secret" "test" {
	okms_id = "%s"
	path    = "%s"
	version = {
		data = jsonencode({ initial = "v1" })
	}
}

data "ovh_okms_secret" "latest" {
	okms_id      = ovh_okms_secret.test.okms_id
	path         = ovh_okms_secret.test.path
	include_data = true
}
`

// Step 2 config: update same resource (cas=1) to create version 2, then read latest, explicit v1 and v2.
const testAccOkmsSecretDataSourceConfigV2 = `
resource "ovh_okms_secret" "test" {
	okms_id = "%s"
	path    = "%s"
	cas     = 1
	version = {
		data = jsonencode({ initial = "v1", second = "v2" })
	}
}

data "ovh_okms_secret" "latest" {
	okms_id      = ovh_okms_secret.test.okms_id
	path         = ovh_okms_secret.test.path
	include_data = true
}

data "ovh_okms_secret" "v1" {
	okms_id      = ovh_okms_secret.test.okms_id
	path         = ovh_okms_secret.test.path
	version      = 1
	include_data = true
}

data "ovh_okms_secret" "v2" {
	okms_id      = ovh_okms_secret.test.okms_id
	path         = ovh_okms_secret.test.path
	version      = 2
	include_data = true
}
`

func TestAccOkmsSecretDataSource_latestAndVersions(t *testing.T) {
	okmsID := os.Getenv("OVH_OKMS")
	if okmsID == "" {
		checkEnvOrSkip(t, "OVH_OKMS")
	}
	path := fmt.Sprintf("tfacc-%s", acctest.RandString(6))

	configV1 := fmt.Sprintf(testAccOkmsSecretDataSourceConfigV1, okmsID, path)
	configV2 := fmt.Sprintf(testAccOkmsSecretDataSourceConfigV2, okmsID, path)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckOkms(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: configV1,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.ovh_okms_secret.latest", "version", "1"),
					resource.TestCheckResourceAttrSet("data.ovh_okms_secret.latest", "data"),
				),
			},
			{
				Config: configV2,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.ovh_okms_secret.latest", "version", "2"),
					resource.TestCheckResourceAttrSet("data.ovh_okms_secret.latest", "data"),
					resource.TestCheckResourceAttr("data.ovh_okms_secret.v1", "version", "1"),
					resource.TestCheckResourceAttrSet("data.ovh_okms_secret.v1", "data"),
					resource.TestCheckResourceAttr("data.ovh_okms_secret.v2", "version", "2"),
					resource.TestCheckResourceAttrSet("data.ovh_okms_secret.v2", "data"),
				),
			},
		},
	})
}
