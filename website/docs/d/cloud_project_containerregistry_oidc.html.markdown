---
subcategory : "Managed Private Registry"
---

# ovh_cloud_project_containerregistry_oidc (Data Source)

Use this data source to get a OVHcloud Managed Private Registry OIDC.

## Example Usage

```hcl
data "ovh_cloud_project_containerregistry_oidc" "my-oidc" {
  service_name = "XXXXXX"
  registry_id  = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxx"
}

output "oidc-client-id" {
  value = data.ovh_cloud_project_containerregistry_oidc.my-oidc.oidc_client_id
}
```

## Argument Reference

The following arguments are supported:

* `service_name` - (Optional) The id of the public cloud project. If omitted, the `OVH_CLOUD_PROJECT_SERVICE` environment variable is used.
* `registry_id` - The id of the Managed Private Registry.

## Attributes Reference

The following attributes are exported:

* `service_name` - The ID of the public cloud project. If omitted, the `OVH_CLOUD_PROJECT_SERVICE` environment variable is used.
* `registry_id` - The ID of the Managed Private Registry.
* `oidc_name` - The name of the OIDC provider.
* `oidc_endpoint` - The URL of an OIDC-compliant server.
* `oidc_client_id` - The client ID with which Harbor is registered as client application with the OIDC provider.
* `oidc_scope` - The scope sent to OIDC server during authentication. It's a comma-separated string that must contain 'openid' and usually also contains 'profile' and 'email'. To obtain refresh tokens it should also contain 'offline_access'.
* `oidc_groups_claim` - The name of Claim in the ID token whose value is the list of group names.
* `oidc_admin_group` - Specify an OIDC admin group name. All OIDC users in this group will have harbor admin privilege. Keep it blank if you do not want to.
* `oidc_verify_cert` - Set it to `false` if your OIDC server is hosted via self-signed certificate.
* `oidc_auto_onboard` - Skip the onboarding screen, so user cannot change its username. Username is provided from ID Token.
* `oidc_user_claim` - The name of the claim in the ID Token where the username is retrieved from. If not specified, it will default to 'name' (only useful when automatic Onboarding is enabled).
