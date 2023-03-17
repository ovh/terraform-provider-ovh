---
subcategory : "Dedicated Server"
---

# ovh_me_ipxe_script (Data Source)

Use this data source to retrieve information about an IPXE Script.

## Example Usage

```hcl
data "ovh_me_ipxe_script" "script" {
  name = "myscript"
}
```

## Argument Reference

* `name` - (Required) The name of the IPXE Script.

## Attributes Reference

* `name` - See Argument Reference above.
* `script` - The content of the script.
