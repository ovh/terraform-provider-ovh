---
subcategory : "Dedicated Server"
---

# ovh_dedicated_server_reboot_task

Reboot your Dedicated Server.

~> **WARNING** After some delay, if the task is marked as `done`, the Provider may purge it. To avoid raising errors when terraform refreshes its plan, 404 errors are ignored on Resource Read, thus some information may be lost after a while.

## Example Usage

```terraform
data "ovh_dedicated_server_boots" "rescue" {
  service_name = "nsxxxxxxx.ip-xx-xx-xx.eu"
  boot_type    = "rescue"
  kernel       = "rescue64-pro"
}

resource "ovh_dedicated_server_update" "server_on_rescue" {
  service_name = "nsxxxxxxx.ip-xx-xx-xx.eu"
  boot_id      = data.ovh_dedicated_server_boots.rescue.result[0]
  monitoring   = true
  state        = "ok"
}

resource "ovh_dedicated_server_reboot_task" "server_reboot" {
  service_name = data.ovh_dedicated_server_boots.rescue.service_name

  keepers = [
     ovh_dedicated_server_update.server_on_rescue.boot_id,
  ]
}
```

## Argument Reference

The following arguments are supported:

* `service_name` - (Required) The service_name of your dedicated server.
* `keepers` - List of values tracked to trigger reboot, used also to form implicit dependencies.

## Attributes Reference

The following attributes are exported:

* `id` - The task id
* `comment` - Details of this task. (should be `Reboot asked`)
* `done_date` - Completion date in RFC3339 format.
* `function` - Function name (should be `hardReboot`).
* `last_update` - Last update in RFC3339 format.
* `start_date` - Task creation date in RFC3339 format.
* `status` - Task status (should be `done`)
