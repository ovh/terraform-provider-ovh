---
name: Report a Bug
about: Let us know about an unexpected error, a crash, or otherwise incorrect behavior.
title: "[BUG]"
labels: ''
assignees: ''

---

### Describe the bug
A clear and concise description of what the bug is.

### Terraform Version

Run `terraform -v` to show the version. If you are not running the latest version of Terraform, please upgrade because your issue may have already been fixed.

### OVH Terraform Provider Version

Run `terraform init` to show the provider(s) version.

### Affected Resource(s)

Please list the resources as a list, for example:
- ovh_cloud_project_kube
- ovh_order_cart

If this issue appears to affect multiple resources, it may be an issue with [Terraform's core](https://github.com/hashicorp/terraform).

### Terraform Configuration Files

```hcl
# Copy-paste your Terraform configurations here - for large Terraform configs,
# please use a service like Dropbox and share a link to the ZIP file. For
# security, you can also encrypt the files using our GPG public key.
```

### Debug Output

Please provider a link to a GitHub Gist containing the complete debug output: https://www.terraform.io/docs/internals/debugging.html. Please do NOT paste the debug output in the issue; just paste a link to the Gist.

### Panic Output

If Terraform produced a panic, please provide a link to a GitHub Gist containing the output of the `crash.log`.

### Expected Behavior

What should have happened?

### Actual Behavior

What actually happened?

### Steps to Reproduce

Please list the steps required to reproduce the issue, for example:
1. `terraform apply`

### References

Are there any other GitHub issues (open or closed) or Pull Requests that should be linked here? For example:
- GH-1234

### Additional context

Add any other context about the problem here.
