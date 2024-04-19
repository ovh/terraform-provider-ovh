---
subcategory : "Account Management"
---

# ovh_me_ssh_key (Data Source)

-> __NOTE__ This data source will be removed in next release.

Use this data source to retrieve information about an SSH key.

## Example Usage

```hcl
data "ovh_me_ssh_key" "mykey" {
  key_name = "mykey"
}
```

## Argument Reference

-> __NOTE__ This data source will be removed in next release.


* `key_name` - (Required) The name of the SSH key.

## Attributes Reference

-> __NOTE__ This data source will be removed in next release.


* `key_name` - See Argument Reference above.
* `key` - The content of the public key.
E.g.: "ssh-ed25519 AAAAC3..."
* `default` - True when this public SSH key is used for rescue mode and reinstallations.
