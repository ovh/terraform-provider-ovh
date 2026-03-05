---
subcategory : "Floating IP"
---

# ovh_cloud_floating_ip

Creates a floating IP in a public cloud project.

## Example Usage

```terraform
resource "ovh_cloud_floating_ip" "fip" {
  service_name = "xxxxxxxxxx"
  region       = "GRA1"
  description  = "my-floating-ip"
}
```

## Argument Reference

The following arguments are supported:

* `service_name` - (Required) Service name of the resource representing the id of the cloud project. **Changing this value recreates the resource.**
* `region` - (Required) Region where the floating IP will be created. **Changing this value recreates the resource.**
* `availability_zone` - (Optional) Availability zone of the floating IP. **Changing this value recreates the resource.**
* `description` - (Optional) Description of the floating IP.

## Attributes Reference

The following attributes are exported:

* `id` - Floating IP ID.
* `checksum` - Computed hash representing the current target specification value.
* `created_at` - Creation date of the floating IP.
* `updated_at` - Last update date of the floating IP.
* `resource_status` - Floating IP readiness in the system (`CREATING`, `DELETING`, `ERROR`, `OUT_OF_SYNC`, `READY`, `UPDATING`).
* `current_state` - Current state of the floating IP:
  * `ip` - Floating IP address.
  * `status` - Floating IP status.
  * `network_id` - Network ID associated with the floating IP.
  * `description` - Description of the floating IP.
  * `region` - Region of the floating IP.
  * `availability_zone` - Availability zone of the floating IP.

## Import

A cloud floating IP can be imported using the `service_name` and `floating_ip_id`, separated by `/`:

```terraform
import {
  to = ovh_cloud_floating_ip.fip
  id = "<service_name>/<floating_ip_id>"
}
```

```bash
$ terraform import ovh_cloud_floating_ip.fip service_name/floating_ip_id
```
