---
subcategory : "Account Management"
---

# ovh_me_identity_provider

Configure SAML Federation (SSO) to an identity provider.

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

    # Local users will still be able to login if set to false.
    # Administrator can always login regardless of this value.
    disable_users = false

    # The assertion must contain the attribute "https://example.org/attributes/role"
    # with the allowed values being "user" or "administrator"
    requested_attributes {
      is_required = true
      name = "https://example.org/attributes/role"
      name_format = "urn:oasis:names:tc:SAML:2.0:attrname-format:uri"
      values = ["user", administrator]
    }
    # If the attribute "https://example.org/attributes/group" is available,
    # we want the IdP to provide it
    requested_attributes {
      is_required = false
      name = "https://example.org/attributes/group"
      name_format = "urn:oasis:names:tc:SAML:2.0:attrname-format:uri"
      values = []
    }
}
```

## Argument Reference

* `group_attribute_name` - The name of the attribute containing the information of which group the connecting users belong to.
* `metadata` - The SAML xml metadata of the Identity Provider to federate to.
* `disable_users` - Whether local users should still be usable as a login method or not (optional, defaults to true). Administrator will always be able to login, regardless of this value.
* `requested_attributes` A SAML 2.0 requested attribute as defined in [SAML-ReqAttrExt-v1.0](http://docs.oasis-open.org/security/saml-protoc-req-attr-req/v1.0/cs01/saml-protoc-req-attr-req-v1.0-cs01.pdf). A RequestedAttribute object will indicate that the Identity Provider should add the described attribute to the SAML assertions that will be given to the Service Provider (OVH).
  * `is_required` Expresses that this Attribute is mandatory. If the requested attribute is not present in the assertion, the user won't be allowed to log in.
  * `name` Name of the SAML Attribute that is required.
  * `name_format` NameFormat of the SAML RequestedAttribute.
  * `values`  List of AttributeValues allowed for this RequestedAttribute.

## Attributes Reference

* `creation` - Creation date of the SAML Federation.
* `last_update` - Date of the last update of the SAML Federation.
