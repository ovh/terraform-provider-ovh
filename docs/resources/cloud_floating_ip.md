---
subcategory : "Public IPs"
---

# ovh_cloud_floating_ip

Creates a floating IP in a public cloud project.

## Example Usage

```terraform
resource "ovh_cloud_floating_ip" "ip" {
  service_name = "<public cloud project ID>"
  region       = "GRA1"
  description  = "my-floating-ip"
}
```

## Argument Reference

The following arguments are supported:

* `service_name` - (Required) Service name of the resource representing the id of the cloud project. If omitted, the `OVH_CLOUD_PROJECT_SERVICE` environment variable is used. **Changing this value recreates the resource.**
* `region` - (Required) Region where the floating IP will be created. **Changing this value recreates the resource.**
* `availability_zone` - (Optional) Availability zone for the floating IP. **Changing this value recreates the resource.**
* `description` - (Optional) Description of the floating IP. This is the only argument that can be updated in place.

## Attributes Reference

The following attributes are exported:

* `id` - IP address of the floating IP.
* `checksum` - Computed hash representing the current target specification value.
* `created_at` - Creation date of the floating IP.
* `updated_at` - Last update date of the floating IP.
* `resource_status` - Floating IP readiness in the system (`CREATING`, `DELETING`, `ERROR`, `OUT_OF_SYNC`, `READY`, `UPDATING`).
* `current_state` - Current state of the floating IP:
  * `id` - OpenStack identifier of the floating IP.
  * `ip` - IP address of the floating IP.
  * `status` - OpenStack status of the floating IP (`ACTIVE`, `DOWN`, `ERROR`).
  * `description` - Description of the floating IP.
  * `network` - External network the floating IP belongs to:
    * `id` - Network ID.
  * `associated_resource` - Resource the floating IP is currently attached to. Null when the floating IP is not attached to any resource:
    * `id` - ID of the associated resource.
    * `type` - Type of the associated resource.
  * `location` - Location details:
    * `region` - Region.
    * `availability_zone` - Availability zone.
* `current_tasks` - Ongoing asynchronous tasks related to the floating IP. Each element exports:
  * `id` - Identifier of the current task.
  * `link` - Link to the task details.
  * `status` - Current global status of the current task.
  * `type` - Type of the current task.
  * `errors` - Errors that occurred on the task:
    * `message` - Error description.

## Import

A floating IP can be imported using the `service_name` and the IP address, separated by `/`:

```terraform
import {
  to = ovh_cloud_floating_ip.ip
  id = "<service_name>/<ip>"
}
```

```bash
$ terraform import ovh_cloud_floating_ip.ip service_name/ip
```
