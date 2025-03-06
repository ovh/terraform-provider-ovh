---
subcategory : "Domain names"
---

# ovh_domain_name

Create and manage a domain name.

## Important

-> **NOTE** To order a product through Terraform, your account needs to have a default payment method defined. This can be done in the [OVHcloud Control Panel](https://www.ovh.com/manager/#/dedicated/billing/payment/method) or via API with the [/me/payment/method](https://api.ovh.com/console/#/me/payment/method~GET) endpoint.

~> **WARNING** `BANK_ACCOUNT` is not supported anymore, please update your default payment method to `SEPA_DIRECT_DEBIT`

## Example Usage

```terraform
resource "ovh_domain_name" "domain" {
  domain_name = "example.com"

  target_spec = {
    dns_configuration = {
      name_servers = [
        {
          name_server = "dns101.ovh.net"
        },
        {
          name_server = "ns101.ovh.net"
        }
      ]
    }
  }
}
```

## Schema

### Required

- `domain_name` (String) Domain name

### Optional

- `checksum` (String) Computed hash used to control concurrent modification requests. Here, it represents the current target specification value
- `ovh_subsidiary` (String) OVH subsidiaries
- `plan` (Attributes List) (see [below for nested schema](#nestedatt--plan))
- `plan_option` (Attributes List) (see [below for nested schema](#nestedatt--plan_option))
- `target_spec` (Attributes) Latest target specification of the domain name resource. (see [below for nested schema](#nestedatt--target_spec))

### Read-Only

- `current_state` (Attributes) Current state of the domain name (see [below for nested schema](#nestedatt--current_state))
- `current_tasks` (Attributes List) Ongoing asynchronous tasks related to the domain name resource (see [below for nested schema](#nestedatt--current_tasks))
- `iam` (Attributes) IAM resource metadata (see [below for nested schema](#nestedatt--iam))
- `id` (String) Unique identifier for the resource. Here, the domain name itself is used as an identifier
- `order` (Attributes) Details about an Order (see [below for nested schema](#nestedatt--order))
- `resource_status` (String) Reflects the readiness of the domain name resource. A new target specification request will be accepted only in `READY`, `UPDATING` or `ERROR` status

<a id="nestedatt--plan"></a>

### Nested Schema for `plan`

Required:

- `duration` (String) Duration selected for the purchase of the product (defaults to "P1Y")
- `plan_code` (String) Identifier of the option offer
- `pricing_mode` (String) Pricing mode selected for the purchase of the product

Optional:

- `configuration` (Attributes List) (see [below for nested schema](#nestedatt--plan--configuration))
- `item_id` (Number) Cart item to be linked
- `quantity` (Number) Quantity of product desired

<a id="nestedatt--plan--configuration"></a>

### Nested Schema for `plan.configuration`

Required:

- `label` (String) Label for your configuration item
- `value` (String) Value or resource URL on API.OVH.COM of your configuration item

<a id="nestedatt--plan_option"></a>

### Nested Schema for `plan_option`

Required:

- `duration` (String) Duration selected for the purchase of the product
- `plan_code` (String) Identifier of the option offer
- `pricing_mode` (String) Pricing mode selected for the purchase of the product
- `quantity` (Number) Quantity of product desired

Optional:

- `configuration` (Attributes List) (see [below for nested schema](#nestedatt--plan_option--configuration))

<a id="nestedatt--plan_option--configuration"></a>

### Nested Schema for `plan_option.configuration`

Required:

- `label` (String) Label for your configuration item
- `value` (String) Value or resource URL on API.OVH.COM of your configuration item

<a id="nestedatt--target_spec"></a>

### Nested Schema for `target_spec`

Optional:

- `dns_configuration` (Attributes) The domain DNS configuration (see [below for nested schema](#nestedatt--target_spec--dns_configuration))

<a id="nestedatt--target_spec--dns_configuration"></a>

### Nested Schema for `target_spec.dns_configuration`

Optional:

- `name_servers` (Attributes List) The name servers to update (see [below for nested schema](#nestedatt--target_spec--dns_configuration--name_servers))

<a id="nestedatt--target_spec--dns_configuration--name_servers"></a>

### Nested Schema for `target_spec.dns_configuration.name_servers`

Optional:

- `ipv4` (String) The IPv4 associated to the name server
- `ipv6` (String) The IPv6 associated to the name server
- `name_server` (String) The host name

<a id="nestedatt--current_state"></a>

### Nested Schema for `current_state`

Read-Only:

- `additional_states` (List of String) Domain additional states
- `dns_configuration` (Attributes) The domain DNS configuration (see [below for nested schema](#nestedatt--current_state--dns_configuration))
- `extension` (String) Extension of the domain name
- `main_state` (String) Domain main state
- `name` (String) Domain name
- `protection_state` (String) Domain protection state
- `suspension_state` (String) Domain suspension state

<a id="nestedatt--current_state--dns_configuration"></a>

### Nested Schema for `current_state.dns_configuration`

Read-Only:

- `configuration_type` (String) The type of DNS configuration of the domain
- `glue_record_ipv6supported` (Boolean) Whether the registry supports IPv6 or not
- `host_supported` (Boolean) Whether the registry accepts hosts or not
- `max_dns` (Number) The maximum number of name servers allowed by the registry
- `min_dns` (Number) The minimum number of name servers allowed by the registry
- `name_servers` (Attributes List) The name servers used by the domain name (see [below for nested schema](#nestedatt--current_state--dns_configuration--name_servers))

<a id="nestedatt--current_state--dns_configuration--name_servers"></a>

### Nested Schema for `current_state.dns_configuration.name_servers`

Read-Only:

- `ipv4` (String) The IPv4 associated to the name server
- `ipv6` (String) The IPv6 associated to the name server
- `name_server` (String) The host name
- `name_server_type` (String) The type of name server

<a id="nestedatt--current_tasks"></a>

### Nested Schema for `current_tasks`

Read-Only:

- `id` (String) Identifier of the current task
- `link` (String) Link to the task details
- `status` (String) Current global status of the current task
- `type` (String) Type of the current task

<a id="nestedatt--iam"></a>

### Nested Schema for `iam`

Read-Only:

- `display_name` (String) Resource display name
- `id` (String) Unique identifier of the resource
- `tags` (Map of String) Resource tags. Tags that were internally computed are prefixed with ovh:
- `urn` (String) Unique resource name used in policies

<a id="nestedatt--order"></a>

### Nested Schema for `order`

Read-Only:

- `date` (String)
- `details` (Attributes List) (see [below for nested schema](#nestedatt--order--details))
- `expiration_date` (String)
- `order_id` (Number)

<a id="nestedatt--order--details"></a>

### Nested Schema for `order.details`

Read-Only:

- `description` (String)
- `detail_type` (String) Product type of item in order
- `domain` (String)
- `order_detail_id` (Number)
- `quantity` (String)

## Import

A domain name can be imported using its `domain_name`.

Using the following configuration:

```terraform
import {
  to = ovh_domain_name.domain
  id = "<domain name>"
}
```

You can then run:

```bash
$ terraform plan -generate-config-out=domain.tf
$ terraform apply
```

The file `domain.tf` will then contain the imported resource's configuration, that can be copied next to the `import` block above. See https://developer.hashicorp.com/terraform/language/import/generating-configuration for more details.
