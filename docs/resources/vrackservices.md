---
subcategory: "Vrack Services"
---

# ovh_vrackservices (Resource)

Orders a Vrack Services.

## Important

-> **NOTE** To order a product through Terraform, your account needs to have a default payment method defined. This can be done in the [OVHcloud Control Panel](https://www.ovh.com/manager/#/dedicated/billing/payment/method) or via API with the [/me/payment/method](https://api.ovh.com/console/#/me/payment/method~GET) endpoint.

-> **NOTE** This resource is created through the OVHcloud order workflow and deleted through the termination workflow, even when the product is free. If your consumer key has restricted access rules, it must allow the routes listed in [Consumer key access rules for resources that place orders](../guides/order_consumer_key_access_rules), in addition to `/vrackServices/*`.

~> **WARNING** `BANK_ACCOUNT` is not supported anymore, please update your default payment_method to `SEPA_DIRECT_DEBIT`

## BGP option

BGP is an option of the vRack Services, ordered through `plan_option` with the `bgp-beta-vrack-services` plan code (`bgp-beta-vrack-services-eu` or `bgp-beta-vrack-services-ca` for US customers ordering in Europe or Canada).

- Adding the BGP option to an existing `ovh_vrackservices` orders it on the delivered service.
- Setting it at creation orders the vRack Services and its BGP option at once.
- Removing it requests the termination of the BGP option, as destroying the vRack Services requests its own: the apply returns once the termination is confirmed, and the BGP is suspended then deleted later. Until then, `current_state.bgp` and `target_spec.bgp` stay visible.

`target_spec.bgp` can only be set while the BGP option is subscribed, and the API refuses to remove it while subscribed: when `target_spec.bgp`, its ASNs or `allowed_prefixes` are omitted, their current value is kept. To reset the BGP configuration, empty `allowed_prefixes` and `bgp_range`.

~> **WARNING** Keep `bgp_range` in your configuration while it holds ranges, including when you only change `allowed_prefixes` or the ASNs. Omitting a non-empty `bgp_range`, or the whole `target_spec.bgp`, shows a change on every plan: Terraform plans its removal, while the provider keeps the current ranges.

A BGP range references an untagged subnet: the subnet `vlan` must be `null`.

## Upgrading from the previous schema

The `ovh_vrackservices` schema follows the vRack Services API revision that adds BGP. An existing configuration must be rewritten, but the Vrack Services is neither replaced nor updated.

| Previous attribute | New attribute |
|---|---|
| `target_spec.subnets[].service_range` | `target_spec.services.service_range[]`, linked to its subnet by `subnet_cidr` |
| `target_spec.subnets[].service_endpoints` | `target_spec.services.service_range[].service_endpoints` |
| — | `target_spec.subnets[].az` (required) |
| `current_state.product_status` | `current_state.configuration_status` |

Before:

```terraform
target_spec = {
  subnets = [
    {
      cidr              = "192.168.0.0/24"
      service_range     = { cidr = "192.168.0.0/29" }
      service_endpoints = []
    },
  ]
}
```

After:

```terraform
target_spec = {
  subnets = [
    {
      az   = "eu-west-lim-a"
      cidr = "192.168.0.0/24"
    },
  ]
  services = {
    service_range = [
      {
        cidr              = "192.168.0.0/29"
        subnet_cidr       = "192.168.0.0/24"
        service_endpoints = []
      },
    ]
  }
}
```

To upgrade:

1. Rewrite `target_spec` as above.
2. Run `terraform plan`: it shows no changes, as the new attributes are read from the API.
3. Run `terraform apply -refresh-only` once to save the state with the new schema. Until then, `terraform show` and `terraform state show` fail.

## Example Usage

### Example 1 - Vrack Services ordered by Terraform

```terraform
locals {
  region = "eu-west-lim"
}

data "ovh_me" "myaccount" {}

resource "ovh_vrackservices" "my-vrackservices" {
  ovh_subsidiary = data.ovh_me.myaccount.ovh_subsidiary
  plan = [
    {
      plan_code = "vrack-services"
      duration = "P1M"
      pricing_mode = "default"

      configuration = [
        {
          label = "region_name"
          value = local.region
        }
      ]
    }
  ]
  target_spec = {
    subnets = []
  }
}
```

### Example 2 - Vrack Services created by Terraform with a basic configuration

```terraform
locals {
  region = "eu-west-lim"
}

data "ovh_me" "myaccount" {}

resource "ovh_vrackservices" "my-vrackservices" {
  ovh_subsidiary = data.ovh_me.myaccount.ovh_subsidiary
  plan = [
    {
      plan_code = "vrack-services"
      duration = "P1M"
      pricing_mode = "default"

      configuration = [
        {
          label = "region_name"
          value = local.region
        }
      ]
    }
  ]
  target_spec = {
    subnets = [
      {
        az   = "eu-west-lim-a"
        cidr = "192.168.0.0/24"
      },
    ]
    services = {
      service_range = [
        {
          cidr              = "192.168.0.0/29"
          subnet_cidr       = "192.168.0.0/24"
          service_endpoints = []
        },
      ]
    }
  }
}
```

### Example 3 - Vrack Services created by Terraform and associated to a vRack

```terraform
locals {
  region = "eu-west-lim"
  vrack_name = "pn-000000"
}

data "ovh_me" "myaccount" {}

resource "ovh_vrackservices" "my-vrackservices" {
  ovh_subsidiary = data.ovh_me.myaccount.ovh_subsidiary
  plan = [
    {
      plan_code = "vrack-services"
      duration = "P1M"
      pricing_mode = "default"

      configuration = [
        {
          label = "region_name"
          value = local.region
        }
      ]
    }
  ]
  target_spec = {
    subnets = []
  }
}

resource "ovh_vrack_vrackservices" "vrack-vrackservices-binding" {
  service_name   = local.vrack_name
  vrack_services = ovh_vrackservices.my-vrackservices.id
}

```

### Example 4 - Vrack Services created by Terraform with a managed service

```terraform
locals {
  region = "eu-west-lim"
  efs_name = "example-efs-service-name-000e75d3d4c1"
}

data "ovh_me" "myaccount" {}

data "ovh_storage_efs" "my-efs" {
  service_name = local.efs_name
}

resource "ovh_vrackservices" "my-vrackservices" {
  ovh_subsidiary = data.ovh_me.myaccount.ovh_subsidiary
  plan = [
    {
      plan_code = "vrack-services"
      duration = "P1M"
      pricing_mode = "default"

      configuration = [
        {
          label = "region_name"
          value = local.region
        }
      ]
    }
  ]
  target_spec = {
    subnets = [
      {
        az           = "eu-west-lim-a"
        cidr         = "192.168.0.0/24"
        display_name = "my.subnet"
        vlan         = 30
      },
    ]
    services = {
      service_range = [
        {
          cidr        = "192.168.0.0/29"
          subnet_cidr = "192.168.0.0/24"
          service_endpoints = [
            {
              managed_service_urn = data.ovh_storage_efs.my-efs.iam.urn
            }
          ]
        },
      ]
    }
  }
}
```

### Example 5 - Vrack Services created by Terraform with a complete configuration

```terraform
# Once this plan executed, your ovh_vrackservices resource must be updated in your state using :
#   `terraform plan -refresh-only`
#   `terraform apply -refresh-only -auto-approve`

locals {
  region = "eu-west-lim"
  efs_name = "example-efs-service-name-000e75d3d4c1"
  vrack_name = "pn-000000"
}

data "ovh_me" "myaccount" {}

data "ovh_storage_efs" "my-efs" {
  service_name = local.efs_name
}

resource "ovh_vrackservices" "my-vrackservices" {
  ovh_subsidiary = data.ovh_me.myaccount.ovh_subsidiary
  plan = [
    {
      plan_code = "vrack-services"
      duration = "P1M"
      pricing_mode = "default"

      configuration = [
        {
          label = "region_name"
          value = local.region
        }
      ]
    }
  ]
  target_spec = {
    subnets = [
      {
        az           = "eu-west-lim-a"
        cidr         = "192.168.0.0/24"
        display_name = "my.subnet"
        vlan         = 30
      },
    ]
    services = {
      service_range = [
        {
          cidr        = "192.168.0.0/29"
          subnet_cidr = "192.168.0.0/24"
          service_endpoints = [
            {
              managed_service_urn = data.ovh_storage_efs.my-efs.iam.urn
            }
          ]
        },
      ]
    }
  }
}

resource "ovh_vrack_vrackservices" "vrack-vrackservices-binding" {
  service_name   = local.vrack_name
  vrack_services = ovh_vrackservices.my-vrackservices.id
}

```

### Example 6 - BGP option ordered on an existing Vrack Services created by Terraform

Add the `plan_option` block to an existing `ovh_vrackservices` to order its BGP option, then configure BGP through `target_spec.bgp`. Removing the `plan_option` block requests the termination of the BGP option.

```terraform
locals {
  region = "eu-west-lim"
}

data "ovh_me" "myaccount" {}

resource "ovh_vrackservices" "my-vrackservices" {
  ovh_subsidiary = data.ovh_me.myaccount.ovh_subsidiary
  plan = [
    {
      plan_code = "vrack-services"
      duration = "P1M"
      pricing_mode = "default"

      configuration = [
        {
          label = "region_name"
          value = local.region
        }
      ]
    }
  ]
  plan_option = [
    {
      plan_code    = "bgp-beta-vrack-services"
      duration     = "P1M"
      pricing_mode = "default"
      quantity     = 1
    }
  ]
  target_spec = {
    subnets = [
      {
        az   = "eu-west-lim-a"
        cidr = "10.0.0.0/24"
      },
    ]
    bgp = {
      client_asn         = 64600
      infrastructure_asn = 65000
      allowed_prefixes = {
        ipv4 = []
        ipv6 = []
      }
      bgp_range = [
        {
          cidr                 = "10.0.0.0/28"
          subnet_cidr          = "10.0.0.0/24"
          client_peers         = ["10.0.0.5"]
          infrastructure_peers = ["10.0.0.1", "10.0.0.2", "10.0.0.3", "10.0.0.4"]
        },
      ]
    }
  }
}
```

### Example 7 - BGP option ordered on an existing Vrack Services imported into Terraform

A single `terraform apply` imports the Vrack Services, orders its BGP option, then applies the BGP configuration.

- Leave `plan` and `ovh_subsidiary` out: they are empty after an import, and setting them forces the Vrack Services to be terminated and ordered again.
- Keep `target_spec.subnets` as it currently is, or the apply also updates the subnets.
- Once imported, the `import` block can be removed.

```terraform
locals {
  vrack_services_id = "vrs-xxx-xxx-xxx-xxx"
}

import {
  to = ovh_vrackservices.my-vrackservices
  id = local.vrack_services_id
}

resource "ovh_vrackservices" "my-vrackservices" {
  plan_option = [
    {
      plan_code    = "bgp-beta-vrack-services"
      duration     = "P1M"
      pricing_mode = "default"
      quantity     = 1
    }
  ]
  target_spec = {
    subnets = [
      {
        az   = "eu-west-lim-a"
        cidr = "10.0.0.0/24"
      },
    ]
    bgp = {
      bgp_range = [
        {
          cidr                 = "10.0.0.0/28"
          subnet_cidr          = "10.0.0.0/24"
          client_peers         = ["10.0.0.5"]
          infrastructure_peers = ["10.0.0.1", "10.0.0.2", "10.0.0.3", "10.0.0.4"]
        },
      ]
    }
  }
}
```

### Example 8 - Vrack Services and BGP option ordered at once by Terraform

```terraform
locals {
  region = "eu-west-lim"
  vrack_name = "pn-000000"
}

data "ovh_me" "myaccount" {}

resource "ovh_vrackservices" "my-vrackservices" {
  ovh_subsidiary = data.ovh_me.myaccount.ovh_subsidiary
  plan = [
    {
      plan_code = "vrack-services"
      duration = "P1M"
      pricing_mode = "default"

      configuration = [
        {
          label = "region_name"
          value = local.region
        }
      ]
    }
  ]
  plan_option = [
    {
      plan_code    = "bgp-beta-vrack-services"
      duration     = "P1M"
      pricing_mode = "default"
      quantity     = 1
    }
  ]
  target_spec = {
    subnets = [
      {
        az   = "eu-west-lim-a"
        cidr = "10.0.0.0/24"
      },
    ]
    bgp = {
      bgp_range = [
        {
          cidr                 = "10.0.0.0/28"
          subnet_cidr          = "10.0.0.0/24"
          client_peers         = ["10.0.0.5"]
          infrastructure_peers = ["10.0.0.1", "10.0.0.2", "10.0.0.3", "10.0.0.4"]
        },
      ]
    }
  }
}

resource "ovh_vrack_vrackservices" "vrack-vrackservices-binding" {
  service_name   = local.vrack_name
  vrack_services = ovh_vrackservices.my-vrackservices.id
}
```

### Example 9 - IPv4 allowed prefix added to the BGP of a Vrack Services created by Terraform with its BGP option

Starting from the Vrack Services of example 8, add `allowed_prefixes` to `target_spec.bgp` to announce a public IP block through BGP. The block must be an OVHcloud IP block of your account: eligible blocks are listed by `GET /v2/vrackServices/resource/{vrackServicesId}/eligibleBgpAllowedPrefixes`. IPv4 prefix length must be between /24 and /30.

```terraform
locals {
  region = "eu-west-lim"
  vrack_name = "pn-000000"
}

data "ovh_me" "myaccount" {}

resource "ovh_vrackservices" "my-vrackservices" {
  ovh_subsidiary = data.ovh_me.myaccount.ovh_subsidiary
  plan = [
    {
      plan_code = "vrack-services"
      duration = "P1M"
      pricing_mode = "default"

      configuration = [
        {
          label = "region_name"
          value = local.region
        }
      ]
    }
  ]
  plan_option = [
    {
      plan_code    = "bgp-beta-vrack-services"
      duration     = "P1M"
      pricing_mode = "default"
      quantity     = 1
    }
  ]
  target_spec = {
    subnets = [
      {
        az   = "eu-west-lim-a"
        cidr = "10.0.0.0/24"
      },
    ]
    bgp = {
      allowed_prefixes = {
        ipv4 = ["x.x.x.x"]
      }
      bgp_range = [
        {
          cidr                 = "10.0.0.0/28"
          subnet_cidr          = "10.0.0.0/24"
          client_peers         = ["10.0.0.5"]
          infrastructure_peers = ["10.0.0.1", "10.0.0.2", "10.0.0.3", "10.0.0.4"]
        },
      ]
    }
  }
}

resource "ovh_vrack_vrackservices" "vrack-vrackservices-binding" {
  service_name   = local.vrack_name
  vrack_services = ovh_vrackservices.my-vrackservices.id
}
```

## Schema

### Required

- `target_spec` (Attributes) Target specification of the vRack Services (see [below for nested schema](#nestedatt--target_spec))

### Optional

- `ovh_subsidiary` (String) OVH subsidiaries
- `plan` (Attributes List) (see [below for nested schema](#nestedatt--plan))
- `plan_option` (Attributes List) (see [below for nested schema](#nestedatt--plan_option))

### Read-Only

- `checksum` (String) Computed hash used to control concurrent modification requests. Here, it represents the target specification value the request is based on
- `created_at` (String) Date of the vRack Services delivery
- `current_state` (Attributes) Current configuration applied to the vRack Services (see [below for nested schema](#nestedatt--current_state))
- `current_tasks` (Attributes List) Asynchronous operations ongoing on the vRack Services (see [below for nested schema](#nestedatt--current_tasks))
- `iam` (Attributes) IAM resource metadata (see [below for nested schema](#nestedatt--iam))
- `id` (String) Unique identifier
- `order` (Attributes) Details about an Order (see [below for nested schema](#nestedatt--order))
- `resource_status` (String) Reflects the readiness of the vRack Services. A new target specification request will be accepted only in `READY` status
- `updated_at` (String) Date of the Last vRack Services update

<a id="nestedatt--target_spec"></a>
### Nested Schema for `target_spec`

Required:

- `subnets` (Attributes List) Target specification of the subnets (see [below for nested schema](#nestedatt--target_spec--subnets))

Optional:

- `bgp` (Attributes) Target specification of the BGP configuration. Can only be set once the BGP option is subscribed (see `plan_option`). When omitted, the current BGP configuration is kept, but a non-empty `bgp_range` must stay in the configuration to avoid a change on every plan (see [below for nested schema](#nestedatt--target_spec--bgp))
- `services` (Attributes) Target specification of the services (see [below for nested schema](#nestedatt--target_spec--services))

<a id="nestedatt--target_spec--subnets"></a>
### Nested Schema for `target_spec.subnets`

Required:

- `az` (String) Availability zone of the subnet
- `cidr` (String) IP address range of the subnet in CIDR format. Must be a private network address (RFC1918)

Optional:

- `display_name` (String) Display name of the subnet. Format must follow `^[ a-zA-Z0-9-_.]{0,40}$`
- `vlan` (Number) Unique inner VLAN that allows subnets segregation. Authorized values: [2 - 4094] and `null` (untagged traffic). Must be `null` when the subnet is referenced by a BGP range


<a id="nestedatt--target_spec--bgp"></a>
### Nested Schema for `target_spec.bgp`

Optional:

- `allowed_prefixes` (Attributes) Allowed prefixes, OVHcloud public IP blocks announced through BGP (see [below for nested schema](#nestedatt--target_spec--bgp--allowed_prefixes))
- `bgp_range` (Attributes List) Target specification of the BGP ranges. At most one BGP range per subnet (see [below for nested schema](#nestedatt--target_spec--bgp--bgp_range))
- `client_asn` (Number) Private client ASN. Authorized values: [64512 - 65533]. Default value: 64512
- `infrastructure_asn` (Number) Private infrastructure ASN. Authorized values: [64513 - 65534]. Default value: 65534

<a id="nestedatt--target_spec--bgp--allowed_prefixes"></a>
### Nested Schema for `target_spec.bgp.allowed_prefixes`

Optional:

- `ipv4` (List of String) Allowed IPv4 prefixes. Prefix length must be between /24 and /30
- `ipv6` (List of String) Allowed IPv6 prefixes. Prefix length must be /56


<a id="nestedatt--target_spec--bgp--bgp_range"></a>
### Nested Schema for `target_spec.bgp.bgp_range`

Required:

- `cidr` (String) BGP range CIDR, contained in the referenced subnet
- `client_peers` (List of String) Client peer IPs within the BGP range CIDR (1 to 4). Must be of the same IP version as the referenced subnet
- `infrastructure_peers` (List of String) Infrastructure peer IPs within the BGP range CIDR. Must be of the same IP version as the referenced subnet
- `subnet_cidr` (String) Referenced subnet CIDR. The referenced subnet must be untagged (`vlan` must be `null`)



<a id="nestedatt--target_spec--services"></a>
### Nested Schema for `target_spec.services`

Required:

- `service_range` (Attributes List) Service range target specifications (see [below for nested schema](#nestedatt--target_spec--services--service_range))

<a id="nestedatt--target_spec--services--service_range"></a>
### Nested Schema for `target_spec.services.service_range`

Required:

- `cidr` (String) Service range CIDR, contained in the referenced subnet
- `service_endpoints` (Attributes List) Service endpoint target specifications (see [below for nested schema](#nestedatt--target_spec--services--service_range--service_endpoints))
- `subnet_cidr` (String) Referenced subnet CIDR

<a id="nestedatt--target_spec--services--service_range--service_endpoints"></a>
### Nested Schema for `target_spec.services.service_range.service_endpoints`

Required:

- `managed_service_urn` (String) IAM Resource URN of the managed service. Managed service Region must match vRack Services Region. Compatible managed service types are listed by /reference/compatibleManagedServiceType call





<a id="nestedatt--plan"></a>
### Nested Schema for `plan`

Required:

- `duration` (String) Duration selected for the purchase of the product, in ISO 8601 format, e.g. `P1M` (one month) or `P1Y` (one year). It must come from the same price as `pricing_mode`, see [Choosing the plan, duration and pricing mode of an ordered product](../guides/order_plan).
- `plan_code` (String) Identifier of the option offer
- `pricing_mode` (String) Pricing mode selected for the purchase of the product, e.g. `default`. It must come from the same price as `duration`.

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

- `duration` (String) Duration selected for the purchase of the product, in ISO 8601 format, e.g. `P1M` (one month) or `P1Y` (one year). It must come from the same price as `pricing_mode`, see [Choosing the plan, duration and pricing mode of an ordered product](../guides/order_plan).
- `plan_code` (String) Identifier of the option offer
- `pricing_mode` (String) Pricing mode selected for the purchase of the product, e.g. `default`. It must come from the same price as `duration`.
- `quantity` (Number) Quantity of product desired

Optional:

- `configuration` (Attributes List) (see [below for nested schema](#nestedatt--plan_option--configuration))

<a id="nestedatt--plan_option--configuration"></a>
### Nested Schema for `plan_option.configuration`

Required:

- `label` (String) Label for your configuration item
- `value` (String) Value or resource URL on API.OVH.COM of your configuration item



<a id="nestedatt--current_state"></a>
### Nested Schema for `current_state`

Read-Only:

- `bgp` (Attributes) BGP configuration of the vRack Services. Null when the BGP option is not subscribed (see [below for nested schema](#nestedatt--current_state--bgp))
- `configuration_status` (Attributes) vRack Services configuration status (see [below for nested schema](#nestedatt--current_state--configuration_status))
- `region` (String) Region of the vRack Services. List of compatible regions can be retrieved from /reference/region
- `services` (Attributes) Services current state (see [below for nested schema](#nestedatt--current_state--services))
- `subnets` (Attributes List) Subnets of the current vRack Services (see [below for nested schema](#nestedatt--current_state--subnets))
- `vrack_id` (String) vRack associated to the vRack Services

<a id="nestedatt--current_state--bgp"></a>
### Nested Schema for `current_state.bgp`

Read-Only:

- `allowed_prefixes` (Attributes) Allowed prefixes, OVHcloud public IP blocks announced through BGP (see [below for nested schema](#nestedatt--current_state--bgp--allowed_prefixes))
- `bandwidth` (Number) BGP ingress bandwidth in Mbps
- `bgp_range` (Attributes List) BGP ranges (see [below for nested schema](#nestedatt--current_state--bgp--bgp_range))
- `client_asn` (Number) Client ASN
- `configuration_status` (Attributes) BGP configuration status (see [below for nested schema](#nestedatt--current_state--bgp--configuration_status))
- `iam` (Attributes) IAM resource metadata of the BGP configuration (see [below for nested schema](#nestedatt--current_state--bgp--iam))
- `id` (String) BGP identifier
- `infrastructure_asn` (Number) Infrastructure ASN

<a id="nestedatt--current_state--bgp--allowed_prefixes"></a>
### Nested Schema for `current_state.bgp.allowed_prefixes`

Read-Only:

- `ipv4` (List of String) Allowed IPv4 prefixes
- `ipv6` (List of String) Allowed IPv6 prefixes


<a id="nestedatt--current_state--bgp--bgp_range"></a>
### Nested Schema for `current_state.bgp.bgp_range`

Read-Only:

- `cidr` (String) BGP range CIDR
- `client_peers` (List of String) Client peer IPs within the BGP range CIDR
- `configuration_status` (Attributes) BGP range configuration status (see [below for nested schema](#nestedatt--current_state--bgp--bgp_range--configuration_status))
- `infrastructure_peers` (List of String) Infrastructure peer IPs within the BGP range CIDR
- `subnet_cidr` (String) Referenced subnet CIDR

<a id="nestedatt--current_state--bgp--bgp_range--configuration_status"></a>
### Nested Schema for `current_state.bgp.bgp_range.configuration_status`

Read-Only:

- `reason` (String) Reason explaining the BGP range configuration status
- `status` (String) BGP range configuration status (`DRAFT`, `READY`)



<a id="nestedatt--current_state--bgp--configuration_status"></a>
### Nested Schema for `current_state.bgp.configuration_status`

Read-Only:

- `reason` (String) Reason for the status. Omitted when status is `READY`
- `status` (String) Configuration status (`DRAFT`, `READY`, `SUSPENDED`)


<a id="nestedatt--current_state--bgp--iam"></a>
### Nested Schema for `current_state.bgp.iam`

Read-Only:

- `display_name` (String) Resource display name
- `id` (String) Unique identifier of the resource
- `state` (String) Resource state
- `urn` (String) Unique resource name used in policies



<a id="nestedatt--current_state--configuration_status"></a>
### Nested Schema for `current_state.configuration_status`

Read-Only:

- `reason` (String) Reason for the status. Omitted when status is `READY`
- `status` (String) Configuration status (`DRAFT`, `READY`, `SUSPENDED`)


<a id="nestedatt--current_state--services"></a>
### Nested Schema for `current_state.services`

Read-Only:

- `service_range` (Attributes List) Service ranges (see [below for nested schema](#nestedatt--current_state--services--service_range))

<a id="nestedatt--current_state--services--service_range"></a>
### Nested Schema for `current_state.services.service_range`

Read-Only:

- `cidr` (String) Service range CIDR
- `remaining_ips` (Number) Number of remaining IPs in the service range
- `reserved_ips` (Number) Number of service range IPs reserved by OVHcloud
- `service_endpoints` (Attributes List) Service endpoints (see [below for nested schema](#nestedatt--current_state--services--service_range--service_endpoints))
- `service_endpoints_configuration_status` (Attributes) Service endpoints configuration status (see [below for nested schema](#nestedatt--current_state--services--service_range--service_endpoints_configuration_status))
- `subnet_cidr` (String) Subnet CIDR
- `used_ips` (Number) Number of service range IPs assigned to the managed services

<a id="nestedatt--current_state--services--service_range--service_endpoints"></a>
### Nested Schema for `current_state.services.service_range.service_endpoints`

Read-Only:

- `endpoints` (Attributes List) Endpoints representing the IPs assigned to the managed services (see [below for nested schema](#nestedatt--current_state--services--service_range--service_endpoints--endpoints))
- `managed_service_urn` (String) IAM Resource URN of the managed service. Compatible managed service types are listed by /reference/compatibleManagedServiceType call.

<a id="nestedatt--current_state--services--service_range--service_endpoints--endpoints"></a>
### Nested Schema for `current_state.services.service_range.service_endpoints.endpoints`

Read-Only:

- `description` (String) IP description defined in the managed service
- `ip` (String) IP address assigned by OVHcloud



<a id="nestedatt--current_state--services--service_range--service_endpoints_configuration_status"></a>
### Nested Schema for `current_state.services.service_range.service_endpoints_configuration_status`

Read-Only:

- `reason` (String) Reason for the status
- `status` (String) Configuration status (`DRAFT`, `READY`)




<a id="nestedatt--current_state--subnets"></a>
### Nested Schema for `current_state.subnets`

Read-Only:

- `az` (String) Availability zone of the subnet
- `cidr` (String) IP address range of the subnet in CIDR format
- `display_name` (String) Display name of the subnet
- `vlan` (Number) Unique inner VLAN that allows subnets segregation



<a id="nestedatt--current_tasks"></a>
### Nested Schema for `current_tasks`

Read-Only:

- `id` (String) Identifier of the current task
- `link` (String) Link to the related resource
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


## Timeouts

```terraform
resource "ovh_vrackservices" "my-vrackservices" {
  # ...

  timeouts {
    create = "1h"
  }
}
```

* `create` - (Default 30m)

## Import

A VrackServices can be imported using the `id`. Using the following configuration:

```terraform
import {
  to = ovh_vrackservices.vrackservices
  id = "vrs-xxx-xxx-xxx-xxx"
}
```

You can then run:

```bash
$ terraform plan -generate-config-out=vrackservices.tf
$ terraform apply
```

The file `vrackservices.tf` will then contain the imported resource's configuration, that can be copied next to the `import` block above. See https://developer.hashicorp.com/terraform/language/import/generating-configuration for more details.
