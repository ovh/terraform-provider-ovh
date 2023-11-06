---
subcategory : "Account Management"
---

# ovh_me_api_oauth2_client (Data Source)

Use this data source to retrieve information about an existing OAuth2 service account.

## Example Usage

```hcl
data "ovh_me_api_oauth2_client" "my_oauth2_client" {
  client_id = "5f8969a993ec8b4b"
}
```

## Argument Reference

* `client_id` - Client ID of an existing OAuth2 service account.
* `client_secret` - Optional: Client secret of an existing service account. Has to be declared if needed as it will not be imported from the API.

## Attributes Reference

* `client_id` - Client ID of the created service account.
* `name` - OAuth2 client name.
* `description` - OAuth2 client description.
* `flow` - The OAuth2 flow to use. `AUTHORIZATION_CODE` or `CLIENT_CREDENTIALS` are supported at the moment.
* `callback_urls` - List of callback urls when configuring the `AUTHORIZATION_CODE` flow.
