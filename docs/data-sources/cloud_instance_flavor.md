---
subcategory: "Instances"
---

# ovh_cloud_instance_flavor (Data Source)

Use this data source to retrieve information about a flavor available in a public cloud project. This is read-only reference data: it describes a hardware sizing (vCPUs, RAM, disks) the project can create instances from, not an existing resource.

## Example Usage

```terraform
data "ovh_cloud_instance_flavor" "flavor" {
  service_name = "<Public cloud project id>"
  id           = "<flavor id>"
}

output "flavor_sizing" {
  value = format(
    "%s: %d vCPUs, %d MB RAM, %d GB disk",
    data.ovh_cloud_instance_flavor.flavor.name,
    data.ovh_cloud_instance_flavor.flavor.vcpus,
    data.ovh_cloud_instance_flavor.flavor.ram,
    data.ovh_cloud_instance_flavor.flavor.disk,
  )
}
```

To look a flavor up by name instead of by id, use
[`ovh_cloud_instance_flavors`](cloud_instance_flavors.md).

## Argument Reference

The following arguments are supported:

* `service_name` - (Required) Service name of the resource representing the id of the cloud project.
* `id` - (Required) The OpenStack/Nova flavor ID to look up. Stable within a region and used to reference the flavor when creating an instance.

## Attributes Reference

The following attributes are exported:

* `name` - The backend flavor name (for example `b2-7`, `c2-15`). This is the commercial/technical name used to identify the sizing in the catalog.
* `vcpus` - Number of virtual CPUs provided by the flavor.
* `ram` - Amount of memory provided by the flavor, expressed in MB.
* `disk` - Size of the flavor's root disk in GB. This is the primary system disk provisioned for instances created from this flavor.
* `swap` - Size of the flavor's swap space in MB. Zero when the flavor provides no swap.
* `ephemeral` - Size of the flavor's ephemeral disk in GB. Ephemeral storage is transient: its contents do not survive a rebuild or deletion of the instance. Zero when the flavor provides no ephemeral disk.
* `is_public` - Whether the flavor is publicly available to the project. Private flavors are only visible to the projects they have been explicitly shared with.
* `description` - Free-form description of the flavor as reported by the backend. May be empty when no description is advertised for this flavor.
* `location` - Region (and, where applicable, availability zone) where this flavor is offered. The flavor catalog is per-region: a flavor returned for one region is not guaranteed to exist, or to carry the same characteristics, in another:
  * `region` - Region code.
  * `availability_zone` - Availability zone within the region.
