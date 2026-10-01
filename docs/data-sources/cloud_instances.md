---
subcategory: "Instances"
---

# ovh_cloud_instances (Data Source)

Use this data source to list the instances of a public cloud project.

~> **NOTE** `current_state.shares` is not populated by the list endpoint. Use the [`ovh_cloud_instance`](cloud_instance.md) data source to retrieve the shares attached to a specific instance.

## Example Usage

```terraform
data "ovh_cloud_instances" "instances" {
  service_name = "<Public cloud project id>"
}

# Names of the instances that are up and running in GRA11.
output "gra11_running_instance_names" {
  value = sort([
    for instance in data.ovh_cloud_instances.instances.instances :
    instance.name
    if instance.region == "GRA11" && instance.current_state.power_state == "ACTIVE"
  ])
}

# Instance id keyed by name, ready to be fed to another resource.
output "instance_ids_by_name" {
  value = {
    for instance in data.ovh_cloud_instances.instances.instances :
    instance.name => instance.id
  }
}
```

## Argument Reference

The following arguments are supported:

* `service_name` - (Required) Service name of the resource representing the id of the cloud project.

## Attributes Reference

The following attributes are exported:

* `instances` - Instances of the project:
  * `id` - Unique identifier of the instance.
  * `name` - Instance name.
  * `region` - Code of the region where the instance is provisioned (for example GRA11, BHS5).
  * `availability_zone` - Availability zone within the region where the instance is placed, null in regions that have none.
  * `flavor_id` - Unique identifier of the flavor.
  * `image_id` - Identifier of the image the instance boots from, null for a boot-from-volume instance.
  * `power_state` - Requested power state of the instance (`ACTIVE`, `SHUTOFF` or `SHELVED`).
  * `checksum` - Computed hash representing the current target specification value. It implements optimistic concurrency control: the value is echoed back on update and the request is rejected when it no longer matches server-side.
  * `created_at` - Creation date of the instance, as an RFC 3339 timestamp.
  * `updated_at` - Last modification date of the instance, as an RFC 3339 timestamp.
  * `resource_status` - Instance readiness in the system (`CREATING`, `DELETING`, `ERROR`, `OUT_OF_SYNC`, `READY`, `UPDATING`). Distinct from `current_state.power_state`, which carries the lower-level OpenStack administrative power state.
  * `current_state` - Observed state of the instance as reported by the compute backend, as opposed to the requested specification exposed at root level. Null while the instance is still being created and no backend state is available yet:
    * `name` - Observed display name of the instance.
    * `power_state` - Observed administrative power state of the instance as reported by OpenStack. It may transiently differ from the requested `power_state` while a power transition is in progress.
    * `locked` - Whether the instance is locked against modifications. While locked, mutating actions are refused until it is unlocked.
    * `ssh_key_name` - Name of the SSH key pair injected into the instance at boot, null when none was provided.
    * `host_id` - Opaque identifier of the physical host the instance is running on, as exposed by OpenStack. Null when not available.
    * `project_id` - Identifier of the Public Cloud project the instance belongs to.
    * `user_id` - Identifier of the OpenStack user that owns the instance.
    * `flavor` - Observed flavor of the instance, with its full sizing details:
      * `id` - Unique identifier of the flavor.
      * `name` - Human-readable flavor name (the commercial flavor label).
      * `vcpus` - Number of virtual CPUs provided by the flavor.
      * `ram` - Amount of RAM provided by the flavor, in MB.
      * `disk` - Size of the flavor's local root disk, in GB.
      * `swap` - Size of the flavor's swap space, in MB.
      * `ephemeral` - Size of the flavor's ephemeral disk, in GB.
    * `image` - Observed image the instance was booted from, null for a boot-from-volume instance which has no image:
      * `id` - Unique identifier of the image.
      * `name` - Human-readable image name, null when the backend does not report it.
      * `size` - Size of the image, in bytes. Null when the backend does not report it.
      * `status` - Lifecycle status of the image as reported by Glance.
      * `deprecated` - Whether the image is flagged as deprecated. A deprecated image still boots existing instances but is no longer recommended for new ones.
    * `location` - Observed region and availability zone where the instance is provisioned:
      * `region` - Code of the region where the instance is provisioned (for example GRA11, BHS5).
      * `availability_zone` - Availability zone within the region where the instance is placed, null in regions that have none.
    * `networks` - Observed network interfaces of the instance: one entry per private network plus at most one entry without a network id for the public (Ext-Net) interface. Entries are ordered by network id, so they do not follow the order of the requested networks:
      * `id` - Identifier of the network this interface is attached to, null for the public (Ext-Net) interface.
      * `subnet_id` - Identifier of the subnet this interface draws its fixed address from, null for an entry without a network id.
      * `gateway_id` - Identifier of the gateway providing egress for this interface, null when none applies.
      * `addresses` - Addresses observed on this interface: its fixed addresses plus, where applicable, its floating IP and any additional IPs routed to it. Each address carries a type of `FIXED`, `FLOATING` or `ADDITIONAL`:
        * `ip` - IP address assigned to the interface (IPv4 or IPv6).
        * `mac` - MAC address of the interface this IP is bound to. Null when the backend reports no interface for the address, which happens for an additional IP routed to an instance whose public interface has no visible Ext-Net address yet.
        * `type` - How this address reaches the instance: `FIXED` for an address assigned to the interface itself, `FLOATING` for a floating IP NAT'd onto it, `ADDITIONAL` for an additional IP routed to the public interface.
        * `version` - IP version of the address (4 for IPv4, 6 for IPv6).
    * `volumes` - Observed block volumes attached to the instance:
      * `id` - Unique identifier of the attached volume.
      * `name` - Display name of the attached volume.
      * `size` - Size of the attached volume, in GB.
    * `shares` - Observed instance-side share attachments, derived from the Manila access rules that target one of the instance's IPs. Only populated on a single-instance read, so always null here:
      * `id` - Identifier of the attached file storage share.
      * `access_level` - Observed access level of the access rule for this instance (`READ_ONLY` or `READ_WRITE`).
      * `access_to` - The instance IP address the Manila access rule targets: its fixed IPv4 address on the share's network.
      * `state` - Observed state of the underlying access rule (`ACTIVE`, `APPLYING`, `DENYING`, `ERROR`). Null while no state has been reported yet.
    * `security_groups` - Security groups currently attached to the instance's ports:
      * `id` - Security group identifier.
    * `group` - Instance (placement) group the instance belongs to, null when it is not part of any group:
      * `id` - Identifier of the instance (placement) group.
