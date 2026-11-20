---
subcategory: "Managed Kubernetes Service (MKS)"
---

# ovh_cloud_managed_kubernetes_nodepool

Creates a nodepool in a Managed Kubernetes cluster created with
[`ovh_cloud_managed_kubernetes`](cloud_managed_kubernetes.md).

## Example Usage

```terraform
resource "ovh_cloud_managed_kubernetes_nodepool" "pool" {
  service_name          = "<Public cloud project id>"
  managed_kubernetes_id = ovh_cloud_managed_kubernetes.cluster.id
  name                  = "my-pool"
  desired_nodes         = 3

  node_config = { flavor = "b3-8" }

  autoscaler = {
    enabled   = true
    min_nodes = 1
    max_nodes = 5
  }

  template = {
    metadata = { labels = { team = "backend" } }
  }

  # The autoscaler rewrites desired_nodes: keep the configured value as the initial size only.
  lifecycle {
    ignore_changes = [desired_nodes]
  }
}
```

## Argument Reference

The following arguments are supported:

* `service_name` - (Required) Service name of the resource representing the id of the cloud project. **Changing this value recreates the resource.**
* `managed_kubernetes_id` - (Required) ID of the Managed Kubernetes cluster. **Changing this value recreates the resource.**
* `node_config` - (Required) Hardware configuration of the nodes.
  * `flavor` - (Required) Node flavor, for example `b3-8`. **Changing this value recreates the resource.**
  * `attach_floating_ip` - (Optional) `enabled`: whether a floating IP is attached to each node. Terraform will only perform drift detection of its value when present in a configuration. **Changing a configured value recreates the resource.**
* `desired_nodes` - (Optional) Desired number of nodes. Defaults to `0`. With the autoscaler enabled, must be at least `autoscaler.min_nodes`. The autoscaler then updates this value: add `desired_nodes` to `lifecycle.ignore_changes` to avoid a diff on every plan.
* `name` - (Optional) Nodepool name. Generated when omitted. **Changing a configured value recreates the resource.**
* `availability_zone` - (Optional) Availability zone the nodes are placed in. Defaults to the zone chosen by the API. Terraform will only perform drift detection of its value when present in a configuration. **Changing a configured value recreates the resource.**
* `billing_type` - (Optional) Billing period of the nodes: `HOURLY` or `MONTHLY`. Defaults to the period chosen by the API. Terraform will only perform drift detection of its value when present in a configuration. **Changing a configured value recreates the resource.**
* `anti_affinity` - (Optional) Anti-affinity placement of the nodes. Terraform will only perform drift detection of its value when present in a configuration. **Changing a configured value recreates the resource.**
  * `enabled` - (Optional) Whether anti-affinity placement is enabled.
  * `policy` - (Optional) `HARD` spreads nodes across distinct hypervisors and caps the nodepool size.
* `autoscaler` - (Optional) Automatic node scaling. Defaults to the thresholds chosen by the API. Terraform will only perform drift detection of its value when present in a configuration.
  * `enabled` - (Optional) Whether automatic scaling is enabled.
  * `min_nodes` - (Optional) Minimum number of nodes.
  * `max_nodes` - (Optional) Maximum number of nodes.
  * `scale_down` - (Optional) `unneeded_time_seconds`, `unready_time_seconds` and `utilization_threshold` (ratio between 0 and 1, defaults to 0.5).
* `template` - (Optional) Kubernetes metadata and scheduling settings applied to every node. Terraform will only perform drift detection of its value when present in a configuration.
  * `metadata` - (Optional) `labels` and `annotations` (maps of strings), `finalizers` (list of strings).
  * `spec` - (Optional) `taints` (list of `{ key, value, effect }`, `effect` being `NoExecute`, `NoSchedule` or `PreferNoSchedule`) and `unschedulable`.

An optional argument omitted from the configuration keeps the value the API applies:
its default, or the value found on import. Terraform only performs drift detection
of an argument when it is present in the configuration, and removing an argument
from the configuration does not reset it: set it explicitly to change it, or set an
empty value (`taints = []`, `labels = {}`) to clear it.

## Attributes Reference

The following attributes are exported:

* `id` - Nodepool ID.
* `checksum` - Computed hash of the target spec.
* `created_at` - Creation date.
* `updated_at` - Last update date.
* `resource_status` - Nodepool readiness in the system (`CREATING`, `DELETING`, `ERROR`, `OUT_OF_SYNC`, `READY`, `UPDATING`…).
* `current_state` - Observed state of the nodepool. Same shape as the arguments above (`name`, `availability_zone`, `billing_type`, `desired_nodes`, `node_config`, `anti_affinity`, `autoscaler`, `template`), plus:
  * `rollout.strategy` - Strategy used to apply node changes: `IN_PLACE` or `ROLLING_UPDATE`.
  * `status` - `current_nodes`, `ready_nodes` and `up_to_date_nodes`.

## Import

A nodepool can be imported using the `service_name`, the `managed_kubernetes_id`
and the nodepool `id`, separated by `/`:

```terraform
import {
  to = ovh_cloud_managed_kubernetes_nodepool.pool
  id = "<service_name>/<managed_kubernetes_id>/<nodepool_id>"
}
```

```bash
$ terraform import ovh_cloud_managed_kubernetes_nodepool.pool service_name/managed_kubernetes_id/nodepool_id
```
