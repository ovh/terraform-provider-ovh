---
subcategory: "Managed Kubernetes Service (MKS)"
---

# ovh_cloud_managed_kubernetes_plans (Data Source)

Use this data source to list the Managed Kubernetes plans available for cluster creation in a region, with their control plane characteristics and limits.

## Example Usage

```terraform
data "ovh_cloud_managed_kubernetes_plans" "plans" {
  service_name = "<Public cloud project id>"
  region       = "EU-WEST-PAR"
}

output "available_plans" {
  value = [for p in data.ovh_cloud_managed_kubernetes_plans.plans.plans : p.name if p.status == "AVAILABLE"]
}
```

## Argument Reference

The following arguments are supported:

* `service_name` - (Required) Service name of the resource representing the id of the cloud project.
* `region` - (Required) Region to list plans for. Plan availability varies by region.

## Attributes Reference

The following attributes are exported:

* `plans` - Plans available in the region:
  * `name` - Plan name (`FREE`, `STANDARD`).
  * `status` - Availability of the plan for cluster creation (`AVAILABLE`, `UNAVAILABLE`).
  * `control_plane` - Control plane characteristics of the plan:
    * `etcd` - Etcd characteristics of the control plane:
      * `dedicated` - Whether etcd runs on storage dedicated to the cluster.
      * `storage` - Etcd storage quota:
        * `limit` - Storage quota, expressed in `unit`.
        * `unit` - Unit of the storage quota (for example `MiB`, `GiB`).
  * `limits` - Limits enforced by the plan:
    * `max_nodes` - Maximum number of nodes allowed across the cluster.
