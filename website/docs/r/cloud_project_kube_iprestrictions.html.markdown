---
layout: "ovh"
page_title: "OVH: cloud_project_kube_iprestrictions"
sidebar_current: "docs-ovh-resource-cloud-project-kube-iprestrictions-x"
description: |-
  Apply IP restrictions to a managed Kubernetes cluster.
---

# ovh_cloud_project_kube_iprestrictions

Apply IP restrictions to a managed Kubernetes cluster.

## Example Usage

```hcl
resource "ovh_cloud_project_kube_iprestrictions" "vrack_only" {
   service_name = "94d423da0e5545f29812836460a19939"
   kube_id      = "9260267d-2bf9-4d9a-bb6e-24b8969f65e2"
   ips          = toset(["10.42.0.0/16"])
}
```

## Argument Reference

The following arguments are supported:

* `service_name` - (Optional) The id of the public cloud project. If omitted,
    the `OVH_CLOUD_PROJECT_SERVICE` environment variable is used.

* `kube_id` - The id of the managed Kubernetes cluster.

* `ips` - List of CIDR authorized to interact with the managed Kubernetes cluster.

## Attributes Reference

No additional attributes than the ones provided are exported.
