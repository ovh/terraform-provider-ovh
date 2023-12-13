---
subcategory : "Logs Data Platform"
---


# ovh_dbaas_logs_clusters (Data Source)

Use this data source to retrieve UUIDs of DBaas logs clusters.

## Example Usage

```hcl
data "ovh_dbaas_logs_clusters" "logstash" {
  service_name = "ldp-xx-xxxxx"
}
```

## Argument Reference

* `service_name` - The service name. It's the ID of your Logs Data Platform instance.

## Attributes Reference

* `uuids` is the cluster id
