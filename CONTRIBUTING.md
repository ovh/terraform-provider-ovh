fdskjfds
# Contributing to `terraform-provider-ovh`

Thanks for wanting to contribute to this project ❤️.

This project accepts contributions. In order to contribute, you should pay attention to a few things:

1. Your code must follow the coding style rules
2. Your code must be fully documented
3. Your code must have acceptance test
4. Every Terraform resource added must be importable by the end-user
5. Your work must be signed (see below)
6. Please test your new resources, datasources and acceptance tests
7. Use GitHub Pull Requests to contribute

## Coding and documentation Style:

- Code must be formatted with `make fmt` command
- Name your resources and datasources according to the API endpoint
- The examples of resources and datasources in the documentation must follow the [Terraform style guidelines](https://developer.hashicorp.com/terraform/language/style)
- Check your documentation through [Terraform Doc Preview Tool](https://registry.terraform.io/tools/doc-preview)
- When adding a documentation page, use the `subcategory:` tag in the [YAML Frontmatter](https://developer.hashicorp.com/terraform/registry/providers/docs#yaml-frontmatter) with a value equals to the product name defined in the OVHcloud [product map](https://www.product-map.ovh/)
- New documentation pages should be added first in the directory `templates/`, with the examples being placed in the `examples/` directory. Once this is done, the content in `docs/` directory must be generated with [tfplugindocs](https://github.com/hashicorp/terraform-plugin-docs?tab=readme-ov-file#usage).

## Acceptance tests:

- Each resource and/or datasource need to have an acceptance test
- If you use new environment variables, document them in `website/docs/index.html.markdown`
- Acceptance tests must be run and must pass
- Don't forget to add or modify existing sweeper method if you think the acceptance tests may leave orphan resources on failure

## Submitting Modifications:

The contributions should be submitted through Github Pull Requests
and follow the DCO which is defined below.

## Licensing for new files

terraform-provider-ovh is licensed under the Mozilla Public License 2.0. Anything
contributed to terraform-provider-ovh must be released under this license.

## Submiting an Issue:

In addition to contributions, we welcome [bug reports](https://github.com/ovh/terraform-provider-ovh/issues/new?template=report-a-bug.md), [resource or datasource requests](https://github.com/ovh/terraform-provider-ovh/issues/new?template=request-a-new-resource-and-or-datasource.md), [documentation errors reports](https://github.com/ovh/terraform-provider-ovh/issues/new?template=report-a-documentation-error.md) and [feature requests](https://github.com/ovh/terraform-provider-ovh/issues/new?template=request-a-feature.md).


## Developer Certificate of Origin (DCO)

To improve tracking of contributions to this project we will use a
process modeled on the modified DCO 1.1 and use a "sign-off" procedure
on patches that are being emailed around or contributed in any other
way.

The sign-off is a simple line at the end of the explanation for the
patch, which certifies that you wrote it or otherwise have the right
to pass it on as an open-source patch.  The rules are pretty simple:
if you can certify the below:

By making a contribution to this project, I certify that:

(a) The contribution was created in whole or in part by me and I have
    the right to submit it under the open source license indicated in
    the file; or

(b) The contribution is based upon previous work that, to the best of
    my knowledge, is covered under an appropriate open source License
    and I have the right under that license to submit that work with
    modifications, whether created in whole or in part by me, under
    the same open source license (unless I am permitted to submit
    under a different license), as indicated in the file; or

(c) The contribution was provided directly to me by some other person
    who certified (a), (b) or (c) and I have not modified it.

(d) The contribution is made free of any other party's intellectual
    property claims or rights.

(e) I understand and agree that this project and the contribution are
    public and that a record of the contribution (including all
    personal information I submit with it, including my sign-off) is
    maintained indefinitely and may be redistributed consistent with
    this project or the open source license(s) involved.


then you just add a line saying

    Signed-off-by: Random J Developer <random@example.org>

using your real name (sorry, no pseudonyms or anonymous contributions.)
