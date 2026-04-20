---
subcategory : "Logs Data Platform"
---

# ovh_dbaas_logs_encryption_key

Creates a DBaaS Logs encryption key.

## Example Usage

```terraform
resource "ovh_dbaas_logs_encryption_key" "key" {
  service_name = "ldp-xx-xxxxx"
  title        = "my-encryption-key"
  content      = file("path/to/pgp-public-key.asc")
  fingerprint  = "ABCD1234EFGH5678IJKL"
}
```

## Argument Reference

The following arguments are supported:

* `service_name` - (Required, ForceNew) The LDP service name
* `title` - (Required) The encryption key title
* `content` - (Required, Sensitive, ForceNew) The PGP public key content
* `fingerprint` - (Required, ForceNew) The PGP key fingerprint

## Attributes Reference

Id is set to the encryption key Id. In addition, the following attributes are exported:

* `encryption_key_id` - The encryption key ID
* `created_at` - The encryption key creation date
* `is_editable` - Indicates if the key is editable

## Import

OVHcloud DBaaS Logs Encryption Key can be imported using the `service_name` and `encryption_key_id` of the key, separated by "/" E.g.,

```bash
$ terraform import ovh_dbaas_logs_encryption_key.key ldp-ra-XX/dc145bc2-eb01-4efe-a802-XXXXXX
```
