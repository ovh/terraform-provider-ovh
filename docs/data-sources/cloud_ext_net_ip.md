---
subcategory : "Public IPs"
---

# ovh_cloud_ext_net_ip (Data Source)

Use this data source to retrieve information about an external network IP in a public cloud project. External network IPs are created and deleted through the lifecycle of other products (typically instances attached to the public network), so this data source is a read-only view.

## Example Usage

```terraform
data "ovh_cloud_ext_net_ip" "ip" {
  service_name = "<public cloud project ID>"
  id           = "<external network IP address>"
}

output "ext_net_ip_status" {
  value = data.ovh_cloud_ext_net_ip.ip.resource_status
}
```

## Argument Reference

The following arguments are supported:

* `service_name` - (Required) Service name of the resource representing the id of the cloud project. If omitted, the `OVH_CLOUD_PROJECT_SERVICE` environment variable is used.
* `id` - (Required) IP address of the external network IP.

## Attributes Reference

The following attributes are exported:

* `checksum` - Checksum field of the API envelope. Always empty for this read-only IP type.
* `created_at` - Creation date of the external network IP.
* `updated_at` - Last update date of the external network IP.
* `resource_status` - External network IP readiness in the system (`CREATING`, `DELETING`, `ERROR`, `OUT_OF_SYNC`, `READY`, `UPDATING`).
* `current_state` - Current state of the external network IP:
  * `id` - Identifier of the external network IP.
  * `ip` - IP address of the external network IP.
  * `associated_resource` - Resource the external network IP is currently attached to. Null when the IP is not attached to any resource:
    * `id` - ID of the associated resource.
    * `type` - Type of the associated resource.
  * `location` - Location details:
    * `region` - Region.
    * `availability_zone` - Availability zone.
* `current_tasks` - Ongoing asynchronous tasks related to the external network IP. Each element exports:
  * `id` - Identifier of the current task.
  * `link` - Link to the task details.
  * `status` - Current global status of the current task.
  * `type` - Type of the current task.
  * `errors` - Errors that occurred on the task:
    * `message` - Error description.
