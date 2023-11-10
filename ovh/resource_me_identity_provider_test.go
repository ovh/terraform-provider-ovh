package ovh

import (
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const samlIDPMetadata string = `<?xml version="1.0" encoding="UTF-8"?>
<EntityDescriptor xmlns="urn:oasis:names:tc:SAML:2.0:metadata" entityID="https://ovhcloud.com/">
  <IDPSSODescriptor protocolSupportEnumeration="urn:oasis:names:tc:SAML:2.0:protocol">
    <KeyDescriptor use="signing">
      <KeyInfo xmlns="http://www.w3.org/2000/09/xmldsig#">
        <X509Data>
          <X509Certificate>MIIFlTCCA32gAwIBAgIUP8WQwHQwrvTa00RU9JROZAJj9ccwDQYJKoZIhvcNAQELBQAwWjELMAkGA1UEBhMCRlIxEzARBgNVBAgMClNvbWUtU3RhdGUxDDAKBgNVBAcMA1JCWDERMA8GA1UECgwIT1ZIY2xvdWQxFTATBgNVBAMMDG92aGNsb3VkLmNvbTAeFw0yMzExMDkxMDA2MjRaFw0zMzExMDYxMDA2MjRaMFoxCzAJBgNVBAYTAkZSMRMwEQYDVQQIDApTb21lLVN0YXRlMQwwCgYDVQQHDANSQlgxETAPBgNVBAoMCE9WSGNsb3VkMRUwEwYDVQQDDAxvdmhjbG91ZC5jb20wggIiMA0GCSqGSIb3DQEBAQUAA4ICDwAwggIKAoICAQC4V3HulFBksxpgkgR6KgDaSSIKkKRgyDGCF06oQN/WPxGDSHTQQHTMN7jnsbr2uJieNKh+iasGvE9JFmd6nutloL1UHoO/ecrE8P2PYgpgezl7WfyoscBDZAjWM8E9FdENnonhvlga2DgV2DGIB4+D7aN6TIPfWukOB2MjfQloA9Iw71+peO9R55S7x7zixgpLO9NovbmaAyClbz06Tsm/7ezM+Vte7BfFqGUnwuNzqgOYfQm88EqXTpCT3QfR8i2IydGgAFLMFs9YvMnCaNLw9PCN7U6VPkY6M6cFQhO/moRb3H/euJnLNRMsXp99K8ruUnQ6902NXpOOnQu5Ewzfahmx0WWvlpFGdJK34oXjaWeTodGuvHtDxCY4tiHr8jCHf9h4cmC20xAyd/V7XBtu1Pc5UAg4I0w5ehWvHtVdxCsuPEh7c4qtuLyN9Qh15r+eRbiqnWTH/xJTwfo6q6iafXXcFOlTn7WoWmmeq0R8whg6XjcxMIzBXjtynTDbQa4LVq3T8iJiGfuDgwv5OwDPRN1CsawxefETsCUQ+jf/Iw/4nZpD/YqCI5xvYtDgPSt3v2TsoOnwOSjOqKmEOoHxGTN3mhbcD+I1QKJW79zqu6GVXVwMkgWdP4pkIWGccB0FqhIVzY19xQ40DbfnCkMTv2XN4t53c/q7CYhtvyN3XwIDAQABo1MwUTAdBgNVHQ4EFgQUC8yuX4Ub/Od5jSaz7NdwHUSlq5wwHwYDVR0jBBgwFoAUC8yuX4Ub/Od5jSaz7NdwHUSlq5wwDwYDVR0TAQH/BAUwAwEB/zANBgkqhkiG9w0BAQsFAAOCAgEAR4MzroH8kEwqcZeB94hetY/NQGZI+kZ26iKnvLaZa8r56UeiIrEGdEeys5JdQh/XJDWsEU6piJ0dwrIkkpgZELmUmToylcxndzjjcHbiKLlqkL+kBu9QeO/r6JTHaNyWs0An2VvCUfo+Frt8hvrJCINlCDylOaWIxHH3P0TG7ThFGWSy8nW+VMMXDS8vQIGRM66HqgYlu6HBryecf0SsCkVYbUb1zYJ+lEhYK0pj4RORainJX+PU+mIMUwQtfBByuI7RP0a2Vny0gffrtPuNfhRJb8Pwt2UYw2niWUDOfXuk9RYgqX/1wLVqk72KJJlD3c7+abZ6BcNEJax5e/icilUrxcs4MymDPjk63kQURRVzcC4hCXYqJVQmRfVT4fdLLKPmeg3ysl+U4eJZ8odmaqoVGqZryncdAC+nT5lnLRm6m2lv3v+YhConctLxzCwV/xA8jU2w9VVRw2gkY8bdkvOb7c2OpXU6J3TYtaltG7foQiuXbRd37GWzzzEspxiAI9y8uIEJTsASaufsEdpR+a1sPy3rYJom/Li3dH9p9Ch+tp51pMYhSRGEiNu9g5918zMbrKvwkl6h/PQlTOlb65qUUoNKC5Baxhz3VkGxSKMUwS4Lj/WHvCGU5OteGFHglDgDm125FDakOYU1dnMm/P55yNhnSUH2sXngybxnw/w=</X509Certificate>
        </X509Data>
      </KeyInfo>
    </KeyDescriptor>
    <NameIDFormat>urn:oasis:names:tc:SAML:2.0:nameid-format:persistent</NameIDFormat>
    <SingleSignOnService Binding="urn:oasis:names:tc:SAML:2.0:bindings:HTTP-POST" Location="https://ovhcloud.com/"></SingleSignOnService>
    <SingleSignOnService Binding="urn:oasis:names:tc:SAML:2.0:bindings:HTTP-Redirect" Location="https://ovhcloud.com/"></SingleSignOnService>
  </IDPSSODescriptor>
</EntityDescriptor>`

func init() {
	resource.AddTestSweepers("ovh_me_identity_provider", &resource.Sweeper{
		Name: "ovh_me_identity_provider",
		F:    testSweepMeIdentityProvider,
	})
}

func testSweepMeIdentityProvider(region string) error {
	client, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting client: %s", err)
	}

	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		log.Printf("[INFO] Deleting identity provider")
		if err := client.Delete("/me/identity/provider", nil); err != nil {
			return resource.RetryableError(err)
		}

		// Successful delete
		return nil
	})

	return err
}

func TestAccMeIdentityProvider_basic(t *testing.T) {
	groupeAttribute := acctest.RandomWithPrefix(test_prefix)
	disableUsers := "false"
	config := fmt.Sprintf(testAccMeIdentityProviderConfig_basic, groupeAttribute, disableUsers, samlIDPMetadata)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCredentials(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					checkIdentityProviderResourceAttr("ovh_me_identity_provider.my_provider", groupeAttribute, disableUsers, samlIDPMetadata, nil)...,
				),
			},
		},
	})
}

func TestAccMeIdentityProvider_requestedAttributes(t *testing.T) {
	groupeAttribute := acctest.RandomWithPrefix(test_prefix)
	disableUsers := "false"
	requestedAttribute := map[string]string{
		"is_required": "false",
		"name":        "test1",
		"name_format": "test2",
		"values":      "test3",
	}
	config := fmt.Sprintf(testAccMeIdentityProviderConfig_requestedAttribute, groupeAttribute, disableUsers, samlIDPMetadata, requestedAttribute["is_required"], requestedAttribute["name"], requestedAttribute["name_format"], requestedAttribute["values"])

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCredentials(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					checkIdentityProviderResourceAttr("ovh_me_identity_provider.my_provider", groupeAttribute, disableUsers, samlIDPMetadata, requestedAttribute)...,
				),
			},
		},
	})
}

func checkIdentityProviderResourceAttr(name, group_attribute, disable_users, metadata string, requestedAttributes map[string]string) []resource.TestCheckFunc {
	checks := []resource.TestCheckFunc{}
	checks = append(checks, resource.TestCheckResourceAttr(name, "group_attribute_name", group_attribute))
	checks = append(checks, resource.TestCheckResourceAttr(name, "disable_users", disable_users))
	checks = append(checks, resource.TestCheckResourceAttr(name, "metadata", metadata+"\n"))
	if requestedAttributes != nil {
		checks = append(checks, resource.TestCheckResourceAttr(name, "requested_attributes.0.is_required", requestedAttributes["is_required"]))
		checks = append(checks, resource.TestCheckResourceAttr(name, "requested_attributes.0.name", requestedAttributes["name"]))
		checks = append(checks, resource.TestCheckResourceAttr(name, "requested_attributes.0.name_format", requestedAttributes["name_format"]))
		checks = append(checks, resource.TestCheckResourceAttr(name, "requested_attributes.0.values.0", requestedAttributes["values"]))
	}
	return checks
}

const testAccMeIdentityProviderConfig_basic = `
resource "ovh_me_identity_provider" "my_provider" {
	group_attribute_name = "%s"
	disable_users = %s
	metadata = <<EOT
%s
EOT
}
`
const testAccMeIdentityProviderConfig_requestedAttribute = `
resource "ovh_me_identity_provider" "my_provider" {
	group_attribute_name = "%s"
	disable_users = %s
	metadata = <<EOT
%s
EOT
	requested_attributes {
		is_required = %s
		name = "%s"
		name_format = "%s"
		values = [ "%s" ]
	}
}
`
