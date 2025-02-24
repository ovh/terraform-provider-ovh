package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccResourceDomainName_basic(t *testing.T) {
	domain := os.Getenv("OVH_TESTACC_ORDER_DOMAIN")
	config := fmt.Sprintf(`
		resource "ovh_domain_name" "domain" {	  
			domain_name = "%s"
		}`,
		domain,
	)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckOrderDomain(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_domain_name.domain", "id", domain),
					resource.TestCheckResourceAttrSet("ovh_domain_name.domain", "checksum"),
					resource.TestCheckResourceAttr("ovh_domain_name.domain", "current_state.dns_configuration.name_servers.#", "2"),
				),
			},
			{
				ResourceName:      "ovh_domain_name.name",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateId:     domain,
			},
		},
	})
}
