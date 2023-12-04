package ovh

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
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
func TestAccMeApiOauth2Client_basic(t *testing.T) {
	const resourceName1 = "ovh_me_api_oauth2_client.service_account_1"
	const resourceName2 = "ovh_me_api_oauth2_client.service_account_2"

	// Successful test with app required parameters
	const okConfigClientCredentials = `
	resource "ovh_me_api_oauth2_client" "service_account_1" {
		description = "tf acc test client credentials"
		name        = "tf acc test client credentials"
		flow        = "CLIENT_CREDENTIALS"
	}`
	const okConfigAuthorizationCode = `
	resource "ovh_me_api_oauth2_client" "service_account_2" {
		description   = "tf acc test authorization code"
		name          = "tf acc test authorization code"
		flow          = "AUTHORIZATION_CODE"
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
			// Verify that the state matches the resource imported with its client_id
			{
				ResourceName:            resourceName1,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"client_secret"}, // Client secrets cannot be imported using client id only
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
				ImportStateVerifyIgnore: []string{"client_secret"}, // Client secrets cannot be imported using client id only
			},
		},
	})
}

// Test invalid oauth2 clients configurations
func TestAccMeApiOauth2Client_configMissingArguments(t *testing.T) {
	const configMissingName = `
	resource "ovh_me_api_oauth2_client" "service_account_1" {
		description = "tf acc import test"
		flow        = "AUTHORIZATION_CODE"
	}`
	const configMissingDescription = `
	resource "ovh_me_api_oauth2_client" "service_account_1" {
		name = "tf acc test authorization code"
		flow = "AUTHORIZATION_CODE"
	}`
	const configMissingFlow = `
	resource "ovh_me_api_oauth2_client" "service_account_1" {
		name        = "tf acc test authorization code"
		description = "tf acc import test"
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
			// Check that we cannot create a resource with a missing flow
			{
				Config:      configMissingFlow,
				ExpectError: regexp.MustCompile("Missing required argument"),
			},
		},
	})
}

// Test invalid oauth2 clients import IDs
func TestAccMeApiOauth2Client_importBasic(t *testing.T) {
	const resourceName = "ovh_me_api_oauth2_client.service_account_1"
	const okConfigClientCredentials = `
	resource "ovh_me_api_oauth2_client" "service_account_1" {
		description = "tf acc test client credentials"
		name        = "tf acc test client credentials"
		flow        = "CLIENT_CREDENTIALS"
	}`
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCredentials(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: okConfigClientCredentials,
			},
			// Verify that the state matches the resource imported with its client_id and client_secret separated by a pipe
			{
				ResourceName: resourceName,
				ImportStateIdFunc: func(state *terraform.State) (string, error) {
					resource_id := state.RootModule().Resources[resourceName].Primary.Attributes["id"]
					resource_secret := state.RootModule().Resources[resourceName].Primary.Attributes["client_secret"]
					return resource_id + "|" + resource_secret, nil
				},
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Check that importing a resource fails when its client_id contains a pipe but is not formatted properly
			{
				ImportState:       true,
				ImportStateVerify: true,
				ResourceName:      resourceName,
				ImportStateId:     "fake_client|fake_secret|extra_data",
				ExpectError:       regexp.MustCompile("Resource IDs with the pipe character should be formatted as"),
			},
		},
	})
}
