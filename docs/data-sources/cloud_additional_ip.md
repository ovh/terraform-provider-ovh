---
subcategory : "Public IPs"
---

# ovh_cloud_additional_ip (Data Source)

Use this data source to retrieve information about an additional IP in a public cloud project. Additional IPs are created and deleted through the lifecycle of other products (typically instances they are attached to), so this data source is a read-only view.

## Example Usage

```terraform
data "ovh_cloud_additional_ip" "ip" {
  service_name = "<public cloud project ID>"
  id           = "<additional IP address>"
}

output "additional_ip_status" {
  value = data.ovh_cloud_additional_ip.ip.resource_status
}
```

## Argument Reference

The following arguments are supported:

* `service_name` - (Required) Service name of the resource representing the id of the cloud project. If omitted, the `OVH_CLOUD_PROJECT_SERVICE` environment variable is used.
* `id` - (Required) IP address of the additional IP.

## Attributes Reference

The following attributes are exported:

* `checksum` - Checksum field of the API envelope. Always empty for this read-only IP type.
* `resource_status` - Additional IP readiness in the system (`CREATING`, `DELETING`, `ERROR`, `OUT_OF_SYNC`, `READY`, `UNKNOWN`, `UPDATING`).
* `current_state` - Current state of the additional IP:
  * `id` - Identifier of the additional IP.
  * `ip` - IP address of the additional IP.
  * `ip_block` - IP block the additional IP belongs to. May be null.
  * `associated_resource` - Resource the additional IP is currently attached to. Null when the IP is not attached to any resource:
    * `id` - ID of the associated resource.
    * `type` - Type of the associated resource.
  * `location` - Location details:
    * `region` - Region.
    * `availability_zone` - Availability zone.
* `current_tasks` - Ongoing asynchronous tasks related to the additional IP. Each element exports:
  * `id` - Identifier of the current task.
  * `link` - Link to the task details.
  * `status` - Current global status of the current task.
  * `type` - Type of the current task.
  * `errors` - Errors that occurred on the task:
    * `message` - Error description.
