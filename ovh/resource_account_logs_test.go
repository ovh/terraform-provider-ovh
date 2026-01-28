package ovh

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Tests valid configurations to create account logs subscriptions
func TestAccAccountLogs_basic(t *testing.T) {
	const resourceName = "ovh_account_logs.audit_logs"

	streamID := os.Getenv("OVH_LDP_STREAM_ID")
	if streamID == "" {
		t.Skip("OVH_LDP_STREAM_ID environment variable is not set")
	}

	config := fmt.Sprintf(`
	resource "ovh_account_logs" "audit_logs" {
		log_type      = "audit"
		stream_id     = "%s"
		kind          = "default"
	}`, streamID)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCredentials(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			// Create the object
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "log_type", "audit"),
					resource.TestCheckResourceAttr(resourceName, "kind", "default"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
		},
	})
}

// Tests valid activity log subscription
func TestAccAccountLogs_activity(t *testing.T) {
	const resourceName = "ovh_account_logs.activity_logs"

	streamID := os.Getenv("OVH_LDP_STREAM_ID")
	if streamID == "" {
		t.Skip("OVH_LDP_STREAM_ID environment variable is not set")
	}

	config := fmt.Sprintf(`
	resource "ovh_account_logs" "activity_logs" {
		log_type      = "activity"
		stream_id = "%s"
		kind          = "default"
	}`, streamID)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCredentials(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			// Create the object
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "log_type", "activity"),
					resource.TestCheckResourceAttr(resourceName, "kind", "default"),
				),
			},
		},
	})
}

// Tests valid access_policy log subscription
func TestAccAccountLogs_accessPolicy(t *testing.T) {
	const resourceName = "ovh_account_logs.access_policy_logs"

	streamID := os.Getenv("OVH_LDP_STREAM_ID")
	if streamID == "" {
		t.Skip("OVH_LDP_STREAM_ID environment variable is not set")
	}

	config := fmt.Sprintf(`
	resource "ovh_account_logs" "access_policy_logs" {
		log_type      = "access_policy"
		stream_id     = "%s"
		kind          = "default"
	}`, streamID)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCredentials(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			// Create the object
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "log_type", "access_policy"),
					resource.TestCheckResourceAttr(resourceName, "kind", "default"),
				),
			},
		},
	})
}

// Test invalid log type configuration
func TestAccAccountLogs_invalidLogType(t *testing.T) {
	const configInvalidLogType = `
	resource "ovh_account_logs" "logs" {
		log_type      = "invalid_log_type"
		stream_id     = "xxx"
		kind          = "default"
	}`

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCredentials(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      configInvalidLogType,
				ExpectError: regexp.MustCompile("unsupported log type"),
			},
		},
	})
}

// Test missing required arguments
func TestAccAccountLogs_configMissingArguments(t *testing.T) {
	const configMissingStreamId = `
	resource "ovh_account_logs" "logs" {
		log_type = "audit"
		kind     = "default"
	}`

	const configMissingKind = `
	resource "ovh_account_logs" "logs" {
		log_type      = "audit"
		stream_id     = "xxx"
	}`

	const configMissingLogType = `
	resource "ovh_account_logs" "logs" {
		stream_id     = "xxx"
		kind          = "default"
	}`

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCredentials(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      configMissingStreamId,
				ExpectError: regexp.MustCompile("Missing required argument"),
			},
			{
				Config:      configMissingKind,
				ExpectError: regexp.MustCompile("Missing required argument"),
			},
			{
				Config:      configMissingLogType,
				ExpectError: regexp.MustCompile("Missing required argument"),
			},
		},
	})
}
