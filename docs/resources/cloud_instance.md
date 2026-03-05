---
subcategory : "Cloud Instances"
---

# ovh_cloud_instance

Creates a cloud instance in a public cloud project.

## Example Usage

```terraform
resource "ovh_cloud_instance" "instance" {
  service_name = "xxxxxxxxxx"
  name         = "my-instance"
  flavor_id    = "068a57cf-b1b4-428f-9b17-4f32a526390c"
  image_id     = "8d75e170-1ef9-4e25-8fc8-d231929e56e8"
  region       = "GRA1"
}
```

### With network

```terraform
resource "ovh_cloud_instance" "instance" {
  service_name = "xxxxxxxxxx"
  name         = "my-instance"
  flavor_id    = "068a57cf-b1b4-428f-9b17-4f32a526390c"
  image_id     = "8d75e170-1ef9-4e25-8fc8-d231929e56e8"
  region       = "GRA1"

  networks {
    id = "fbdf6240-8b56-4626-b57c-25e4af487606"
  }
}
```

## Argument Reference

The following arguments are supported:

* `service_name` - (Required) Service name of the resource representing the id of the cloud project. **Changing this value recreates the resource.**
* `name` - (Required) Instance name.
* `flavor_id` - (Required) Flavor ID for the instance.
* `image_id` - (Required) Image ID for the instance.
* `region` - (Required) Region where the instance will be created. **Changing this value recreates the resource.**
* `availability_zone` - (Optional) The availability zone where the instance will be created. **Changing this value recreates the resource.**
* `networks` - (Optional) List of networks to attach to the instance.
  * `id` - (Optional) Network ID to attach.
  * `public` - (Optional) Whether to attach the public network.
  * `subnet_id` - (Optional) Subnet ID for the private network.
  * `floating_ip_id` - (Optional) Floating IP ID for the network.
* `volume_ids` - (Optional) List of volume IDs to attach to the instance.
* `ssh_key_name` - (Optional) SSH key name to associate with the instance.

## Attributes Reference

The following attributes are exported:

* `id` - Instance ID.
* `checksum` - Computed hash representing the current target specification value.
* `created_at` - Creation date of the instance.
* `updated_at` - Last update date of the instance.
* `resource_status` - Instance readiness in the system (`CREATING`, `DELETING`, `ERROR`, `READY`, `SUSPENDED`, `UPDATING`).
* `current_state` - Current state of the instance:
  * `flavor` - Flavor details:
    * `id` - Flavor identifier.
    * `name` - Flavor name.
    * `vcpus` - Number of vCPUs.
    * `ram` - RAM in MB.
    * `disk` - Local disk size in GB.
  * `image` - Image details:
    * `id` - Image identifier.
    * `name` - Image name.
    * `status` - Image status.
  * `name` - Instance name.
  * `host_id` - Host identifier.
  * `ssh_key_name` - Associated SSH key name.
  * `project_id` - Project identifier.
  * `user_id` - User identifier.
  * `networks` - List of instance networks:
    * `id` - Network identifier.
    * `public` - Whether the network is public.
    * `subnet_id` - Associated subnet identifier.
    * `gateway_id` - Gateway identifier.
    * `floating_ip_id` - Floating IP identifier.
    * `addresses` - IP addresses:
      * `ip` - IP address.
      * `mac` - MAC address.
      * `type` - Address type (fixed, floating, ...).
      * `version` - IP version.
  * `volumes` - Attached block volumes:
    * `id` - Volume identifier.
    * `name` - Volume name.
    * `size` - Volume size in GB.
  * `security_groups` - Security groups attached to the instance.

## Import

A cloud instance can be imported using the `service_name` and `id`, separated by `/`:

```terraform
import {
  to = ovh_cloud_instance.instance
  id = "<service_name>/<instance_id>"
}
```

```bash
$ terraform plan -generate-config-out=instance.tf
$ terraform apply
```
