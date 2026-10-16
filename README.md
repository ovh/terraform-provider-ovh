Terraform OVH Provider
==================

The OVH Provider allows Terraform to manage [OVH](https://www.ovhcloud.com/) resources.

[![License](https://img.shields.io/github/license/ovh/terraform-provider-ovh)](LICENSE)
[![GitHub release (latest by date)](https://img.shields.io/github/v/release/ovh/terraform-provider-ovh)](https://github.com/ovh/terraform-provider-ovh/releases/latest)
[![Release Date](https://img.shields.io/github/release-date/ovh/terraform-provider-ovh)](https://github.com/ovh/terraform-provider-ovh/releases)
![Terraform Provider Downloads](https://img.shields.io/terraform/provider/dm/786)
![Website](https://img.shields.io/badge/website-terraform?color=blue)

Requirements
------------

- [Terraform](https://developer.hashicorp.com/terraform/install) 1.0+
- [Go](https://go.dev/doc/install) 1.27 (to build the provider plugin)

Building The Provider
---------------------

Clone the repository anywhere on your machine, then build the provider:

```sh
$ git clone git@github.com:ovh/terraform-provider-ovh
$ cd terraform-provider-ovh
$ make build
```

`make build` runs `go install`, which puts the `terraform-provider-ovh` binary in `$(go env GOPATH)/bin` (or in `$GOBIN` if it is set).

Using the provider
----------------------

Please see the documentation in the [Terraform registry](https://registry.terraform.io/providers/ovh/ovh/latest/docs).

Or you can browse the documentation within this repo [here](https://github.com/ovh/terraform-provider-ovh/tree/master/docs).

Using the locally built provider
----------------------

Configure a dev override in your .terraformrc to be able to test the local version of the provider you built.
If not already present, create a terraformrc file and add a dev override for the ovh provider.

```sh
$ cat > ~/.terraformrc <<EOF
provider_installation {
  dev_overrides {
    "ovh/ovh" = "$(go env GOPATH)/bin/"
  }

  direct {}
}
EOF
```

Build the provider.

```sh
$ make build
```

Create a Terraform configuration with the following provider configuration.

```sh
$ mkdir ~/test-terraform-provider-ovh
$ cd ~/test-terraform-provider-ovh
$ cat > providers.tf <<EOF
# Configure the OVHcloud Provider
terraform {
  required_providers {
    ovh = {
      source = "ovh/ovh"
      version = ">= 2.7.0"
    }
  }
}

provider "ovh" {
}
EOF

# Export OVHcloud API credentials
$ export OVH_ENDPOINT="..."
$ export OVH_APPLICATION_KEY="..."
$ export OVH_APPLICATION_SECRET="..."
$ export OVH_CONSUMER_KEY="..."

# Initialize your project and remove existing dependencies lock file
$ rm -f .terraform.lock.hcl && terraform init
...

# Apply your resources & datasources
$ terraform apply
...
```


Developing the Provider
---------------------------

If you wish to work on the provider, you'll first need [Go](https://go.dev/doc/install) installed on your machine (version 1.27+ is *required*).

To compile the provider, run `make build`. This will build the provider and put the provider binary in the `$(go env GOPATH)/bin` directory.

```sh
$ make build
...
$ $(go env GOPATH)/bin/terraform-provider-ovh
...
```

Testing the Provider
--------------------

In order to test the provider, you can simply run `make test`.

```sh
$ make test
```

In order to run the full suite of Acceptance tests you will need to have the following list of OVH products attached to your account:

- a [Vrack](https://www.ovhcloud.com/en-ie/network/vrack/)
- a [Load Balancer](https://www.ovhcloud.com/en-ie/network/load-balancer/)
- a registered [Domain](https://www.ovhcloud.com/en-ie/domains/)
- a [Cloud Project](https://www.ovhcloud.com/en-ie/public-cloud/)
- a [KMS](https://www.ovhcloud.com/en-gb/identity-security-operations/key-management-service/)

You will also need to setup your [OVHcloud API](https://api.ovh.com) credentials (see [documentation](https://registry.terraform.io/providers/ovh/ovh/latest/docs#configuration-reference)).

Once setup, please follow these steps to prepare an environment for running the Acceptance tests:

```sh
$ cat > ~/.ovhrc <<EOF
# setup ovh api credentials
export OVH_ENDPOINT="ovh-eu"
export OVH_APPLICATION_KEY="..."
export OVH_APPLICATION_SECRET="..."
export OVH_CONSUMER_KEY="..."
EOF
$ source ~/.ovhrc
```

The acceptance tests read their test services and values from environment variables. The most common ones are listed below, and the [Testing and Development](./docs/index.md#testing-and-development) section of the provider documentation describes more of them:

```sh
export OVH_IP_TEST="..."
export OVH_IP_BLOCK_TEST="..."
export OVH_IP_REVERSE_TEST="..."
export OVH_IP_MOVE_SERVICE_NAME_TEST="..."
export OVH_IPLB_SERVICE_TEST="..."
export OVH_CLOUD_PROJECT_SERVICE_TEST="..."
export OVH_CLOUD_PROJECT_REGION_TEST="..."
export OVH_CLOUD_PROJECT_LOADBALANCER_TEST="..."
export OVH_CLOUD_PROJECT_FAILOVER_IP_TEST="..."
export OVH_CLOUD_PROJECT_FAILOVER_IP_ROUTED_TO_1_TEST="..."
export OVH_CLOUD_PROJECT_FAILOVER_IP_ROUTED_TO_2_TEST="..."
export OVH_VRACK_SERVICE_TEST="..."
export OVH_VRACK_PUBLIC_ROUTING_PRIORITY_REGION_TEST="..."
export OVH_VRACK_PUBLIC_ROUTING_PRIORITY_ID_TEST="..."
export OVH_VRACK_PUBLIC_ROUTING_PRIORITY_AZ_A_TEST="..."
export OVH_VRACK_PUBLIC_ROUTING_PRIORITY_AZ_B_TEST="..."
export OVH_VRACK_PUBLIC_ROUTING_PRIORITY_AZ_C_TEST="..."
export OVH_VRACKSERVICES_SERVICE_TEST="..."
export OVH_ZONE_TEST="..."
export OVH_DOMAIN_NS1_HOST_TEST="..."
export OVH_DOMAIN_NS1_IP_TEST="..."
export OVH_DOMAIN_NS2_HOST_TEST="..."
export OVH_DOMAIN_NS3_HOST_TEST="..."
export OVH_DOMAIN_DS_RECORD_ALGORITHM_TEST="..."
export OVH_DOMAIN_DS_RECORD_PUBLIC_KEY_TEST="..."
export OVH_DOMAIN_DS_RECORD_TAG_TEST="..."
export OVH_OKMS="..."
export OVH_STORAGE_EFS_SERVICE_TEST="..."
export OVH_DBAAS_LOGS_SERVICE_TEST="..."

$ make testacc
```

To run only one acceptance test, you can run:

```sh
$ make testacc TESTARGS="-run TestAccCloudProjectKubeUpdateVersion_basic"
```

To run one acceptance test and bypass go test caching:

```sh
$ TF_ACC=1 go test -count=1 $(go list ./... |grep -v 'vendor') -v -run  TestAccCloudProjectKubeUpdateVersion_basic -timeout 600m -p 10
```

To remove dangling resources, you can run:

```sh
$ make testacc TESTARGS="-sweep"
```

# Contributing

Please read the [contributing guide](./CONTRIBUTING.md) to learn about how you can contribute to the OVHcloud Terraform provider ;-).<br/>
There is no small contribution, don't hesitate!

Our awesome contributors:

<a href="https://github.com/ovh/terraform-provider-ovh/graphs/contributors">
  <img src="https://contrib.rocks/image?repo=ovh/terraform-provider-ovh" />
</a>
