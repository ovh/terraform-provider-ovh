---
subcategory : "Cloud Instances"
---

# ovh_cloud_instance_events (Data Source)

List events (action history) for a compute instance in a public cloud project.

Events are returned from the Nova `os-instance-actions` API and mapped to the `common.Event` schema. Each event represents an action performed on the instance (create, reboot, lock, stop, etc.) with its result status.

## Example Usage

```terraform
data "ovh_cloud_instance_events" "events" {
  service_name = "xxxxxxxxxx"
  instance_id  = "00000000-0000-0000-0000-000000000001"
}

# Output the most recent event
output "latest_event" {
  value = length(data.ovh_cloud_instance_events.events.events) > 0 ? data.ovh_cloud_instance_events.events.events[0] : null
}

# Filter for errors only
output "error_events" {
  value = [for e in data.ovh_cloud_instance_events.events.events : e if e.type == "TASK_ERROR"]
}
```

## Argument Reference

The following arguments are supported:

* `service_name` - (Required) Service name of the resource representing the id of the cloud project.
* `instance_id` - (Required) ID of the compute instance.

## Attributes Reference

The following attributes are exported:

* `events` - List of events for the instance, sorted by most recent first. Each event has the following attributes:
  * `created_at` - Creation date of the event.
  * `kind` - Nature of the event (e.g. `LOCK`, `UNLOCK`, `REBOOT`, `CREATE`, `STOP`, `START`, `RESIZE`, `REBUILD`, `SHELVE`, `UNSHELVE`, `RESCUE`, `UNRESCUE`).
  * `link` - Link to the event related resource (may be null).
  * `message` - Description of what happened on the event.
  * `type` - Type of the event. Possible values are:
    * `TASK_START` — The action has started.
    * `TASK_SUCCESS` — The action completed successfully.
    * `TASK_ERROR` — The action failed.
    * `TARGET_SPEC_UPDATE` — The target specification was updated.
