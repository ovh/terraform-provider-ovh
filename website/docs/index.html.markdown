---
layout: "ovh"
page_title: "Provider: OVH"
sidebar_current: "docs-ovh-index"
description: |-
  The OVH provider is used to interact with the many resources supported by OVHcloud. The provider needs to be configured with the proper credentials before it can be used.
---

# OVH Provider

The OVH provider is the entry point to interact with the resources provided by OVHcloud. 

-> __NOTE__ According on your needs, you may need to use additional providers. This [documentation page](https://help.ovhcloud.com/csm/en-gb-terraform-at-ovhcloud?id=kb_article_view&sysparm_article=KB0054612) provides the mapping between the control panel concepts and the terraform providers / ressources.

Use the navigation to the left to read about the available resources.

## Provider configuration

The provider needs to be configured with the proper credentials before it can be used. Requests to OVHcloud APIs require a set of secrets keys and the definition of the API end point. See [First Steps with the API](https://docs.ovh.com/gb/en/customer/first-steps-with-ovh-api/) (or the French version, [Premiers pas avec les API OVHcloud](https://docs.ovh.com/fr/api/api-premiers-pas/)) for a detailed explanation.

Besides the API end-point, the required keys are the `application_key`, the `application_secret`, and the `consumer_key`.
These keys can be generated via the [OVHcloud token generation page](https://api.ovh.com/createToken/?GET=/*&POST=/*&PUT=/*&DELETE=/*). 

These parameters can be configured directly in the provider block as shown hereafter.

Terraform 0.13 and later:

```hcl
terraform {
  required_providers {
    ovh = {
      source = "ovh/ovh"
    }
  }
}

provider "ovh" {
  endpoint           = "ovh-eu"
  application_key    = "xxxxxxxxx"
  application_secret = "yyyyyyyyy"
  consumer_key       = "zzzzzzzzzzzzzz"
}
```

Alternatively the secret keys can be retrieved from your environment.

 * `OVH_ENDPOINT`
 * `OVH_APPLICATION_KEY`
 * `OVH_APPLICATION_SECRET`
 * `OVH_CONSUMER_KEY` 

This later method (or a similar alternative) is recommended to avoid storing secret data in a source repository.


## Example Usage

```hcl
variable "service_name" {
  default = "wwwwwww"
}

# Create an OVHcloud Managed Kubernetes cluster
resource "ovh_cloud_project_kube" "my_kube_cluster" {
  service_name = var.service_name
  name         = "my-super-kube-cluster"
  region       = "GRA5"
  version      = "1.22"
}

# Create a Node Pool for our Kubernetes clusterx
resource "ovh_cloud_project_kube_nodepool" "node_pool" {
  service_name  = var.service_name
  kube_id       = ovh_cloud_project_kube.my_kube_cluster.id
  name          = "my-pool" //Warning: "_" char is not allowed!
  flavor_name   = "b2-7"
  desired_nodes = 3
  max_nodes     = 3
  min_nodes     = 3
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
* `OVH_IPLB_IPFO_TEST`- An array of FailOver IPs (also known as Additional IPs) that shall be associated with the IPLB Service

* `OVH_VRACK_SERVICE_TEST` - The ID of the vRack to use.

* `OVH_CLOUD_PROJECT_SERVICE_TEST` - The ID of your public cloud project.

* `OVH_CLOUD_PROJECT_DATABASE_ENGINE_TEST` - The name of the database engine to test.

* `OVH_CLOUD_PROJECT_DATABASE_VERSION_TEST` - The version of the database engine to test.

* `OVH_CLOUD_PROJECT_DATABASE_KAFKA_VERSION_TEST` - The version of the kafka to test. if not set `OVH_CLOUD_PROJECT_DATABASE_VERSION_TEST` is use.

* `OVH_CLOUD_PROJECT_DATABASE_M3DB_VERSION_TEST` - The version of the M3DB to test. if not set `OVH_CLOUD_PROJECT_DATABASE_VERSION_TEST` is use.

* `OVH_CLOUD_PROJECT_DATABASE_MONGODB_VERSION_TEST` - The version of the mongodb to test. if not set `OVH_CLOUD_PROJECT_DATABASE_VERSION_TEST` is use.

* `OVH_CLOUD_PROJECT_DATABASE_OPENSEARCH_VERSION_TEST` - The version of the opensearch to test. if not set `OVH_CLOUD_PROJECT_DATABASE_VERSION_TEST` is use.

* `OVH_CLOUD_PROJECT_DATABASE_POSTGRESQL_VERSION_TEST` - The version of the postgresql to test. if not set `OVH_CLOUD_PROJECT_DATABASE_VERSION_TEST` is use.

* `OVH_CLOUD_PROJECT_DATABASE_REDIS_VERSION_TEST` - The version of the redis to test. if not set `OVH_CLOUD_PROJECT_DATABASE_VERSION_TEST` is use.

* `OVH_CLOUD_PROJECT_DATABASE_REGION_TEST` - The region of the database service to test.

* `OVH_CLOUD_PROJECT_DATABASE_FLAVOR_TEST` - The node flavor of the database service to test.

* `OVH_CLOUD_PROJECT_DATABASE_IP_RESTRICTION_IP_TEST` - The IP restriction to test.

* `OVH_CLOUD_PROJECT_FAILOVER_IP_TEST` - The ip address of your public cloud failover ip.

* `OVH_CLOUD_PROJECT_FAILOVER_IP_ROUTED_TO_1_TEST` - The GUID of an instance to which failover IP addresses can be attached

* `OVH_CLOUD_PROJECT_FAILOVER_IP_ROUTED_TO_2_TEST` - The GUID of a secondary instance to which failover IP addresses can be attached. There must be 2 as associations can only be updated not removed. To test effectively, the failover ip address must be moved between instances 

* `OVH_CLOUD_PROJECT_KUBE_REGION_TEST` - The region of your public cloud kubernetes project.

* `OVH_CLOUD_PROJECT_KUBE_VERSION_TEST` - The version of your public cloud kubernetes project.

* `OVH_DEDICATED_SERVER` - The name of the dedicated server to test dedicated_server_networking resource.

* `OVH_NASHA_SERVICE_TEST` - The name of your HA-NAS service.

* `OVH_ZONE_TEST` - The domain you own to test the domain_zone resource.

* `OVH_IP_TEST`, `OVH_IP_BLOCK_TEST`, `OVH_IP_REVERSE_TEST` - The values you have to set for testing ip reverse resources.

* `OVH_DBAAS_LOGS_SERVICE_TEST` - The name of your Dbaas logs service.

* `OVH_DBAAS_LOGS_LOGSTASH_VERSION_TEST` - The name of your Dbaas logs Logstash version.

* `OVH_TESTACC_ORDER_VRACK` - set this variable to "yes" will order vracks.

* `OVH_TESTACC_ORDER_CLOUDPROJECT` - set this variable to "yes" will order cloud projects.

* `OVH_TESTACC_ORDER_DOMAIN` - set this variable to "mydomain.ovh" to run tests for domain zones.

* `OVH_HOSTING_PRIVATEDATABASE_SERVICE_TEST` - Set a Webhosting database service name to test

* `OVH_HOSTING_PRIVATEDATABASE_NAME_TEST` - The database name created of your service name

* `OVH_HOSTING_PRIVATEDATABASE_USER_TEST` - The username of your private database to test

* `OVH_HOSTING_PRIVATEDATABASE_PASSWORD_TEST` - The password of your private database's user to test

* `OVH_HOSTING_PRIVATEDATABASE_GRANT_TEST` - The grant of your private database's user to test

* `OVH_HOSTING_PRIVATEDATABASE_WHITELIST_IP_TEST` - Whitelist an IP address to connect to your instance

* `OVH_HOSTING_PRIVATEDATABASE_WHITELIST_NAME_TEST` - Set a custom label to your whitelisted IP address

* `OVH_HOSTING_PRIVATEDATABASE_WHITELIST_SERVICE_TEST` - Set this variable to `true` to authorize service access the whitelisted IP address

* `OVH_HOSTING_PRIVATEDATABASE_WHITELIST_SFTP_TEST` - Set this variable to `true` to authorize SFTP access to a whitelisted IP address

* `OVH_CLOUD_PROJECT_WORKFLOW_BACKUP_REGION_TEST` - The openstack region in which the workflow will be defined
* `OVH_CLOUD_PROJECT_WORKFLOW_BACKUP_INSTANCE_ID_TEST` - The openstack id of the instance to backup

### Used by OVHcloud internal account only:

* `OVH_TESTACC_ORDER_IPLOADBALANCING` - set this variable to "yes" will order ip loadbalancing.

* `OVH_TESTACC_IP` - set this variable to "yes" will order public ip blocks.


### Using a locally built terraform-provider-ovh

If you wish to test the provider from the local version you just built, you can try the following method.

First install the terraform provider binary into your local plugin repository:

```sh
# Set your target environment (OS_architecture): linux_amd64, darwin_amd64...
$ export ENV="linux_amd64"
$ make build
...
$ mkdir -p ~/.terraform.d/plugins/terraform.local/local/ovh/0.0.1/$ENV
$ cp $GOPATH/bin/terraform-provider-ovh ~/.terraform.d/plugins/terraform.local/local/ovh/0.0.1/$ENV/terraform-provider-ovh_v0.0.1
```

Then create a terraform configuration using this exact provider:

```hcl
terraform {
  required_providers {
    ovh = {
      source = "terraform.local/local/ovh"
      version = "0.0.1"
    }
  }
}

data "ovh_me" "me" {}

output "me" {
  value = data.ovh_me.me
}
```

This allows you to use your unreleased version of the provider.
The version number is not important and you can use whatever you like in this example but you need to stay coherent between the configuration, the directory structure and the binary filename.
