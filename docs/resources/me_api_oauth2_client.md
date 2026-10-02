---
subcategory : "Account Management (IAM)"
---

# ovh_me_api_oauth2_client

Creates an OAuth2 service account.

## Example Usage

An OAuth2 client for an app hosted at `my-app.com`, that uses the authorization code flow to authenticate.

```terraform
resource "ovh_me_api_oauth2_client" "my_oauth2_client_auth_code" {
  name = "OAuth2 authorization code service account"
  flow = "AUTHORIZATION_CODE"
  description = "An OAuth2 client using the authorization code flow for my-app.com"
  callback_urls = ["https://my-app.com/callback"]
}
```

An OAuth2 client for an app hosted at `my-app.com`, that uses the client credentials flow to authenticate.

```terraform
resource "ovh_me_api_oauth2_client" "my_oauth2_client_client_creds" {
  name = "client credentials service account"
  description = "An OAuth2 client using the client credentials flow for my app"
  flow = "CLIENT_CREDENTIALS"
}
```

An OAuth2 client whose client secret is not kept in the Terraform state, for example a service account that only accesses a KMS with a certificate (see `ovh_okms_credential`).

```terraform
resource "ovh_me_api_oauth2_client" "my_oauth2_client_no_secret" {
  name                  = "service account without client secret"
  description           = "An OAuth2 client whose secret is not kept in the Terraform state"
  flow                  = "CLIENT_CREDENTIALS"
  discard_client_secret = true
}
```

## Argument Reference

* `name` - OAuth2 client name.
* `description` - OAuth2 client description.
* `flow` - The OAuth2 flow to use. `AUTHORIZATION_CODE` or `CLIENT_CREDENTIALS` are supported at the moment.
* `callback_urls` - List of callback urls when configuring the `AUTHORIZATION_CODE` flow.
* `discard_client_secret` - (Optional) If `true`, the client secret is not kept in the Terraform state and cannot be retrieved afterwards. Defaults to `false`. Switching from `true` to `false` creates a new OAuth2 client to get a new client secret.

## Attributes Reference

* `client_id` - Client ID of the created service account.
* `client_secret` - Client secret of the created service account. Empty when `discard_client_secret` is `true`.
* `name` - OAuth2 client name.
* `description` - OAuth2 client description.
* `flow` - The OAuth2 flow to use. `AUTHORIZATION_CODE` or `CLIENT_CREDENTIALS` are supported at the moment.
* `callback_urls` - List of callback urls when configuring the `AUTHORIZATION_CODE` flow.
* `discard_client_secret` - Whether the client secret is discarded from the Terraform state.
* `identity` - Identity URN of the service account to be used inside an IAM policy.

## Import

OAuth2 clients can be imported using their `client_id`:

```bash
$ terraform import ovh_me_api_oauth2_client.my_oauth2_client client_id
```

Because the client_secret is only available for resources created using terraform, OAuth2 clients can also be imported using a `client_id` and a `client_secret` with a pipe separator:

```bash
$ terraform import ovh_me_api_oauth2_client.my_oauth2_client 'client_id|client_secret'
```

Imported OAuth2 clients have `discard_client_secret` set to `false`. If your configuration sets it to `true`, the next `terraform apply` removes the client secret from the state without recreating the OAuth2 client.
