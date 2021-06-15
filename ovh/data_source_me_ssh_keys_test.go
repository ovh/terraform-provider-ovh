package ovh

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccMeSshKeysDataSource_basic(t *testing.T) {
	sshKey1Name := acctest.RandomWithPrefix(test_prefix)
	sshKey1 := "ssh-ed25519 AAAAC3NzaC1xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
	sshKey2Name := acctest.RandomWithPrefix(test_prefix)
	sshKey2 := "ssh-ed25519 AAAAC3NzaC1yyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyy"

	preSetup := fmt.Sprintf(
		testAccMeSshKeysDatasourceConfig_preSetup,
		sshKey1Name,
		sshKey1,
		sshKey2Name,
		sshKey2,
	)
	config := fmt.Sprintf(
		testAccMeSshKeysDatasourceConfig_keys,
		sshKey1Name,
		sshKey1,
		sshKey2Name,
		sshKey2,
	)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCredentials(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: preSetup,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"ovh_me_ssh_key.key_1", "key_name", sshKey1Name),
					resource.TestCheckResourceAttr(
						"ovh_me_ssh_key.key_1", "key", sshKey1),
					resource.TestCheckResourceAttr(
						"ovh_me_ssh_key.key_2", "key_name", sshKey2Name),
					resource.TestCheckResourceAttr(
						"ovh_me_ssh_key.key_2", "key", sshKey2),
				),
			}, {
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckOutput(
						"keys_present", "true"),
				),
			},
		},
	})
}

const testAccMeSshKeysDatasourceConfig_preSetup = `
resource "ovh_me_ssh_key" "key_1" {
	key_name = "%s"
	key      = "%s"
}

resource "ovh_me_ssh_key" "key_2" {
	key_name = "%s"
	key      = "%s"
}
`

const testAccMeSshKeysDatasourceConfig_keys = `
resource "ovh_me_ssh_key" "key_1" {
	key_name = "%s"
	key      = "%s"
}

resource "ovh_me_ssh_key" "key_2" {
	key_name = "%s"
	key      = "%s"
}

data "ovh_me_ssh_keys" "keys" {}

output "keys_present" {
    value = tostring(contains(data.ovh_me_ssh_keys.keys.names, ovh_me_ssh_key.key_1.key_name) && contains(data.ovh_me_ssh_keys.keys.names, ovh_me_ssh_key.key_2.key_name))
}
`
