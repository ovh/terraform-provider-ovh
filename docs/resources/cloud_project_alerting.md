---
subcategory : "Cloud Project"
---

# ovh_cloud_project_alerting

Creates an alert on a public cloud project.

## Example Usage

```terraform
resource "ovh_cloud_project_alerting" "my_alert" {
  service_name = "XXX"
  delay = 3600
  email = "aaa.bbb@domain.com"
  monthly_threshold = 1000
}
```

## Argument Reference

* `service_name` - The id of the public cloud project. If omitted, the `OVH_CLOUD_PROJECT_SERVICE` environment variable is used.

* `delay` - Delay between two alerts in seconds
* `email` - Email to contact
* `monthly_threshold` - Monthly threshold for this alerting in currency

## Attributes Reference

The following attributes are exported:

* `id` - Alert ID
* `creationDate` - Alerting creation date
* `delay` - Delay between two alerts in seconds
* `email` - Email to contact
* `monthly_threshold` - Monthly threshold for this alerting in currency
* `formatted_monthly_threshold` - Formatted monthly threshold for this alerting
  * `currency_code` - Currency of the monthly threshold
  * `text`: Text representation of the monthly threshold
  * `value`: Value of the monthly threshold
