---
subcategory : "Load Balancer"
---

# ovh_cloud_loadbalancer_l7policies (Data Source)

Use this data source to list the L7 policies of a load balancer listener in a public cloud project.

## Example Usage

```terraform
data "ovh_cloud_loadbalancer_l7policies" "policies" {
  service_name    = "<public cloud project ID>"
  loadbalancer_id = "<load balancer ID>"
  listener_id     = "<listener ID>"
}

output "l7policy_actions" {
  value = [for p in data.ovh_cloud_loadbalancer_l7policies.policies.l7policies : p.action]
}
```

## Argument Reference

The following arguments are supported:

* `service_name` - (Required) Service name of the resource representing the id of the cloud project.
* `loadbalancer_id` - (Required) ID of the load balancer.
* `listener_id` - (Required) ID of the listener.

## Attributes Reference

The following attributes are exported:

* `l7policies` - List of L7 policies. Each element exports:
  * `id` - L7 policy ID.
  * `name` - Name of the L7 policy.
  * `description` - Description of the L7 policy.
  * `action` - Action of the L7 policy (`REDIRECT_PREFIX`, `REDIRECT_TO_POOL`, `REDIRECT_TO_URL`, `REJECT`).
  * `position` - Position of the L7 policy.
  * `redirect_prefix` - Redirect prefix for `REDIRECT_PREFIX` action.
  * `redirect_url` - Redirect URL for `REDIRECT_TO_URL` action.
  * `redirect_http_code` - HTTP redirect code (`301`, `302`, `303`, `307`, `308`).
  * `redirect_pool_id` - ID of the pool for `REDIRECT_TO_POOL` action.
  * `rules` - List of L7 rules for this policy:
    * `type` - Type of the L7 rule (`COOKIE`, `FILE_TYPE`, `HEADER`, `HOST_NAME`, `PATH`).
    * `compare_type` - Comparison type (`CONTAINS`, `ENDS_WITH`, `EQUAL_TO`, `REGEX`, `STARTS_WITH`).
    * `value` - Value to compare against.
    * `key` - Key for `COOKIE` and `HEADER` rule types.
    * `invert` - Whether to invert the rule match.
  * `checksum` - Computed hash representing the current target specification value.
  * `created_at` - Creation date of the L7 policy.
  * `updated_at` - Last update date of the L7 policy.
  * `resource_status` - L7 policy readiness in the system (`CREATING`, `DELETING`, `ERROR`, `OUT_OF_SYNC`, `READY`, `UPDATING`).
  * `current_state` - Current state of the L7 policy:
    * `name` - L7 policy name.
    * `description` - L7 policy description.
    * `action` - L7 policy action.
    * `position` - L7 policy position.
    * `redirect_prefix` - Redirect prefix.
    * `redirect_url` - Redirect URL.
    * `redirect_http_code` - HTTP redirect code.
    * `redirect_pool_id` - Redirect pool ID.
    * `operating_status` - Operating status of the L7 policy.
    * `provisioning_status` - Provisioning status of the L7 policy.
    * `rules` - Current state of the L7 rules (same schema as `rules`, plus `id`, `operating_status` and `provisioning_status`).
