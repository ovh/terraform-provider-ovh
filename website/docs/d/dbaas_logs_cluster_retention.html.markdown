---
subcategory : "Logs Data Platform"
---


# ovh_dbaas_logs_cluster_retention (Data Source)

Use this data source to retrieve information about a DBaas logs cluster retention.

## Example Usage

```hcl
data "ovh_dbaas_logs_cluster_retention" "retention" {
  service_name = "ldp-xx-xxxxx"
  cluster_id   = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
  retention_id = "yyyyyyyy-yyyy-yyyy-yyyy-yyyyyyyyyyyy"
}
```

It is also possible to retrieve a retention using its duration:

```hcl
data "ovh_dbaas_logs_cluster_retention" "retention" {
  service_name = "ldp-xx-xxxxx"
  cluster_id   = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
  duration     = "P14D"
}
```

Additionnaly, you can filter retentions on their type:

```hcl
data "ovh_dbaas_logs_cluster_retention" "retention" {
  service_name   = "ldp-xx-xxxxx"
  cluster_id     = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
  duration       = "P14D"
  retention_type = "LOGS_INDEXING"
}
```

## Argument Reference

* `service_name` - (Required) The service name. It's the ID of your Logs Data Platform instance.
* `cluster_id` - (Required) Cluster ID
* `retention_id` - ID of the retention object. Cannot be used if `duration` or `retention_type` is defined.
* `retention_type` - Type of the retention (LOGS_INDEXING | LOGS_COLD_STORAGE | METRICS_TENANT). Cannot be used if `retention_id` is defined. Defaults to `LOGS_INDEXING` if not defined.
* `duration` - Indexed duration expressed in ISO-8601 format. Cannot be used if `retention_id` is defined.

## Attributes Reference

* `retention_id` - ID of the retention that can be used when creating a stream
* `duration` - Indexed duration expressed in ISO-8601 format
* `retention_type` - Type of the retention (LOGS_INDEXING | LOGS_COLD_STORAGE | METRICS_TENANT)
* `is_supported` - Indicates if a new stream can use it