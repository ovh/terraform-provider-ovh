package ovh

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type CloudProjectInstanceActiveMonthlyBillingCreateOpts struct {
}

func (p *CloudProjectInstanceActiveMonthlyBillingCreateOpts) String() string {
	return fmt.Sprintf("")
}

func (p *CloudProjectInstanceActiveMonthlyBillingCreateOpts) FromResource(d *schema.ResourceData) *CloudProjectInstanceActiveMonthlyBillingCreateOpts {
	params := &CloudProjectInstanceActiveMonthlyBillingCreateOpts{}
	return params
}

type CloudProjectInstanceActiveMonthlyBillingResponseMonthlyBilling struct {
	Since  string `json:"since"`
	Status string `json:"status"`
}

func (p *CloudProjectInstanceActiveMonthlyBillingResponseMonthlyBilling) String() string {
	return fmt.Sprintf("since: %s, status: %s", p.Since, p.Status)
}

type CloudProjectInstanceActiveMonthlyBillingResponse struct {
	MonthlyBilling *CloudProjectInstanceActiveMonthlyBillingResponseMonthlyBilling `json:"monthlyBilling"`
}

func (p *CloudProjectInstanceActiveMonthlyBillingResponse) String() string {
	return fmt.Sprintf("monthlyBilling: %s", p.MonthlyBilling)
}
