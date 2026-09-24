---
page_title: "Consumer key access rules for resources that place orders"
---

# Consumer key access rules for resources that place orders

Some resources cannot be created through their product's own API. To create them, the provider places an order through the OVHcloud order workflow (cart, checkout, payment), and to delete them it goes through the termination workflow (terminate, then confirm with a token received by email). This is also the case for products that are free, such as `ovh_okms`, `ovh_vrack` or `ovh_vrackservices`.

If you authenticate with an application key and a consumer key whose access rules are restricted, the consumer key must allow these workflow routes in addition to the product's own routes. Otherwise, `terraform apply` or `terraform destroy` fails with `This call has not been granted`.

This page lists the routes used by each of these resources. See [#1448](https://github.com/ovh/terraform-provider-ovh/issues/1448) for the related discussion.

## Affected resources

| Resource | Order product (`POST /order/cart/{cartId}/{product}`) | Product API path | Terminate | Confirm termination |
|---|---|---|---|---|
| `ovh_okms` | `okms` | `/okms` | `POST /services/{serviceId}/terminate` | `POST /services/{serviceId}/terminate/confirm` |
| `ovh_vrack` | `vrack` | `/vrack` | `POST /vrack/{serviceName}/terminate` | `POST /vrack/{serviceName}/confirmTermination` |
| `ovh_vrackservices` | `vrackServices` | `/vrackServices` | `POST /services/{serviceId}/terminate` | `POST /services/{serviceId}/terminate/confirm` |
| `ovh_cloud_project` | `cloud` | `/cloud/project` | `POST /cloud/project/{serviceName}/terminate` | `POST /cloud/project/{serviceName}/confirmTermination` |
| `ovh_domain_zone` | `dns` | `/domain/zone` | `POST /domain/zone/{serviceName}/terminate` | `POST /domain/zone/{serviceName}/confirmTermination` |
| `ovh_domain_name` | `domain` | `/domain` | `POST /services/{serviceId}/terminate` | `POST /services/{serviceId}/terminate/confirm` |
| `ovh_ip_service` | `ip` | `/ip/service` | `POST /ip/service/{serviceName}/terminate` | `POST /ip/service/{serviceName}/confirmTermination` |
| `ovh_iploadbalancing` | `ipLoadbalancing` | `/ipLoadbalancing` | `POST /ipLoadbalancing/{serviceName}/terminate` | `POST /ipLoadbalancing/{serviceName}/confirmTermination` |
| `ovh_hosting_privatedatabase` | `privateSQL` | `/hosting/privateDatabase` | `POST /hosting/privateDatabase/{serviceName}/terminate` | `POST /hosting/privateDatabase/{serviceName}/confirmTermination` |
| `ovh_storage_efs` | `netapp` | `/storage/netapp` | `POST /storage/netapp/{serviceName}/terminate` | `POST /storage/netapp/{serviceName}/confirmTermination` |
| `ovh_vps` | `vps` | `/vps` | `POST /vps/{serviceName}/terminate` | `POST /vps/{serviceName}/confirmTermination` |
| `ovh_dedicated_server` ¹ | `baremetalServers` / `eco` | `/dedicated/server` | `POST /dedicated/server/{serviceName}/terminate` | `POST /dedicated/server/{serviceName}/confirmTermination` |

¹ Only when `service_name` is not set. When an existing server is given, no order is placed.

## Routes used on create

| Route | When |
|---|---|
| `GET /me` | Only if `ovh_subsidiary` is not set |
| `POST /order/cart` | Always |
| `POST /order/cart/{cartId}/assign` | Always |
| `POST /order/cart/{cartId}/{product}` | Always |
| `POST /order/cart/{cartId}/item/{itemId}/configuration` | If the plan has `configuration` entries |
| `POST /order/cart/{cartId}/{product}/options` | If `plan_option` is set |
| `GET /me/payment/method` | Always (looks up the default payment method) |
| `GET /order/cart/{cartId}/checkout` | Only if there is no default payment method (checks whether the order is free) |
| `POST /order/cart/{cartId}/checkout` | Always |
| `POST /me/order/{orderId}/pay` | Paid orders (default payment method) |
| `POST /me/order/{orderId}/payWithRegisteredPaymentMean` | Free orders (`fidelityAccount`) |
| `GET /me/order/{orderId}/status` | While waiting for delivery (all except `ovh_dedicated_server`) |

Once the order is placed, the provider reads it to find the created service:

| Route | When |
|---|---|
| `GET /me/order/{orderId}` | `ovh_cloud_project` only |
| `GET /me/order/{orderId}/details` | All except `ovh_domain_name` |
| `GET /me/order/{orderId}/details/{orderDetailId}` | All except `ovh_domain_name`, on `ovh-eu` / `ovh-ca`; always for `ovh_cloud_project` |
| `GET /me/order/{orderId}/details/{orderDetailId}/extension` | All except `ovh_domain_name`, on `ovh-eu` / `ovh-ca`; always for `ovh_cloud_project` |
| `GET /me/order/{orderId}/details/{orderDetailId}/operations` | All except `ovh_domain_name` and `ovh_cloud_project`, on `ovh-us` only |
| `GET /me/order/{orderId}/details/{orderDetailId}/operations/{operationId}` | All except `ovh_domain_name` and `ovh_cloud_project`, on `ovh-us` only |

`ovh_domain_name` does not read the order details: its service name is the domain itself.

## Routes used on destroy

| Route | When |
|---|---|
| `GET /me/notification/email/history` | Always (polls for up to 30 minutes for the termination email) |
| `GET /me/notification/email/history/{id}` | Always (the confirmation token is read from the email body) |
| `GET /services?resourceName={name}` | `ovh_okms`, `ovh_vrackservices` |
| `GET /services?resourceName={name}&routes=/domain/{serviceName}` | `ovh_domain_name` |
| Terminate and confirm termination routes | See the table above |

## Example access rules

Access rules are path prefixes that accept a `*` wildcard, and they do not include the API version: a rule on `/okms/*` also applies to `/v2/okms/...`, while a rule written as `/v2/okms/*` grants nothing.

Replace `<product-path>` with the product API path from the table above. The first block is used on create, the second on destroy, and the last one by the resource itself.

```json
{
  "accessRules": [
    { "method": "GET",    "path": "/me" },
    { "method": "GET",    "path": "/me/payment/method" },
    { "method": "POST",   "path": "/order/cart" },
    { "method": "GET",    "path": "/order/cart/*" },
    { "method": "POST",   "path": "/order/cart/*" },
    { "method": "GET",    "path": "/me/order/*" },
    { "method": "POST",   "path": "/me/order/*/pay" },
    { "method": "POST",   "path": "/me/order/*/payWithRegisteredPaymentMean" },

    { "method": "GET",    "path": "/me/notification/email/history" },
    { "method": "GET",    "path": "/me/notification/email/history/*" },
    { "method": "GET",    "path": "/services" },
    { "method": "POST",   "path": "/services/*/terminate" },
    { "method": "POST",   "path": "/services/*/terminate/confirm" },

    { "method": "GET",    "path": "<product-path>/*" },
    { "method": "POST",   "path": "<product-path>/*" },
    { "method": "PUT",    "path": "<product-path>/*" },
    { "method": "DELETE", "path": "<product-path>/*" }
  ]
}
```

Notes:

- `GET /me` can be left out if `ovh_subsidiary` is set in your configuration.
- The destroy block can be left out if you never delete these resources with Terraform. `/services/*` is only used by `ovh_okms`, `ovh_vrackservices` and `ovh_domain_name`; the other resources terminate through their product API path, which is already covered by the last block.
- A consumer key with restricted access rules is created with [`POST /auth/credential`](https://eu.api.ovh.com/console/?section=%2Fauth&branch=v1#post-/auth/credential) and then validated once in a browser.

~> **WARNING** The order and termination routes are account-wide and cannot be restricted to a single service. A consumer key with the rules above can order and pay for any product, read the account's orders and notification emails, and terminate services. Avoid storing it as a long-lived CI secret: use it for the runs that create or delete these resources, and use a key restricted to `<product-path>/*` for day-to-day runs.

## Creating the service once, then managing it with a restricted key

For services that are created once and rarely deleted, such as a Public Cloud project, you can avoid giving the order and termination routes to your CI credentials:

1. Create the service in the [OVHcloud Control Panel](https://www.ovh.com/manager/), or with a broader consumer key used only once.
2. Import it into your Terraform state, using an `import` block or `terraform import` (see the Import section of each resource page).
3. Use a consumer key restricted to the product API path, for example `/cloud/project/<project_id>/*`.
4. Prevent Terraform from terminating it: set `deletion_protection = true` on `ovh_cloud_project`, or add `lifecycle { prevent_destroy = true }` on the other resources.
