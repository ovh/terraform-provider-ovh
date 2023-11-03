package ovh

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// Returns an error if a string is empty,
// Used in the calls to resource.TestCheckResourceAttrWith() below
func apiOauth2ClientStringNotEmpty(s string) error {
	if len(s) > 0 {
		return nil
	}
	return fmt.Errorf("Value is empty")
}

// Tests valid configurations to create oauth2 clients
func TestAccMeApiOauth2Client_importBasic(t *testing.T) {
	const resourceName1 = "ovh_me_api_oauth2_client.service_account_1"
	const resourceName2 = "ovh_me_api_oauth2_client.service_account_2"

	// Successful test with app required parameters
	const okConfigClientCredentials = `
	resource "ovh_me_api_oauth2_client" "service_account_1" {
		description = "tf acc test client credentials"
		name        = "tf acc test client credentials"
		flow = "CLIENT_CREDENTIALS"
	}`
	const okConfigAuthorizationCode = `
	resource "ovh_me_api_oauth2_client" "service_account_2" {
		description = "tf acc test authorization code"
		name        = "tf acc test authorization code"
		flow = "AUTHORIZATION_CODE"
		callback_urls = ["https://localhost:8080"]
	}`
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCredentials(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			// Create the object, check that the client secret is not empty after creation
			{
				Config: okConfigClientCredentials,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrWith(resourceName1, "client_id", apiOauth2ClientStringNotEmpty),
					resource.TestCheckResourceAttrWith(resourceName1, "client_secret", apiOauth2ClientStringNotEmpty),
				),
			},
			// Verify that state matches the imported resource
			{
				ResourceName:            resourceName1,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"client_secret"}, // Client secrets cannot be imported
			},
			// Update the object with the second configuration, check that the client secret is not empty after creation
			{
				Config: okConfigAuthorizationCode,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrWith(resourceName2, "client_id", apiOauth2ClientStringNotEmpty),
					resource.TestCheckResourceAttrWith(resourceName2, "client_secret", apiOauth2ClientStringNotEmpty),
				),
			},
			// Verify that state matches the imported resource
			{
				ResourceName:            resourceName2,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"client_secret"}, // Client secrets cannot be imported
			},
		},
	})
}

// Test invalid oauth2 clients configurations
func TestAccMeApiOauth2Client_configMissingArguments(t *testing.T) {
	const configMissingName = `
	resource "ovh_me_api_oauth2_client" "service_account_1" {
		description = "tf acc import test"
		flow = "AUTHORIZATION_CODE"
	}`
	const configMissingDescription = `
	resource "ovh_me_api_oauth2_client" "service_account_1" {
		name        = "tf acc test authorization code"
	}`
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCredentials(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			// Check that we cannot create a resource with a missing name
			{
				Config:      configMissingName,
				ExpectError: regexp.MustCompile("Missing required argument"),
			},
			// Check that we cannot create a resource with a missing description
			{
				Config:      configMissingDescription,
				ExpectError: regexp.MustCompile("Missing required argument"),
			},
		},
	})
}
