---
subcategory : "Managed Kubernetes Service (MKS)"
---

# ovh_cloud_project_kube_oidc (Data Source)

Use this data source to get a OVHcloud Managed Kubernetes Service cluster OIDC.

## Example Usage

```hcl
data "ovh_cloud_project_kube_oidc" "oidc" {
  service_name = "XXXXXX"
  kube_id      = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxx"
}

output "oidc-val" {
  value = data.ovh_cloud_project_kube_oidc.oidc.client_id
}
```

## Argument Reference

The following arguments are supported:

* `service_name` - (Optional) The id of the public cloud project. If omitted,
  the `OVH_CLOUD_PROJECT_SERVICE` environment variable is used.

* `kube_id` - The id of the managed kubernetes cluster.

## Attributes Reference

The following attributes are exported:

* `service_name` - See Argument Reference above.
* `kube_id` - See Argument Reference above.

* `client_id` - The OIDC client ID.

* `issuer_url` - The OIDC issuer url.

* `oidc_username_claim` - JWT claim to use as the user name. By default sub, which is expected to be a unique identifier of the end user. Admins can choose other claims, such as email or name, depending on their provider. However, claims other than email will be prefixed with the issuer URL to prevent naming clashes with other plugins.

* `oidc_username_prefix` - Prefix prepended to username claims to prevent clashes with existing names (such as system: users). For example, the value oidc: will create usernames like oidc:jane.doe. If this field isn't set and `oidc_username_claim` is a value other than email the prefix defaults to ( Issuer URL )# where ( Issuer URL ) is the value of oidcIssuerUrl. The value - can be used to disable all prefixing.

* `oidc_groups_claim` - Array of JWT claim to use as the user's group. If the claim is present it must be an array of strings.

* `oidc_groups_prefix` - Prefix prepended to group claims to prevent clashes with existing names (such as system: groups). For example, the value oidc: will create group names like oidc:engineering and oidc:infra.

* `oidc_required_claim` - Array of key=value pairs that describe required claims in the ID Token. If set, the claims are verified to be present in the ID Token with a matching value."

* `oidc_signing_algs` - Array of signing algorithms accepted. Default is \"RS256\".

* `oidc_ca_content` - Content of the certificate for the CA, in base64 format, that signed your identity provider's web certificate. Defaults to the host's root CAs.
