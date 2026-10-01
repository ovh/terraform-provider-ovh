---
subcategory : "Logs Data Platform"
---

# ovh_dbaas_logs_encryption_key (Data Source)

Use this data source to retrieve information about a DBaaS Logs encryption key.

## Example Usage

### By title

```terraform
data "ovh_dbaas_logs_encryption_key" "key" {
  service_name = "ldp-xx-xxxxx"
  title        = "my-encryption-key"
}
```

### By encryption key ID

```terraform
data "ovh_dbaas_logs_encryption_key" "key" {
  service_name      = "ldp-xx-xxxxx"
  encryption_key_id = "dc145bc2-eb01-4efe-a802-XXXXXX"
}
```

## Argument Reference

* `service_name` - (Required) The LDP service name
* `title` - (Optional) The encryption key title. Conflicts with `encryption_key_id`.
* `encryption_key_id` - (Optional) The encryption key ID. Conflicts with `title`.

At least one of `title` or `encryption_key_id` must be specified.

## Attributes Reference

* `fingerprint` - The PGP key fingerprint
* `created_at` - The encryption key creation date
* `is_editable` - Indicates if the key is editable
