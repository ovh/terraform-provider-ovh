package ovh

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func init() {
	resource.AddTestSweepers("ovh_me_logs_audit_log_subscription", &resource.Sweeper{
		Name: "ovh_me_logs_audit_log_subscription",
		F: func(region string) error {
			return testSweepAccountLogSubscription(region, "/me/logs/audit/log/subscription")
		},
	})
}

func TestAccMeLogsAuditLogSubscription_basic(t *testing.T) {
	testAccAccountLogSubscription(t, "ovh_me_logs_audit_log_subscription")
}
