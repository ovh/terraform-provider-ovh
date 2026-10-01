---
subcategory: "Instances"
---

# ovh_cloud_instance_flavors (Data Source)

Use this data source to list the flavors available in a public cloud project. This is read-only reference data: each entry describes a hardware sizing (vCPUs, RAM, disks) the project can create instances from, not an existing resource.

## Example Usage

```terraform
# Omit `region` to list the catalog of every region at once: the same flavor
# then appears once per region it is offered in.
data "ovh_cloud_instance_flavors" "flavors" {
  service_name = "<Public cloud project id>"
  region       = "GRA11"
}

# Resolve a flavor id by its commercial name, to feed ovh_cloud_instance.flavor_id.
output "b3_8_flavor_id" {
  value = one([
    for flavor in data.ovh_cloud_instance_flavors.flavors.flavors :
    flavor.id if flavor.name == "b3-8"
  ])
}

# Every public flavor of the region with at least 16 GB of RAM.
output "large_flavor_names" {
  value = sort([
    for flavor in data.ovh_cloud_instance_flavors.flavors.flavors :
    flavor.name if flavor.is_public && flavor.ram >= 16384
  ])
}
```

### Feeding an instance

```terraform
data "ovh_cloud_instance_flavors" "catalog" {
  service_name = "<Public cloud project id>"
  region       = "GRA11"
}

resource "ovh_cloud_instance" "instance" {
  service_name = "<Public cloud project id>"
  region       = "GRA11"
  name         = "my-instance"
  image_id     = "<image id>"

  flavor_id = one([
    for flavor in data.ovh_cloud_instance_flavors.catalog.flavors :
    flavor.id if flavor.name == "b3-8"
  ])

  networks = [
    { auto_assign_public_ip = true },
  ]
}
```

## Argument Reference

The following arguments are supported:

* `service_name` - (Required) Service name of the resource representing the id of the cloud project.
* `region` - (Optional) Restrict the listing to the flavors offered in this region. The catalog is per-region: a flavor returned for one region is not guaranteed to exist, or to carry the same characteristics, in another.

## Attributes Reference

The following attributes are exported:

* `flavors` - Flavors advertised by the backend for this project, one entry per region the flavor is offered in:
  * `id` - The OpenStack/Nova flavor ID. Stable within a region and used to reference the flavor when creating an instance.
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
