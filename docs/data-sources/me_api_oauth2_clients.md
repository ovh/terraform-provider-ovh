---
subcategory : "Account Management (IAM)"
---

# ovh_me_api_oauth2_client (Data Source)

Use this data source to retrieve information the list of existing OAuth2 service account IDs.

## Example Usage

```terraform
data "ovh_me_api_oauth2_client" "my_oauth2_clients" {
}
```

## Argument Reference

This datasource takes no arguments.

## Attributes Reference

* `client_ids` - The list of all the existing client IDs.
