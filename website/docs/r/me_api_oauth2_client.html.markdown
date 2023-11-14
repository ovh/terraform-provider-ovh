---
subcategory : "Account Management"
---

# ovh_me_api_oauth2_client

Creates an OAuth2 service account.

## Example Usage

An OAuth2 client for an app hosted at `my-app.com`, that uses the authorization code flow to authenticate.

```hcl
resource "ovh_me_api_oauth2_client" "my_oauth2_client_auth_code" {
  name = "OAuth2 authorization code service account"
  flow = "AUTHORIZATION_CODE"
  description = "An OAuth2 client using the authorization code flow for my-app.com"
  callback_urls = ["https://my-app.com/callback"]
}
```

An OAuth2 client for an app hosted at `my-app.com`, that uses the client credentials flow to authenticate.

```hcl
resource "ovh_me_api_oauth2_client" "my_oauth2_client_client_creds" {
  name = "client credentials service account"
  description = "An OAuth2 client using the client credentials flow for my app"
  flow = "CLIENT_CREDENTIALS"
}
```

## Argument Reference

* `name` - OAuth2 client name.
* `description` - OAuth2 client description.
* `flow` - The OAuth2 flow to use. `AUTHORIZATION_CODE` or `CLIENT_CREDENTIALS` are supported at the moment.
* `callback_urls` - List of callback urls when configuring the `AUTHORIZATION_CODE` flow.

## Attributes Reference

* `client_id` - Client ID of the created service account.
* `client_secret` - Client secret of the created service account.
* `name` - OAuth2 client name.
* `description` - OAuth2 client description.
* `flow` - The OAuth2 flow to use. `AUTHORIZATION_CODE` or `CLIENT_CREDENTIALS` are supported at the moment.
* `callback_urls` - List of callback urls when configuring the `AUTHORIZATION_CODE` flow.


## Import

OAuth2 clients can be imported using their `client_id`:

```bash
$ terraform import ovh_me_api_oauth2_client.my_oauth2_client client_id
```

Because the client_secret is only available for resources created using terraform, OAuth2 clients can also be imported using a `client_id` and a `client_secret` with a pipe separator:

```bash
$ terraform import ovh_me_api_oauth2_client.my_oauth2_client 'client_id|client_secret'
```
