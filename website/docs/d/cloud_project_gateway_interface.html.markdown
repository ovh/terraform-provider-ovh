---
subcategory: "Public Cloud Network"
---

# ovh_cloud_project_gateway_interface

Use this datasource to get a public cloud project Gateway Interface.

## Example Usage

```hcl
data "ovh_cloud_project_gateway_interface" "interface" {
	service_name = "XXXXXX"
	region       = "GRA11"
	id           = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxx"
	interface_id = "yyyyyyyy-yyyy-yyyy-yyyy-yyyyyyyy"
}
```

## Argument Reference

The following arguments are supported:

* `service_name` - (Required) ID of the cloud project
* `region` - (Required) Region of the gateway
* `id` - (Required) ID of the gateway
* `interface_id` - (Required) ID of the interface

## Attributes Reference

The following attributes are exported:

* `service_name` - ID of the cloud project
* `region` - Region of the gateway
* `id` - ID of the gateway
* `subnet_id` - ID of the subnet to add
* `interface_id` - ID of the interface
* `ip` - IP of the interface
* `network_id` - Network ID of the interface
