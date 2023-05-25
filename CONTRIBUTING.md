Thanks for wanting to contribute to this project ❤️.

This project accepts contributions. In order to contribute, you should pay attention to a few things:

1. Your code must follow the coding style rules
2. Your code must be fully documented
3. Your code must have acceptance test
4. Your new resource need to be imported by the user
5. Please test your new resources, datasources and acceptance tests
6. GitHub Pull Requests

## Coding and documentation Style:

- Code must be formatted with `make fmt` command
- Name your resources and datasources according to the API endpoint
- Check your documentation through [Terraform Doc Preview Tool](https://registry.terraform.io/tools/doc-preview)
- When adding a documentation page, use the `subcategory:` tag in the [YAML Frontmatter](https://developer.hashicorp.com/terraform/registry/providers/docs#yaml-frontmatter) with a value equals to the product name defined in the OVHcloud [product map](https://www.product-map.ovh/)

## Acceptance tests:

- Each resource and/or datasource need to have an acceptance test
- If you use new environment variables, document them in `website/docs/index.html.markdown`
- Acceptance tests must be run and must pass
- Don't forget to add or modify existing sweeper method if you think the acceptance tests may leave oprhan resources on failure

## Submitting Modifications:

The contributions should be submitted through new GitHub Pull Requests.

## Submiting an Issue:

In addition to contributions, we welcome [bug reports](https://github.com/ovh/terraform-provider-ovh/issues/new?assignees=&labels=&template=Bug_Report.md) and [feature requests](https://github.com/ovh/terraform-provider-ovh/issues/new?assignees=&labels=enhancement&template=Feature_Request.md).
