# Deploy an OVHcloud Manager Kubernetes cluster, a Node Pool and an application

Through this procedure, you can [create multiple OVHcloud Managed Kubernetes, through Terraform](https://docs.ovh.com/gb/en/kubernetes/creating-a-cluster-through-terraform/).

You will create:

* in parallel, several Kubernetes clusters in OVH (depending on the number setted in `scripts/create.sh` script, by default is one cluster)
* a node pool with 3 nodes
* a functional deployment (an app)
* and a service of type Load Balancer

# Prerequisites

Generate [OVH API credentials](https://api.ovh.com/createToken/?GET=/*&POST=/*&PUT=/*&DELETE=/*) and then export in environment variables in your machine like this:

```
$ export OVH_ENDPOINT="ovh-eu"
$ export OVH_APPLICATION_KEY="xxx"
$ export OVH_APPLICATION_SECRET="xxxxx"
$ export OVH_CONSUMER_KEY="xxxxx"
```

Or you can directly put them in `provider.tf` in OVH provider definition:

```
provider "ovh" {
  version            = "~> 0.16"
  endpoint           = "ovh-eu"
  application_key    = "xxx"
  application_secret = "xxx"
  consumer_key       = "xxx"
}
```

Set in `variables.tf` your service_name parameter (Public Cloud project ID):

```
variable "service_name" {
  default = "xxxxx"
}
```

# How To

Create the Kubernetes clusters and for each, apply a deployment a service and when the OVHcloud Load Balancer is created, curl the app:

```
./scripts/create.sh
```

Output are writted in `logs` file.
Display the logs in realtime:

```sh
$ tail -f logs
```

# Clean

You can remove/destroy generated files and OVHcloud resources:

```
./scripts/clean.sh
```
