---
subcategory : "Workflow Management"
---

# ovh_cloud_project_workflow_backup

Manage a worflow that schedules backups of public cloud instance. Note that upon deletion, the workflow is deleted but any backups that have been created by this workflow are not.

## Example Usage

```terraform
resource "ovh_cloud_project_workflow_backup" "my_backup" {
  service_name        = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
  region_name         = "GRA11"
  cron                = "50 4 * * *"
  instance_id         = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx"
  max_execution_count = "0"
  name                = "Backup workflow for instance"
  rotation            = "7"
}
```

## Argument Reference

The following arguments are supported:

* `service_name` - (Optional) The id of the public cloud project. If omitted, the `OVH_CLOUD_PROJECT_SERVICE` environment variable is used.

* `region_name` - (Mandatory) The name of the openstack region.

* `cron` - (Mandatory) The cron periodicity at which the backup workflow is scheduled

* `instanceId` the id of the instance to back up

* `max_execution_count` - (Optional) The number of times the worflow is run. Default value is `0` which means that the workflow will be scheduled continously until its deletion

* `name` - (Mandatory) The worflow name that is used in the UI
* `rotation`- (Mandatory) The number of backup that are retained.
* `backup_name` - (Optional) The name of the backup files that are created. If empty, the `name` attribute is used.

## Attributes Reference

The following attributes are exported:

* `id` - The workflow ID.
* `service_name` - See Argument Reference above.
* `region_name` - See Argument Reference above.
* `cron` - See Argument Reference above.
* `instance_id` - See Argument Reference above.
* `max_execution_count` - See Argument Reference above.
* `name` - See Argument Reference above.
* `rotation` - See Argument Reference above.
* `backup_name` - See Argument Reference above.
* `created_at` - The creation date of the workflow.

## Import

A cloud project workflow backup can be imported using the `id` (workflow ID) E.g.,

```bash
$ terraform import ovh_cloud_project_workflow_backup.my_workflow id
```
