---
layout: "ovh"
page_title: "OVH: cloud_project_kube_oidc"
sidebar_current: "docs-ovh-resource-cloud-project-kube-oidc"
description: |-
Creates an OIDC configuration in an OVHcloud Managed Kubernetes cluster.
---

# ovh_cloud_project_kube_oidc

Creates an OIDC configuration in an OVHcloud Managed Kubernetes cluster.

## Example Usage

```hcl
resource "ovh_cloud_project_kube_oidc" "my-oidc" {
  service_name = var.projectid
  kube_id      = ovh_cloud_project_kube.k8stf.id
  client_id    = "xxx"
  issuer_url   = "https://ovh.com"
}
```

## Argument Reference

The following arguments are supported:

* `service_name` - (required) The id of the public cloud project. If omitted,
  the `OVH_CLOUD_PROJECT_SERVICE` environment variable is used.

* `kube_id` - The id of the managed kubernetes cluster.

* `client_id` - The OIDC client ID.

* `issuer_url` - The OIDC issuer url.
