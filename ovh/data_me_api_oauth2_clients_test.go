package ovh

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Tests valid configurations to read all oauth2 clients
func TestAccMeApiOauth2Clients_data(t *testing.T) {
	// Successful test with app required parameters,
	// that validates the use of the data instruction
	// We create a resource beforehand, to make sure that at least one service account exists
	const okConfigClientCredentials = `
	resource "ovh_me_api_oauth2_client" "service_account_1" {
		description = "tf acc test client credentials"
		name        = "tf acc test client credentials"
		flow        = "CLIENT_CREDENTIALS"
	}
	# A data source listing all the client ids in the account
	data "ovh_me_api_oauth2_clients" "all_clients_ref" {
		depends_on = [ovh_me_api_oauth2_client.service_account_1]
	}

	output "is_array_length_gt_0" {
		value = length(data.ovh_me_api_oauth2_clients.all_clients_ref.client_ids) > 0
	}
	`
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCredentials(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			// Create the object, check that the client secret is not empty after creation
			{
				Config: okConfigClientCredentials,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckOutput(
						"is_array_length_gt_0", "true",
					),
				),
			},
		},
	})
}
