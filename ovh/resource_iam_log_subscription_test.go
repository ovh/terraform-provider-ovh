package ovh

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func init() {
	resource.AddTestSweepers("ovh_iam_log_subscription", &resource.Sweeper{
		Name: "ovh_iam_log_subscription",
		F: func(region string) error {
			return testSweepAccountLogSubscription(region, "/v2/iam/log/subscription")
		},
	})
}

func TestAccIamLogSubscription_basic(t *testing.T) {
	testAccAccountLogSubscription(t, "ovh_iam_log_subscription")
}
