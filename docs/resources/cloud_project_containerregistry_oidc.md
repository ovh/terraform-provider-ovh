---
subcategory : "Managed Private Registry (MPR)"
---

# ovh_cloud_project_containerregistry_oidc

Creates an OIDC configuration in an OVHcloud Managed Private Registry.

## Example Usage

```terraform
resource "ovh_cloud_project_containerregistry_oidc" "my_oidc" {
  service_name = "XXXXXX"
  registry_id  = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxx"

  #required field
  oidc_name          = "my-oidc-provider"
  oidc_endpoint      = "https://xxxx.yyy.com"
  oidc_client_id     = "xxx"
  oidc_client_secret = "xxx"
  oidc_scope         = "openid,profile,email,offline_access"

  #optional field
  oidc_group_filter = "harbor-admin"
  oidc_groups_claim = "groups"
  oidc_admin_group  = "harbor-admin"
  oidc_verify_cert  = true
  oidc_auto_onboard = true
  oidc_user_claim   = "preferred_username"
  delete_users      = false
}

output "oidc_client_secret" {
  value = ovh_cloud_project_containerregistry_oidc.my_oidc.oidc_client_secret
  sensitive = true
}
```

## Argument Reference

The following arguments are supported:

* `service_name` - The ID of the public cloud project. If omitted, the `OVH_CLOUD_PROJECT_SERVICE` environment variable is used. **Changing this value recreates the resource.**
* `registry_id` - The ID of the Managed Private Registry. **Changing this value recreates the resource.**
* `oidc_name` - The name of the OIDC provider.
* `oidc_endpoint` - The URL of an OIDC-compliant server.
* `oidc_client_id` - The client ID with which Harbor is registered as client application with the OIDC provider.
* `oidc_client_secret` - The secret for the Harbor client application.
* `oidc_scope` - The scope sent to OIDC server during authentication. It's a comma-separated string that must contain 'openid' and usually also contains 'profile' and 'email'. To obtain refresh tokens it should also contain 'offline_access'.
* `oidc_group_filter` - The regular expression to select matching groups from the Group Claim Name list. Matching groups are added to Harbor. This filter does not limit the users’ capability to log in into Harbor.
* `oidc_groups_claim` - The name of Claim in the ID token whose value is the list of group names.
* `oidc_admin_group` - Specify an OIDC admin group name. All OIDC users in this group will have harbor admin privilege. Keep it blank if you do not want to.
* `oidc_verify_cert` - Set it to `false` if your OIDC server is hosted via self-signed certificate.
* `oidc_auto_onboard` - Skip the onboarding screen, so user cannot change its username. Username is provided from ID Token.
* `oidc_user_claim` - The name of the claim in the ID Token where the username is retrieved from. If not specified, it will default to 'name' (only useful when automatic Onboarding is enabled).
* `delete_users` - Delete existing users from Harbor. OIDC can't be enabled if there is at least one user already created. This parameter is only used at OIDC configuration creation. **Changing this value recreates the resource.**

## Attributes Reference

The following attributes are exported:

* `service_name` - See Argument Reference above.
* `registry_id` - See Argument Reference above.
* `oidc_name` - See Argument Reference above.
* `oidc_endpoint` - See Argument Reference above.
* `oidc_client_id` - See Argument Reference above.
* `oidc_client_secret` - (Sensitive) See Argument Reference above.
* `oidc_scope` - See Argument Reference above.
* `oidc_group_filter` - See Argument Reference above.
* `oidc_groups_claim` - See Argument Reference above.
* `oidc_admin_group` - See Argument Reference above.
* `oidc_verify_cert` - See Argument Reference above.
* `oidc_auto_onboard` - See Argument Reference above.
* `oidc_user_claim` - See Argument Reference above.

## Timeouts

```terraform
resource "ovh_cloud_project_containerregistry_oidc" "my_oidc" {
  # ...

  timeouts {
    create = "1h"
    update = "45m"
    delete = "50s"
  }
}
```
* `create` - (Default 10m)
* `update` - (Default 10m)
* `delete` - (Default 10m)

## Import

OVHcloud Managed Private Registry OIDC can be imported using the tenant `service_name` and registry id `registry_id` separated by "/" E.g.,

```bash
$ terraform import ovh_cloud_project_containerregistry_oidc.my-oidc service_name/registry_id
```
