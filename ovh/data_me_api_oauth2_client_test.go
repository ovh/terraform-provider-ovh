package ovh

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// Tests valid configurations to create oauth2 clients
func TestAccMeApiOauth2Client_data(t *testing.T) {
	const resourceName = "ovh_me_api_oauth2_client.service_account_1"

	// Successful test with app required parameters,
	// that validates the use of the data instruction
	const okConfigClientCredentials = `
	resource "ovh_me_api_oauth2_client" "service_account_1" {
		description = "tf acc test client credentials"
		name        = "tf acc test client credentials"
		flow        = "CLIENT_CREDENTIALS"
	}
	data "ovh_me_api_oauth2_client" "service_account_1" {
		client_id  = ovh_me_api_oauth2_client.service_account_1.client_id
		depends_on = [ovh_me_api_oauth2_client.service_account_1]
	}
	output "oauth2_client_name" {
		value = data.ovh_me_api_oauth2_client.service_account_1.name
	}`
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCredentials(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			// Create the object, check that the client secret is not empty after creation
			{
				Config: okConfigClientCredentials,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrWith(resourceName, "client_id", apiOauth2ClientStringNotEmpty),
					resource.TestCheckResourceAttrWith(resourceName, "client_secret", apiOauth2ClientStringNotEmpty),
					resource.TestCheckOutput(
						"oauth2_client_name", "tf acc test client credentials",
					),
				),
			},
		},
	})
}
