---
subcategory : "Load Balancer"
---

# ovh_cloud_loadbalancer_l7policy

Creates an L7 policy on a load balancer listener in a public cloud project. The policy's rules are managed inline with the policy.

## Example Usage

### Redirect to a Pool

```terraform
resource "ovh_cloud_loadbalancer_l7policy" "api" {
  service_name     = "xxxxxxxxxx"
  loadbalancer_id  = ovh_cloud_loadbalancer.lb.id
  listener_id      = ovh_cloud_loadbalancer_listener.http.id
  name             = "api-to-api-pool"
  action           = "REDIRECT_TO_POOL"
  redirect_pool_id = ovh_cloud_loadbalancer_pool.api.id
  position         = 1

  rules = [
    {
      type         = "PATH"
      compare_type = "STARTS_WITH"
      value        = "/api"
    },
    {
      type         = "HEADER"
      compare_type = "EQUAL_TO"
      key          = "X-Version"
      value        = "v2"
    },
  ]
}
```

### Redirect to a URL

```terraform
resource "ovh_cloud_loadbalancer_l7policy" "https_redirect" {
  service_name       = "xxxxxxxxxx"
  loadbalancer_id    = ovh_cloud_loadbalancer.lb.id
  listener_id        = ovh_cloud_loadbalancer_listener.http.id
  action             = "REDIRECT_TO_URL"
  redirect_url       = "https://www.example.com"
  redirect_http_code = 301

  rules = [
    {
      type         = "HOST_NAME"
      compare_type = "EQUAL_TO"
      value        = "example.com"
    },
  ]
}
```

## Argument Reference

The following arguments are supported:

* `service_name` - (Required) Service name of the resource representing the id of the cloud project. **Changing this value recreates the resource.**
* `loadbalancer_id` - (Required) ID of the load balancer. **Changing this value recreates the resource.**
* `listener_id` - (Required) ID of the listener. **Changing this value recreates the resource.**
* `action` - (Required) Action of the L7 policy (`REDIRECT_PREFIX`, `REDIRECT_TO_POOL`, `REDIRECT_TO_URL`, `REJECT`).
* `name` - (Optional) Name of the L7 policy.
* `description` - (Optional) Description of the L7 policy.
* `position` - (Optional) Position of the L7 policy in the listener's policy list. If omitted, the value assigned by the API is stored in the state.
* `redirect_prefix` - (Optional) Redirect prefix for `REDIRECT_PREFIX` action.
* `redirect_url` - (Optional) Redirect URL for `REDIRECT_TO_URL` action.
* `redirect_http_code` - (Optional) HTTP redirect code (`301`, `302`, `303`, `307`, `308`) for the `REDIRECT_PREFIX` and `REDIRECT_TO_URL` actions. If omitted, the value assigned by the API (`302`) is stored in the state.
* `redirect_pool_id` - (Optional) ID of the pool for `REDIRECT_TO_POOL` action.
* `rules` - (Optional) List of L7 rules for this policy. All rules must match for the policy to apply:
  * `type` - (Required) Type of the L7 rule (`COOKIE`, `FILE_TYPE`, `HEADER`, `HOST_NAME`, `PATH`).
  * `compare_type` - (Required) Comparison type (`CONTAINS`, `ENDS_WITH`, `EQUAL_TO`, `REGEX`, `STARTS_WITH`).
  * `value` - (Required) Value to compare against.
  * `key` - (Optional) Key for `COOKIE` and `HEADER` rule types.
  * `invert` - (Optional) Whether to invert the rule match. Defaults to the value assigned by the API.

## Attributes Reference

The following attributes are exported:

* `id` - L7 policy ID.
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
  * `rules` - Current state of the L7 rules:
    * `id` - Rule ID.
    * `type` - Rule type.
    * `compare_type` - Comparison type.
    * `value` - Value to compare against.
    * `key` - Key for `COOKIE` and `HEADER` rule types.
    * `invert` - Whether the rule match is inverted.
    * `operating_status` - Operating status of the rule.
    * `provisioning_status` - Provisioning status of the rule.

## Import

A cloud load balancer L7 policy can be imported using the `service_name`, `loadbalancer_id`, `listener_id`, and `l7policy_id`, separated by `/`:

```terraform
import {
  to = ovh_cloud_loadbalancer_l7policy.api
  id = "<service_name>/<loadbalancer_id>/<listener_id>/<l7policy_id>"
}
```

```bash
$ terraform import ovh_cloud_loadbalancer_l7policy.api service_name/loadbalancer_id/listener_id/l7policy_id
```
