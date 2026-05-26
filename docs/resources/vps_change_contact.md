---
subcategory : "VPS"
---

# ovh_vps_change_contact (Resource)

Submit a contact change request for an OVHcloud VPS by calling `POST /vps/{serviceName}/changeContact`. This is a one-shot action: at least one of `contact_admin`, `contact_billing`, `contact_tech` must be provided. Changing any argument forces a new resource (a new request is submitted to the OVH API).

~> **NOTE** Contact changes cannot be rolled back from this resource. `terraform destroy` only removes the resource from state.

## Example Usage

```terraform
resource "ovh_vps_change_contact" "to_new_admin" {
  service_name  = "vpsXXXXX.ovh.net"
  contact_admin = "ab12345-ovh"
}
```

## Argument Reference

* `service_name` - (Required, ForceNew) The service name of your VPS.
* `contact_admin` - (Optional, ForceNew) The OVH nichandle of the new admin contact.
* `contact_billing` - (Optional, ForceNew) The OVH nichandle of the new billing contact.
* `contact_tech` - (Optional, ForceNew) The OVH nichandle of the new tech contact.

At least one of `contact_admin`, `contact_billing` or `contact_tech` must be set.

## Attributes Reference

* `task_ids` - The list of task IDs returned by the OVH API for the submitted change request(s).
