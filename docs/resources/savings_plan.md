---
subcategory : "Savings Plan"
---

# ovh_savings_plan

Create and manage an OVHcloud Savings Plan

## Example Usage

```terraform
resource "ovh_savings_plan" "plan" {
  service_name = "<public cloud project ID>"
  flavor = "Rancher"
  period = "P1M"
  size = 2
  display_name = "one_month_rancher_savings_plan"
  auto_renewal = true
}
```

## Argument Reference

The following arguments are supported:

* `service_name` - (Required) ID of the public cloud project
* `flavor` - (Required) Savings Plan flavor. The list of available flavors can be retrieved in the [next section](#available-flavors).
* `period` - (Required) Periodicity of the Savings Plan
* `size` - (Required) Size of the Savings Plan
* `display_name` - (Required) Custom display name, used in invoices
* `auto_renewal` - Whether Savings Plan should be renewed at the end of the period (defaults to false)

### Available flavors

The following flavors can be used in savings plans (the value is case-insensitive):

* Rancher: `rancher`, `rancher_standard`, `rancher_ovhcloud_edition` (please note that `rancher` is an alias for `rancher_standard`)
* General purpose instances: `b3-8`, `b3-16`, `b3-32`, `b3-64`, `b3-128`, `b3-256`
* Compute optimized instances: `c3-4`, `c3-8`, `c3-16`, `c3-32`, `c3-64`, `c3-128`
* Memory optimized instances: `r3-16`, `r3-32`, `r3-64`, `r3-128`, `r3-256`, `r3-512`

## Attributes Reference

The following attributes are exported:

* `id` - ID of the Savings Plan
* `service_name` - ID of the public cloud project
* `flavor` - Savings Plan flavor (e.g. Rancher, C3-4, any instance flavor, ...)
* `period` - Periodicity of the Savings Plan
* `size` - Size of the Savings Plan
* `display_name` - Custom display name, used in invoices
* `auto_renewal` - Whether Savings Plan should be renewed at the end of the period
* `service_id` - Billing ID of the service
* `status` - Status of the Savings Plan
* `start_date` - Start date of the Savings Plan
* `end_date` - End date of the Savings Plan
* `period_end_action` - Action performed when reaching the end of the period (controles by the `auto_renewal` parameter)
* `period_start_date` - Start date of the current period
* `period_end_date` - End date of the current period

## Import

A savings plan can be imported using the following format: `service_name` and `id` of the savings plan, separated by "/" e.g.

```bash
$ terraform import ovh_savings_plan.plan service_name/savings_plan_id
```
