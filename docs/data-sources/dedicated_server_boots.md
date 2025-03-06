---
subcategory : "Dedicated Server"
---

# ovh_dedicated_server_boots (Data Source)

Use this data source to get the list of compatible netboots for a dedicated server associated with your OVHcloud Account.

## Example Usage

```terraform
data "ovh_dedicated_server_boots" "netboots" {
  service_name = "myserver"
  boot_type    = "harddisk"
}
```

## Argument Reference

* `service_name` - (Required) The internal name of your dedicated server.

* `boot_type` - (Optional) Filter the value of bootType property (harddisk, rescue, internal, network)

* `kernel` - (Optional) Filter the value of kernel property (iPXE script name)

## Attributes Reference

The following attributes are exported:

* `result` - The list of dedicated server netboots.
