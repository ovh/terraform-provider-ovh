---
subcategory: "VPS"
page_title: "OVHcloud: ovh_vps_console_url"
description: |-
  Requests a KVM console URL for an OVHcloud VPS.
---

# ovh_vps_console_url (Resource)

Requests a fresh, single-use KVM console URL for an OVHcloud VPS by calling
`POST /vps/{serviceName}/getConsoleUrl`.

~> **WARNING** The URL is **one-shot**: every new resource creation calls
`POST /vps/{serviceName}/getConsoleUrl` and asks OVHcloud for a brand new
signed URL. The URL expires very quickly server-side and is invalidated
once consumed. Do not persist it, share it via long-lived outputs, or
treat it as stable across reads.

## Example Usage

```hcl
resource "ovh_vps_console_url" "console" {
  service_name = "vps-xxxxxx.vps.ovh.net"
}

output "vps_console_url" {
  value     = ovh_vps_console_url.console.url
  sensitive = true
}
```

## Argument Reference

* `service_name` - (Required, ForceNew) The service_name of your VPS.
* `triggers` - (Optional, ForceNew) Arbitrary map of values. Changing any
  value forces a new resource, which re-issues a fresh console URL.

## Attributes Reference

* `id` - Set to the `service_name`.
* `url` - The freshly issued, single-use signed console URL. Sensitive.
