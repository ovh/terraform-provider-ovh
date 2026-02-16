---
subcategory: "Vrack Services"
---

# ovh_vrackservices (Resource)

Orders a Vrack Services.

## Important

-> **NOTE** To order a product through Terraform, your account needs to have a default payment method defined. This can be done in the [OVHcloud Control Panel](https://www.ovh.com/manager/#/dedicated/billing/payment/method) or via API with the [/me/payment/method](https://api.ovh.com/console/#/me/payment/method~GET) endpoint.

~> **WARNING** `BANK_ACCOUNT` is not supported anymore, please update your default payment_method to `SEPA_DIRECT_DEBIT`

## Example Usage

### Example 1 - Simple Vrack Services order

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

### Example 2 - Vrack Services basic configuration

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
        cidr         = "192.168.0.0/24"
        service_range = {
          cidr = "192.168.0.0/29"
        }
        service_endpoints = []
      },
    ]
  }
}
```

### Example 3 - Vrack Services associated to a vRack

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

### Example 4 - Vrack Services configuration with a managed service

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
        cidr         = "192.168.0.0/24"
        display_name = "my.subnet"
        service_range = {
          cidr = "192.168.0.0/29"
        }
        vlan = 30
        service_endpoints = [
          {
            managed_service_urn = data.ovh_storage_efs.my-efs.iam.urn
          }
        ]
      },
    ]
  }
}
```

### Example 5 - Vrack Services complete configuration

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
        cidr         = "192.168.0.0/24"
        display_name = "my.subnet"
        service_range = {
          cidr = "192.168.0.0/29"
        }
        vlan = 30
        service_endpoints = [
          {
            managed_service_urn = data.ovh_storage_efs.my-efs.iam.urn
          }
        ]
      },
    ]
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

- `subnets` (Attributes List) Target specification of the subnets. Maximum one subnet per vRack Services (see [below for nested schema](#nestedatt--target_spec--subnets))

<a id="nestedatt--target_spec--subnets"></a>
### Nested Schema for `target_spec.subnets`

Required:

- `cidr` (String) IPv4 CIDR notation (e.g., 192.0.2.0/24)
- `service_endpoints` (Attributes List) Target specification of the Service Endpoints (see [below for nested schema](#nestedatt--target_spec--subnets--service_endpoints))
- `service_range` (Attributes) Target specification of the range dedicated to the subnet's services (see [below for nested schema](#nestedatt--target_spec--subnets--service_range))

Optional:

- `display_name` (String) Display name of the subnet. Format must follow `^[a-zA-Z0-9-_.]{0,40}$`
- `vlan` (Number) Unique inner VLAN that allows subnets segregation. Authorized values: [2 - 4094] and `null` (untagged traffic)

<a id="nestedatt--target_spec--subnets--service_endpoints"></a>
### Nested Schema for `target_spec.subnets.service_endpoints`

Required:

- `managed_service_urn` (String) IAM Resource URN of the managed service. Managed service Region must match vRack Services Region. Compatible managed service types are listed by /reference/compatibleManagedServiceType call


<a id="nestedatt--target_spec--subnets--service_range"></a>
### Nested Schema for `target_spec.subnets.service_range`

Required:

- `cidr` (String) IPv4 CIDR notation (e.g., 192.0.2.0/24)




<a id="nestedatt--plan"></a>
### Nested Schema for `plan`

Required:

- `duration` (String) Duration selected for the purchase of the product
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



<a id="nestedatt--current_state"></a>
### Nested Schema for `current_state`

Read-Only:

- `product_status` (String) Product status of the vRack Services
- `region` (String) Region of the vRack Services. List of compatible regions can be retrieved from /reference/region
- `subnets` (Attributes List) Subnets of the current vRack Services (see [below for nested schema](#nestedatt--current_state--subnets))

<a id="nestedatt--current_state--subnets"></a>
### Nested Schema for `current_state.subnets`

Read-Only:

- `cidr` (String) IP address range of the subnet in CIDR format
- `display_name` (String) Display name of the subnet
- `service_endpoints` (Attributes List) Service endpoints of the subnet (see [below for nested schema](#nestedatt--current_state--subnets--service_endpoints))
- `service_range` (Attributes) Defines a smaller subnet dedicated to the managed services IPs (see [below for nested schema](#nestedatt--current_state--subnets--service_range))
- `vlan` (Number) Unique inner VLAN that allows subnets segregation

<a id="nestedatt--current_state--subnets--service_endpoints"></a>
### Nested Schema for `current_state.subnets.service_endpoints`

Read-Only:

- `endpoints` (Attributes List) Endpoints representing the IPs assigned to the managed services (see [below for nested schema](#nestedatt--current_state--subnets--service_endpoints--endpoints))
- `managed_service_urn` (String) IAM Resource URN of the managed service. Compatible managed service types are listed by /reference/compatibleManagedServiceType call.

<a id="nestedatt--current_state--subnets--service_endpoints--endpoints"></a>
### Nested Schema for `current_state.subnets.service_endpoints.endpoints`

Read-Only:

- `description` (String) IP description defined in the managed service
- `ip` (String) IP address assigned by OVHcloud



<a id="nestedatt--current_state--subnets--service_range"></a>
### Nested Schema for `current_state.subnets.service_range`

Read-Only:

- `cidr` (String) CIDR dedicated to the subnet's services
- `remaining_ips` (Number) Number of remaining IPs in the service range
- `reserved_ips` (Number) Number of service range IPs reserved by OVHcloud
- `used_ips` (Number) Number of service range IPs assigned to the managed services




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
