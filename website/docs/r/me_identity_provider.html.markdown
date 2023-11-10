---
subcategory : "Account Management"
---

# ovh_me_identity_provider

Configure SAML Fedration (SSO) to an identity provider.

## Example Usage

```hcl
resource "ovh_me_identity_provider" "sso" {
    group_attribute_name = "http://schemas.xmlsoap.org/claims/Group"
    metadata = <<EOT
<?xml version="1.0" encoding="UTF-8"?>
<EntityDescriptor xmlns="urn:oasis:names:tc:SAML:2.0:metadata" entityID="https://ovhcloud.com/">
  <IDPSSODescriptor protocolSupportEnumeration="urn:oasis:names:tc:SAML:2.0:protocol">
    <KeyDescriptor use="signing">
      <KeyInfo xmlns="http://www.w3.org/2000/09/xmldsig#">
        <X509Data>
          <X509Certificate>MIIFlTCCA32gAwIBAgIUP8WQwHQwrvTa00RU9JROZAJj9ccwDQYJKoZIhvcNAQELBQAwWjELMAkGA1UEBhMCRlIxEzARBgNVBAgMClNvbWUtU3RhdGUxDDAKBgNVBAcMA1JCWDERMA8GA1UECgwIT1ZIY2xvdWQxFTATBgNVBAMMDG92aGNsb3VkLmNvbTAeFw0yMzExMDkxMDA2[...]xA8jU2w9VVRw2gkY8bdkvOb7c2OpXU6J3TYtaltG7foQiuXbRd37GWzzzEspxiAI9y8uIEJTsASaufsEdpR+a1sPy3rYJom/Li3dH9p9Ch+tp51pMYhSRGEiNu9g5918zMbrKvwkl6h/PQlTOlb65qUUoNKC5Baxhz3VkGxSKMUwS4Lj/WHvCGU5OteGFHglDgDm125FDakOYU1dnMm/P55yNhnSUH2sXngybxnw/w=</X509Certificate>
        </X509Data>
      </KeyInfo>
    </KeyDescriptor>
    <NameIDFormat>urn:oasis:names:tc:SAML:2.0:nameid-format:persistent</NameIDFormat>
    <SingleSignOnService Binding="urn:oasis:names:tc:SAML:2.0:bindings:HTTP-POST" Location="https://ovhcloud.com/"></SingleSignOnService>
    <SingleSignOnService Binding="urn:oasis:names:tc:SAML:2.0:bindings:HTTP-Redirect" Location="https://ovhcloud.com/"></SingleSignOnService>
  </IDPSSODescriptor>
</EntityDescriptor>
EOT

    disable_users = false

    requested_attributes {
      is_required = false
      name = "group"
      name_format = "urn:oasis:names:tc:SAML:2.0:attrname-format:basic"
      values = ["test"]
    }
    requested_attributes {
      is_required = false
      name = "email"
      name_format = "urn:oasis:names:tc:SAML:2.0:attrname-format:basic"
      values = ["test@example.org"]
    }
}
```

## Argument Reference

* `group_attribute_name` - The name of the attribute containing the information of which group the connecting users belong to.
* `metadata` - The SAML xml metadata of the Identity Provider to federate to.
* `disable_users` - Whether account users should still be usable as a login method or not (optional, defaults to true).
* `requested_attributes` A SAML 2.0 requested attribute that should be added to SAML requests when using this provider (optional).
  * `is_required` Expresses that this RequestedAttribute is mandatory.
  * `name` Name of the SAML RequestedAttribute.
  * `name_format` NameFormat of the SAML RequestedAttribute.
  * `values`  List of AttributeValues allowed for this RequestedAttribute

## Attributes Reference

* `creation` - Creation date of the SAML Federation.
* `last_update` - Date of the last update of the SAML Federation.
