package ovh

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDedicatedInstallationTemplatesDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCredentials(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: "data ovh_dedicated_installation_templates templates {}",
				Check: resource.TestCheckResourceAttrSet(
					"data.ovh_dedicated_installation_templates.templates",
					"result.#",
				),
			},
		},
	})
}
