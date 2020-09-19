---
layout: "ovh"
page_title: "OVH: Enterprise Cloud DB"
sidebar_current: "docs-ovh-datasource-enterprise-cloud-db-x"
description: |-
  Get information & status of an Enterprise Cloud DB Instance
---

# ovh_enterprise_cloud_db

Use this data source to retrieve information about an Enterprise Cloud DB Instance.

## Example Usage

```hcl
data "ovh_enterprise_cloud_db" "my-db" {
   cluster_id = "XXXXXX"
}
```

## Argument Reference


* `cluster_id` - (Required) The cluster ID of the Enterprise Cloud DB.


## Attributes Reference

* `region` - Enterprise Cloud DB region 
* `status` - Cluster status
  * created
  * creating
  * deleting
  * reopening
  * restarting
  * scaling
  * suspended
  * suspending
  * updating