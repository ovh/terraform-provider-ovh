Terraform OVH Provider
==================

The OVH Provider allows Terraform to manage [OVH](https://www.ovhcloud.com/) resources.

- Website: [registry.terraform.io](https://registry.terraform.io/providers/ovh/ovh/latest/docs)

Requirements
------------

- [Terraform](https://www.terraform.io/downloads.html) 0.12.x
- [Go](https://golang.org/doc/install) 1.20 (to build the provider plugin)

Building The Provider
---------------------

Clone repository to: `$GOPATH/src/github.com/ovh/terraform-provider-ovh`

```sh
$ mkdir -p $GOPATH/src/github.com/terraform-providers/; cd $GOPATH/src/github.com/terraform-providers/
$ git clone git@github.com:ovh/terraform-provider-ovh
```

Enter the provider directory and build the provider

```sh
$ cd $GOPATH/src/github.com/ovh/terraform-provider-ovh
$ make build
```

Using the provider
----------------------

Please see the documentation in the [Terraform registry](https://www.terraform.io/docs/providers/ovh/index.html).

Or you can browse the documentation within this repo [here](https://github.com/ovh/terraform-provider-ovh/tree/master/website/docs).

Using the locally built provider
----------------------

If you wish to test the provider from the local version you just built, you can try the following method.

First install the Terraform Provider binary into your local plugin repository:

```sh
# Set your target environment (OS_architecture): linux_amd64, darwin_amd64...
$ export ENV="linux_amd64"
$ make build
$ mkdir -p ~/.terraform.d/plugins/terraform.local/local/ovh/0.0.1/$ENV
$ cp $GOPATH/bin/terraform-provider-ovh ~/.terraform.d/plugins/terraform.local/local/ovh/0.0.1/$ENV/terraform-provider-ovh_v0.0.1
```

Then create a Terraform configuration using this exact provider:

```sh
$ mkdir ~/test-terraform-provider-ovh
$ cd ~/test-terraform-provider-ovh
$ cat > main.tf <<EOF
# Configure the OVHcloud Provider
terraform {
  required_providers {
    ovh = {
      source = "terraform.local/local/ovh"
      version = "0.0.1"
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
$ rm .terraform.lock.hcl && terraform init
...

# Apply your resources & datasources
$ terraform apply
...
```


Developing the Provider
---------------------------

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (version 1.8+ is *required*). You'll also need to correctly setup a [GOPATH](http://golang.org/doc/code.html#GOPATH), as well as adding `$GOPATH/bin` to your `$PATH`.

To compile the provider, run `make build`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

```sh
$ make build
...
$ $GOPATH/bin/terraform-provider-ovh
...
```

Testing the Provider
--------------------

In order to test the provider, you can simply run `make test`.

```sh
$ make test
```

In order to run the full suite of Acceptance tests you will need to have the following list of OVH products attached to your account:

- a [Vrack](https://www.ovh.ie/solutions/vrack/)
- a [Load Balancer](https://www.ovh.ie/solutions/load-balancer/)
- a registered [Domain](https://www.ovh.ie/domains/)
- a [Cloud Project](https://www.ovh.ie/public-cloud/instances/)

You will also need to setup your [OVH API](https://api.ovh.com) credentials. (see [documentation](https://www.terraform.io/docs/providers/ovh/index.html#configuration-reference))

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

In order for all the tests to pass you can run:

```sh
export OVH_IP_TEST="..."
export OVH_IP_BLOCK_TEST="..."
export OVH_IP_REVERSE_TEST="..."
export OVH_IP_MOVE_SERVICE_NAME_TEST="..."
export OVH_IPLB_SERVICE_TEST="..."
export OVH_CLOUD_PROJECT_SERVICE_TEST="..."
export OVH_CLOUD_PROJECT_FAILOVER_IP_TEST="..."
export OVH_CLOUD_PROJECT_FAILOVER_IP_ROUTED_TO_1_TEST="..."
export OVH_CLOUD_PROJECT_FAILOVER_IP_ROUTED_TO_2_TEST="..."
export OVH_VRACK_SERVICE_TEST="..."
export OVH_ZONE_TEST="..."

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
