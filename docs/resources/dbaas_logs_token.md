---
subcategory : "Logs Data Platform"
---

# ovh_dbaas_logs_token

Allows to manipulate LDP tokens.

## Example Usage

```terraform
resource "ovh_dbaas_logs_token" "token" {
  service_name     = "ldp-xx-xxxxx"
  name             = "ExampleToken"
}
```

## Argument Reference

The following arguments are supported:

* `service_name` - (Required) The LDP service name
* `cluster_id` - Cluster ID. If not provided, the default cluster_id is used
* `name` - Name of the token

## Attributes Reference

* `service_name` - The LDP service name
* `cluster_id` - Cluster ID
* `name` - Name of the token
* `token_id` - ID of the token
* `value` - Token value
* `created_at` - Token creation date
* `updated_at` - Token last update date
