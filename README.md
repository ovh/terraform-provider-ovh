Terraform OVH Provider
==================

- Website: https://www.terraform.io
- [![Gitter chat](https://badges.gitter.im/hashicorp-terraform/Lobby.png)](https://gitter.im/hashicorp-terraform/Lobby)
- Mailing list: [Google Groups](http://groups.google.com/group/terraform-tool)

Requirements
------------

- [Terraform](https://www.terraform.io/downloads.html) 0.12.x
- [Go](https://golang.org/doc/install) 1.13 (to build the provider plugin)

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

Please see the documentation at [terraform.io](https://www.terraform.io/docs/providers/ovh/index.html).

Or you can browse the documentation within this repo [here](https://github.com/ovh/terraform-provider-ovh/tree/master/website/docs).

Using the locally built provider
----------------------

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

```sh
$ mkdir ~/test-terraform-provider-ovh
$ cd ~/test-terraform-provider-ovh
$ cat > main.tf <<EOF
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
EOF
$ export OVH_ENDPOINT="..."
$ export OVH_APPLICATION_KEY="..."
$ export OVH_APPLICATION_SECRET="..."
$ export OVH_CONSUMER_KEY="..."
$ terraform init
...
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
- a [cloud project](https://www.ovh.ie/public-cloud/instances/)

You will also need to setup your [OVH api](https://api.ovh.com) credentials. (see [documentation](https://www.terraform.io/docs/providers/ovh/index.html#configuration-reference))

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
export OVH_IPLB_SERVICE_TEST="..."
export OVH_CLOUD_PROJECT_SERVICE_TEST="..."
export OVH_CLOUD_PROJECT_FAILOVER_IP_TEST="..."
export OVH_CLOUD_PROJECT_FAILOVER_IP_ROUTED_TO_1_TEST="..."
export OVH_CLOUD_PROJECT_FAILOVER_IP_ROUTED_TO_2_TEST="..."
export OVH_VRACK_SERVICE_TEST="..."
export OVH_ZONE_TEST="..."

$ make testacc
```

To filter acceptance test, you can run:

```sh
$ make testacc TESTARGS="-run TestAccCloudProjectPrivateNetwork"
```

To remove dangling resources, you can run:

```sh
$ make testacc TESTARGS="-sweep"
```
