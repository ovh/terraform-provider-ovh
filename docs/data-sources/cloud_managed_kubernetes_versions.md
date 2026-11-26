---
subcategory: "Managed Kubernetes Service (MKS)"
---

# ovh_cloud_managed_kubernetes_versions (Data Source)

Use this data source to list the Kubernetes minor versions available for Managed Kubernetes cluster creation, and the default version.

## Example Usage

```terraform
data "ovh_cloud_managed_kubernetes_versions" "versions" {
  service_name = "<Public cloud project id>"
}

output "default_version" {
  value = data.ovh_cloud_managed_kubernetes_versions.versions.default_version
}
```

## Argument Reference

The following arguments are supported:

* `service_name` - (Required) Service name of the resource representing the id of the cloud project.

## Attributes Reference

The following attributes are exported:

* `default_version` - Minor version applied when a creation request does not specify one (for example `1.34`).
* `available_versions` - Minor versions available for cluster creation, newest first. Creation resolves a minor version to its latest patch:
  * `version` - Kubernetes minor version (for example `1.34`).
  * `status` - Support status of the version (`SUPPORTED`, `END_OF_SUPPORT`).
