package ovh

import (
	"errors"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func dedicatedServerResourceTestConfig(updated bool) string {
	var (
		monitoring        = true
		noIntervention    = false
		operatingSystem   = "debian11_64"
		displayName       = "First display name"
		efiBootloaderPath = ""
	)

	if updated {
		monitoring = false
		noIntervention = true
		operatingSystem = "debian12_64"
		displayName = "Second display name"
		efiBootloaderPath = `\\efi\\debian\\grubx64.efi`
	}

	return fmt.Sprintf(`
	data "ovh_me" "account" {}

	resource "ovh_dedicated_server" "server" {
		ovh_subsidiary = data.ovh_me.account.ovh_subsidiary
		monitoring = %t
		no_intervention = %t
		os = "%s"
		display_name = "%s"
		efi_bootloader_path = "%s"
		plan = [
			{
				plan_code = "24rise01"
				duration = "P1M"
				pricing_mode = "default"

				configuration = [
					{
						label = "dedicated_datacenter"
						value = "bhs"
					},
					{
						label = "dedicated_os"
						value = "none_64.en"
					},
					{
						label = "region"
						value = "canada"
					}
				]
			}
		]

		plan_option = [
			{
				duration = "P1M"
				plan_code = "ram-32g-rise13"
				pricing_mode = "default"
				quantity = 1
			},
			{
				duration = "P1M"
				plan_code = "bandwidth-500-included-rise"
				pricing_mode = "default"
				quantity = 1
			},
			{
				duration = "P1M"
				plan_code = "softraid-2x512nvme-rise"
				pricing_mode = "default"
				quantity = 1
			},
			{
				duration = "P1M"
				plan_code = "vrack-bandwidth-100-rise-included"
				pricing_mode = "default"
				quantity = 1
			}
		]
	}
	`, monitoring, noIntervention, operatingSystem, displayName, efiBootloaderPath)
}

func TestAccDedicatedServer_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckOrderDedicatedServer(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: dedicatedServerResourceTestConfig(false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"ovh_dedicated_server.server", "monitoring", "true"),
					resource.TestCheckResourceAttr(
						"ovh_dedicated_server.server", "no_intervention", "false"),
					resource.TestCheckResourceAttr(
						"ovh_dedicated_server.server", "display_name", "First display name"),
					resource.TestCheckResourceAttr(
						"ovh_dedicated_server.server", "iam.display_name", "First display name"),
					resource.TestCheckResourceAttr(
						"ovh_dedicated_server.server", "os", "debian11_64"),
					resource.TestCheckResourceAttr(
						"ovh_dedicated_server.server", "efi_bootloader_path", ""),
				),
			},
			{
				Config: dedicatedServerResourceTestConfig(true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"ovh_dedicated_server.server", "monitoring", "false"),
					resource.TestCheckResourceAttr(
						"ovh_dedicated_server.server", "no_intervention", "true"),
					resource.TestCheckResourceAttr(
						"ovh_dedicated_server.server", "display_name", "Second display name"),
					resource.TestCheckResourceAttr(
						"ovh_dedicated_server.server", "iam.display_name", "Second display name"),
					resource.TestCheckResourceAttr(
						"ovh_dedicated_server.server", "os", "debian12_64"),
					resource.TestCheckResourceAttr(
						"ovh_dedicated_server.server", "efi_bootloader_path", "\\efi\\debian\\grubx64.efi"),
				),
			},
			{
				ResourceName:                         "ovh_dedicated_server.server",
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "service_name",
				ImportStateVerifyIgnore: []string{
					"display_name", "order", "ovh_subsidiary", "plan", "plan_option",
				},
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					service, ok := s.RootModule().Resources["ovh_dedicated_server.server"]
					if !ok {
						return "", errors.New("ovh_dedicated_server.server not found")
					}
					return service.Primary.Attributes["service_name"], nil
				},
			},
		},
	})
}
