---
subcategory : "Cloud Project"
---

# ovh_cloud_ssh_key (Data Source)

**This data source uses a Beta API.** Use this data source to retrieve information about a
single SSH key of a Public Cloud project (API v2), looked up by its name.

## Example Usage

```terraform
data "ovh_cloud_ssh_key" "my_key" {
  service_name = "aaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"
  name         = "my-deploy-key"
}
```

## Argument Reference

* `service_name` - (Optional) The id of the public cloud project. If omitted, the `OVH_CLOUD_PROJECT_SERVICE` environment variable is used.
* `name` - (Required) Name of the SSH key to retrieve. SSH key names are unique within a project.

## Attributes Reference

The following attributes are exported:

* `service_name` - The id of the public cloud project.
* `name` - SSH key name.
* `public_key` - SSH public key content.
* `created_at` - Creation date of the SSH key (RFC 3339 format).
* `updated_at` - Last update date of the SSH key (RFC 3339 format).
