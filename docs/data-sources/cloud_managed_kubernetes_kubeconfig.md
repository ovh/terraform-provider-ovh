---
subcategory: "Managed Kubernetes Service (MKS)"
---

# ovh_cloud_managed_kubernetes_kubeconfig (Data Source)

Use this data source to retrieve the admin kubeconfig of a Managed Kubernetes cluster.

## Example Usage

```terraform
data "ovh_cloud_managed_kubernetes_kubeconfig" "cluster" {
  service_name          = "<Public cloud project id>"
  managed_kubernetes_id = "<cluster id>"
}

resource "local_sensitive_file" "kubeconfig" {
  content  = data.ovh_cloud_managed_kubernetes_kubeconfig.cluster.content
  filename = "${path.module}/kubeconfig.yaml"
}
```

## Argument Reference

* `service_name` - (Required) Service name of the resource representing the id of the cloud project.
* `managed_kubernetes_id` - (Required) ID of the [`ovh_cloud_managed_kubernetes`](../resources/cloud_managed_kubernetes.md) cluster.

## Attributes Reference

* `content` - (Sensitive) Admin kubeconfig of the cluster, in YAML. Reading it requires the cluster to be `READY`.

~> **Note** The kubeconfig is stored in the Terraform state. Protect the state accordingly.
