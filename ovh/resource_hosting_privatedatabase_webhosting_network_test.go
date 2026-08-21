package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const testAccHostingPrivateDatabaseWebhostingNetworkBasic = `
resource "ovh_hosting_privatedatabase_webhosting_network" "network" {
    service_name = "%s"
    enabled      = %t
}
`

func TestAccHostingPrivateDatabaseWebhostingNetwork_basic(t *testing.T) {
	serviceName := os.Getenv("OVH_HOSTING_PRIVATEDATABASE_SERVICE_TEST")
	resourceName := "ovh_hosting_privatedatabase_webhosting_network.network"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckHostingPrivateDatabaseWebhostingNetwork(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccHostingPrivateDatabaseWebhostingNetworkBasic, serviceName, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "service_name", serviceName),
					resource.TestCheckResourceAttr(resourceName, "enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "status", "disabled"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateId:     serviceName,
				ImportStateVerify: true,
			},
			{
				Config: fmt.Sprintf(testAccHostingPrivateDatabaseWebhostingNetworkBasic, serviceName, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "status", "enabled"),
				),
			},
		},
	})
}
