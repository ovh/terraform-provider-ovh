# OVH Terraform Provider examples

These examples are used in the [OVH provider documentation](https://registry.terraform.io/providers/ovh/ovh/latest/docs) to create and work with [OVHcloud](https://www.ovhcloud.com/en/) services.

## Provider configuration

The examples at the root of this directory are shown on the main page of the provider documentation:

1. [OAuth2](./example_1.tf): configure the provider with a `client_id` and a `client_secret`.
2. [Application Key/Application Secret](./example_2.tf): configure the provider with an `application_key`, an `application_secret` and a `consumer_key`.
3. [Example Usage](./example_3.tf): create an OVHcloud Managed Kubernetes cluster and a node pool.
4. [Using a locally built terraform-provider-ovh](./example_4.tf): use a provider built from this repository.

## Resources and data sources

Each resource and data source has its own directory, with the examples shown on its documentation page:

- [`resources/`](./resources): one directory per resource, named after the resource without the `ovh_` prefix. For example, [`resources/cloud_project_kube`](./resources/cloud_project_kube) holds the examples of `ovh_cloud_project_kube`.
- [`data-sources/`](./data-sources): one directory per data source, with the same naming.

## Contributing

New examples go in these directories, and are included in the documentation from the `templates/` directory. See [CONTRIBUTING.md](../CONTRIBUTING.md).
