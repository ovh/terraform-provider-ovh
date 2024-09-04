---
subcategory : "Domain names"
---

# ovh_domain_zone_import

Handle a whole DNS zone using a zone file.

~> __WARNING__ This resource and resource `ovh_domain_zone_record` should not be used together as `ovh_domain_zone_import` controls the whole DNS zone at once.

## Example Usage

```hcl
resource "ovh_domain_zone_import" "import" {
  zone_name = "mysite.ovh"
  zone_file = file("./example.zone")
}
```

## Argument Reference

* `zone_name` - (Required) The name of the domain zone
* `zone_file` - (Required) Content of the zone file to import

## Attributes Reference

* `exported_content` - Zone file exported from the API