package ovh

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func init() {
	resource.AddTestSweepers("ovh_me_api_log_subscription", &resource.Sweeper{
		Name: "ovh_me_api_log_subscription",
		F: func(region string) error {
			return testSweepAccountLogSubscription(region, "/me/api/log/subscription")
		},
	})
}

func TestAccMeApiLogSubscription_basic(t *testing.T) {
	testAccAccountLogSubscription(t, "ovh_me_api_log_subscription")
}
