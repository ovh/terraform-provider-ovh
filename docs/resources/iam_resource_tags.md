---
subcategory : "Account Management (IAM)"
---

# ovh_iam_resource_tags

Manages tags for OVHcloud IAM resources.

This resource allows you to apply and manage tags on OVHcloud resources identified by their URN (Uniform Resource Name). Tags are key-value pairs that help organize and categorize your resources.

## Example Usage

### Basic Usage

```terraform
resource "ovh_iam_resource_tags" "project_tags" {
  resource_urn = "urn:v1:eu:resource:cloudProject:1234567890abcdef"
  
  tags = {
    environment    = "staging"
    cost_center    = "engineering"
    project        = "web-app"
    owner          = "team@example.com"
    backup_policy  = "daily"
    compliance     = "gdpr"
  }
}
```

### Using with Data Source

```terraform
data "ovh_cloud_project" "my_project" {
  service_name = "cef572b51c534a8f9fee73ac81957fbd"
}

resource "ovh_iam_resource_tags" "project_tags" {
  resource_urn = data.ovh_cloud_project.my_project.iam.urn

  tags = {
    environment = "production"
    team        = "platform"
    managed_by  = "terraform"
  }
}
```

## Argument Reference

The following arguments are supported:

* `resource_urn` - (Required, ForceNew) The URN (Uniform Resource Name) of the resource to manage tags for. Changing this forces a new resource to be created. The URN format is typically `urn:v1:{region}:resource:{resourceType}:{resourceId}`.

* `tags` - (Optional) A map of tags to apply to the resource. Each tag consists of a key-value pair. Tag keys must match the pattern `^[a-zA-Z0-9_.:/=+@-]{1,128}$` (1-128 characters) and values must match `^[a-zA-Z0-9_.:/=+@-]{0,256}$` (0-256 characters). Both keys and values can contain alphanumeric characters, underscores, dots, colons, slashes, equals signs, plus signs, at signs, and hyphens. **Note:** Tags with keys prefixed by `ovh:` are managed by OVH and cannot be set via the API.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the resource, which is the same as the `urn`.

## Import

IAM resource tags can be imported using the resource URN:

```bash
terraform import ovh_iam_resource_tags.my_tags "urn:v1:eu:resource:cloudProject:1234567890abcdef"
```

After importing, you should update your Terraform configuration to match the imported state. The import will bring in all tags currently applied to the resource, but only tags defined in your configuration will be managed by Terraform.

## Notes

* This resource only manages tags that are explicitly defined in the `tags` argument. Other tags on the resource that are not managed by this Terraform resource will not be affected.
* When the resource is destroyed, only the tags defined in the configuration will be removed from the OVHcloud resource. The resource itself will not be deleted.
* Tag keys and values are case-sensitive.
* Empty tag values are allowed (minimum length is 0 characters for values).
* If you need to remove all tags, set `tags = {}` in your configuration.
* Tags prefixed with `ovh:` are system-managed and cannot be modified through this resource. They will be visible when reading the resource but attempting to set them will result in a validation error.
* The resource uses the OVH API v2 `/iam/resource/{urn}` endpoint with PUT operations to update all tags in a single operation, preserving any unmanaged tags (including `ovh:` prefixed tags).
