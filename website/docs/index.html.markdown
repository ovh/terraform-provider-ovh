---
layout: "ovh"
page_title: "Provider: OVH"
sidebar_current: "docs-ovh-index"
description: |-
  The OVH provider is used to interact with the many resources supported by OVH. The provider needs to be configured with the proper credentials before it can be used.
---

# OVH Provider

The OVH provider is used to interact with the
many resources supported by OVH. The provider needs to be configured
with the proper credentials before it can be used.

Use the navigation to the left to read about the available resources.

## Configuration of the provider

Requests to OVH APIs need to configure secrets keys in the provider, either fetching them from `~/.ovh.conf` file, in configuration of OVH provider or from your environment.

It is recommend to install [ovh-cli](https://github.com/ovh/ovh-cli) to handle and manage all your secret keys.

Follow [installation](https://github.com/ovh/ovh-cli#installation) then [setup](https://github.com/ovh/ovh-cli#getting-started) steps of ovh-cli to initialize your environment (secret keys and `~/.ovh.conf` file).

Then, you can just declare a minimal configuration of the OVH provider:

```hcl
# Configure the OVH Provider
provider "ovh" {
  endpoint = "ovh-eu"
}
```
Secret keys `endpoint`, `application_key`, `application_secret` or
`consumer_key` will be fetched from the `~/.ovh.conf` file.

Or you can declare them in provider configuration:

```hcl
# Configure the OVH Provider
provider "ovh" {
  endpoint           = "ovh-eu"
  application_key    = "yyyyyy"
  application_secret = "xxxxxxxxxxxxxx"
  consumer_key       = "zzzzzzzzzzzzzz"
}
```

Or let the provider fetching them from your environment (see "[Configuration reference](#configuration-reference)").


## Example Usage

```
# Create a public cloud user
resource "ovh_publiccloud_user" "user-test" {
  # ...
}
```

## Configuration Reference

The following arguments are supported:

* `endpoint` - (Required) Specify which API endpoint to use.
  It can be set using the `OVH_ENDPOINT` environment
  variable. e.g. `ovh-eu` or `ovh-ca`.

* `application_key` - (Optional) The API Application Key. If omitted,
  the `OVH_APPLICATION_KEY` environment variable is used.

* `application_secret` - (Optional) The API Application Secret. If omitted,
  the `OVH_APPLICATION_SECRET` environment variable is used.

* `consumer_key` - (Optional) The API Consumer key. If omitted,
  the `OVH_CONSUMER_KEY` environment variable is used.

## Testing and Development

In order to run the Acceptance Tests for development, the following environment
variables must also be set:

* `OVH_ENDPOINT` - possible value are: `ovh-eu`, `ovh-ca`, `ovh-us`, `soyoustart-eu`, `soyoustart-ca`, `kimsufi-ca`, `kimsufi-eu`, `runabove-ca`

* `OVH_IPLB_SERVICE` - The ID of the IP Load Balancer to use

* `OVH_VRACK` - The ID of the vRack to use.

* `OVH_PUBLIC_CLOUD` - The ID of your public cloud project.

* `OVH_ZONE` - The domain you own to test the domain_zone resource.

You will also need to [generate an OVH token](https://api.ovh.com/createToken/?GET=/*&POST=/*&PUT=/*&DELETE=/*) and use it to set the following environment variables:

 * `OVH_APPLICATION_KEY`

 * `OVH_APPLICATION_SECRET`

 * `OVH_CONSUMER_KEY`

You should be able to use any OVH environment to develop on as long as the above environment variables are set.
