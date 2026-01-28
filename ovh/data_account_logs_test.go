package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Tests reading account logs subscriptions via data source
func TestAccDataSourceAccountLogs_basic(t *testing.T) {
	const resourceName = "ovh_account_logs.audit_logs"
	const dataSourceName = "data.ovh_account_logs.audit_logs"

	streamID := os.Getenv("OVH_LDP_STREAM_ID")
	if streamID == "" {
		t.Skip("OVH_LDP_STREAM_ID environment variable is not set")
	}

	config := fmt.Sprintf(`
	resource "ovh_account_logs" "audit_logs" {
		log_type      = "audit"
		stream_id = "%s"
		kind          = "default"
	}

	data "ovh_account_logs" "audit_logs" {
		log_type        = "audit"
		subscription_id = ovh_account_logs.audit_logs.id
	}`, streamID)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCredentials(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "log_type", "audit"),
					resource.TestCheckResourceAttr(resourceName, "kind", "default"),
					// Verify data source attributes match resource
					resource.TestCheckResourceAttrPair(dataSourceName, "id", resourceName, "id"),
					resource.TestCheckResourceAttrPair(dataSourceName, "log_type", resourceName, "log_type"),
					resource.TestCheckResourceAttrPair(dataSourceName, "kind", resourceName, "kind"),
				),
			},
		},
	})
}

// Tests reading activity logs via data source
func TestAccDataSourceAccountLogs_activity(t *testing.T) {
	const dataSourceName = "data.ovh_account_logs.activity_logs"

	streamID := os.Getenv("OVH_LDP_STREAM_ID")
	if streamID == "" {
		t.Skip("OVH_LDP_STREAM_ID environment variable is not set")
	}

	config := fmt.Sprintf(`
	resource "ovh_account_logs" "activity_logs" {
		log_type      = "activity"
		stream_id     = "%s"
		kind          = "default"
	}

	data "ovh_account_logs" "activity_logs" {
		log_type        = "activity"
		subscription_id = ovh_account_logs.activity_logs.id
	}`, streamID)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCredentials(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "log_type", "activity"),
					resource.TestCheckResourceAttr(dataSourceName, "kind", "default"),
				),
			},
		},
	})
}
