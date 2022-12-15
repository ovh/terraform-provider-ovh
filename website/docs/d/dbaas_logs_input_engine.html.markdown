---
layout: "ovh"
page_title: "OVH: dbaas_logs_input_engine"
sidebar_current: "docs-ovh-datasource-dbaas-logs-input-engine"
description: |-
  Get information & status of a DBaas logs input engine.
---

# ovh_dbaas_logs_input_engine (Data Source)

Use this data source to retrieve information about a DBaas logs input engine.

## Example Usage

```hcl
data "ovh_dbaas_logs_input_engine" "logstash" {
  service_name = "ldp-xx-xxxxx"
  name          = "logstash"
  version       = "6.8"
  is_deprecated = true
}
```

## Argument Reference

* `service_name` - The service name. It's the ID of your Logs Data Platform instance.
* `name` - The name of the logs input engine.
* `is_deprecated` - Indicates if engine will soon not be supported.
* `version` - Software version

## Attributes Reference

`id` is set to input engine ID