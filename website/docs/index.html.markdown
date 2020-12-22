---
layout: "ovh"
page_title: "Provider: OVH"
sidebar_current: "docs-ovh-index"
description: |-
  The OVH provider is used to interact with the many resources supported by OVH. The provider needs to be configured with the proper credentials before it can be used.
---

# OVH Provider

The OVH provider is used to interact with the many resources supported by OVH. 
The provider needs to be configured with the proper credentials before it can be used.

Use the navigation to the left to read about the available resources.

## Configuration of the provider

Requests to OVH APIs require a set of secrets keys and the definition of the API end point. 
See [First Steps with the API](https://docs.ovh.com/gb/en/customer/first-steps-with-ovh-api/) (or the French version, [Premiers pas avec les API OVHcloud](https://docs.ovh.com/fr/api/api-premiers-pas/)) for a detailed explanation.

Besides the API end-point, the required keys are the `application_key`, the `application_secret`, and the `consumer_key`.
These keys can be generated via the [OVH token generation page](https://api.ovh.com/createToken/?GET=/*&POST=/*&PUT=/*&DELETE=/*). 

These parameters can be configured directly in the provider block as shown hereafter.


```hcl
# Configure the OVH Provider
provider "ovh" {
  endpoint           = "ovh-eu"
  application_key    = "yyyyyy"
  application_secret = "xxxxxxxxxxxxxx"
  consumer_key       = "zzzzzzzzzzzzzz"
}
```

Alternatively the secret keys can be retrieved from your environment.

 * `OVH_APPLICATION_KEY`
 * `OVH_APPLICATION_SECRET`
 * `OVH_CONSUMER_KEY` 

This later method (or a similar alternative) is recommended to avoid storing secret data in a source repository.


## Example Usage

```hcl
# Create a public cloud user
resource "ovh_cloud_project_user" "user-test" {
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

* `OVH_IPLB_SERVICE_TEST` - The ID of the IP Load Balancer to use

* `OVH_VRACK_SERVICE_TEST` - The ID of the vRack to use.

* `OVH_CLOUD_PROJECT_SERVICE_TEST` - The ID of your public cloud project.

* `OVH_ZONE_TEST` - The domain you own to test the domain_zone resource.

* `OVH_IP_TEST`, `OVH_IP_BLOCK_TEST`, `OVH_IP_REVERSE_TEST` - The values you have to set for testing ip reverse resources.

You will also need to [generate an OVH token](https://api.ovh.com/createToken/?GET=/*&POST=/*&PUT=/*&DELETE=/*) and use it to set the following environment variables:

 * `OVH_APPLICATION_KEY`

 * `OVH_APPLICATION_SECRET`

 * `OVH_CONSUMER_KEY`

You should be able to use any OVH environment to develop on as long as the above environment variables are set.
