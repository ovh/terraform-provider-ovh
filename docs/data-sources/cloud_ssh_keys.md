---
subcategory : "Cloud Project"
---

# ovh_cloud_ssh_keys (Data Source)

**This data source uses a Beta API.** Use this data source to list the SSH keys of a Public
Cloud project (API v2).

## Example Usage

```terraform
data "ovh_cloud_ssh_keys" "my_keys" {
  service_name = "aaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"
}
```

## Argument Reference

* `service_name` - (Optional) The id of the public cloud project. If omitted, the `OVH_CLOUD_PROJECT_SERVICE` environment variable is used.

## Attributes Reference

The following attributes are exported:

* `service_name` - The id of the public cloud project.
* `ssh_keys` - The list of SSH keys of the project. Each element contains:
  * `name` - SSH key name.
  * `public_key` - SSH public key content.
  * `created_at` - Creation date of the SSH key (RFC 3339 format).
  * `updated_at` - Last update date of the SSH key (RFC 3339 format).
