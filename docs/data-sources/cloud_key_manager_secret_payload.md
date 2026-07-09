---
subcategory : "Key Manager"
---

# ovh_cloud_key_manager_secret_payload (Data Source)

Use this data source to retrieve the payload (secret material) of a Barbican Key Manager secret.

~> **WARNING:** The `payload` attribute contains the secret material and is marked as sensitive. It will be stored in clear text in the Terraform state. Read more about [sensitive data in state](https://www.terraform.io/language/state/sensitive-data).

## Example Usage

```terraform
data "ovh_cloud_key_manager_secret_payload" "payload" {
  service_name = "Public cloud project ID"
  secret_id    = "00000000-0000-0000-0000-000000000000"
}

output "secret_payload" {
  value     = data.ovh_cloud_key_manager_secret_payload.payload.payload
  sensitive = true
}
```

## Argument Reference

* `service_name` - (Required) The id of the public cloud project. If omitted, the `OVH_CLOUD_PROJECT_SERVICE` environment variable is used.
* `secret_id` - (Required) UUID of the secret.

## Attributes Reference

The following attributes are exported:

* `payload` - The payload (secret material) of the secret. This value is sensitive.
