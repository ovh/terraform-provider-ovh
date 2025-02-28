---
subcategory : "Domain names"
---

# ovh_domain_ds_records

Use this resource to manage a domain's DS records.

## Example Usage

```hcl
resource "ovh_domain_ds_records" "ds_records" {
  domain = "mydomain.ovh"
  
  ds_records {
      algorithm = "RSASHA1_NSEC3_SHA1"
      flags = "KEY_SIGNING_KEY"
      public_key = "my_base64_encoded_public_key"
      tag = 12345
  }
}
```

## Argument Reference

The following arguments are supported:

* `domain` - (Required) Domain name for which to manage DS records
* `ds_records` - (Required) Details about a DS record
  * `algorithm` - (Required) The record algorithm (`RSASHA1`, `RSASHA1_NSEC3_SHA1`, `RSASHA256`, `RSASHA512`, `ECDSAP256SHA256`, `ECDSAP384SHA384`, `ED25519`)
  * `flags` - (Required) The record flag (`ZONE_SIGNING_KEY`, `KEY_SIGNING_KEY`)
  * `public_key` - (Required) The record base64 encoded public key
  * `tag` - (Required) The record tag


## Attributes Reference

* `domain` - Domain name and resource ID
* `ds_records` - Details about a DS record
  * `algorithm` - The record algorithm (`RSASHA1`, `RSASHA1_NSEC3_SHA1`, `RSASHA256`, `RSASHA512`, `ECDSAP256SHA256`, `ECDSAP384SHA384`, `ED25519`)
  * `flags` - The record flag (`ZONE_SIGNING_KEY`, `KEY_SIGNING_KEY`)
  * `public_key` - The record base64 encoded public key
  * `tag` - The record tag

## Import

DS records can be imported using their `domain`.

Using the following configuration:

```hcl
import {
  to = ovh_domain_ds_records.ds_records
  id = "<domain name>"
}
```

You can then run:

```bash
$ terraform plan -generate-config-out=ds_records.tf
$ terraform apply
```

The file `ds_records.tf` will then contain the imported resource's configuration, that can be copied next to the `import` block above.
See https://developer.hashicorp.com/terraform/language/import/generating-configuration for more details.
