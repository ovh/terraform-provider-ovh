---
subcategory: "Managed Kubernetes Service (MKS)"
---

# ovh_cloud_managed_kubernetes_regions (Data Source)

Use this data source to list the regions where a Managed Kubernetes cluster can be created, with their availability zones and the plans offered there.

## Example Usage

```terraform
data "ovh_cloud_managed_kubernetes_regions" "regions" {
  service_name = "<Public cloud project id>"
}

# Regions offering the STANDARD plan.
output "standard_regions" {
  value = [for r in data.ovh_cloud_managed_kubernetes_regions.regions.regions : r.name if contains(r.available_plans, "STANDARD")]
}
```

## Argument Reference

The following arguments are supported:

* `service_name` - (Required) Service name of the resource representing the id of the cloud project.

## Attributes Reference

The following attributes are exported:

* `regions` - Regions available for cluster creation:
  * `name` - Region identifier (for example `EU-WEST-PAR`).
  * `availability_zones` - Availability zones in the region. Empty for a single-zone region.
  * `available_plans` - Plans that can be selected when creating a cluster in this region.
