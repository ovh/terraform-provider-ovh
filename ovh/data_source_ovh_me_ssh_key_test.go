package ovh

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccMeSshKeyDataSource_basic(t *testing.T) {
	sshKeyName := acctest.RandomWithPrefix(test_prefix)
	sshKey := "ssh-ed25519 AAAAC3NzaC1yyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyy"
	config := fmt.Sprintf(testAccMeSshKeyDatasourceConfig, sshKeyName, sshKey)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCredentials(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.ovh_me_ssh_key.key_1", "key_name", sshKeyName),
					resource.TestCheckResourceAttr(
						"data.ovh_me_ssh_key.key_1", "key", sshKey),
				),
			},
		},
	})
}

const testAccMeSshKeyDatasourceConfig = `
resource "ovh_me_ssh_key" "key_1" {
	key_name = "%s"
	key      = "%s"
}

data "ovh_me_ssh_key" "key_1" {
  key_name = ovh_me_ssh_key.key_1.key_name
  depends_on    = [ovh_me_ssh_key.key_1]
}
`
