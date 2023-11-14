package ovh

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccMeIdentityProviderDataSource_basic(t *testing.T) {
	groupAttributeName := "http://schemas.xmlsoap.org/claims/Group"
	disableUsers := "false"
	reqAttributeRequired := "false"
	reqAttributeName := "identity"
	reqAttributeNameFormat := "urn:oasis:names:tc:SAML:2.0:attrname-format:basic"
	reqAttributeValue := "foobar"

	preSetup := fmt.Sprintf(
		testAccMeIdentityProviderDataSourceConfig_preSetup,
		groupAttributeName,
		samlIDPMetadata,
		disableUsers,
		reqAttributeRequired,
		reqAttributeName,
		reqAttributeNameFormat,
		reqAttributeValue,
	)
	config := fmt.Sprintf(
		testAccMeIdentityProviderDataSourceConfig_keys,
		groupAttributeName,
		samlIDPMetadata,
		disableUsers,
		reqAttributeRequired,
		reqAttributeName,
		reqAttributeNameFormat,
		reqAttributeValue,
	)

	userAttributeName := "http://schemas.xmlsoap.org/ws/2005/05/identity/claims/upn"
	ssoServiceUrl := "https://ovhcloud.com/"
	certificateExpiration := "2033-11-06T10:06:24Z"
	certificateSubject := "CN=ovhcloud.com,O=OVHcloud,L=RBX,ST=Some-State,C=FR"

	requestedAttributes := map[string]string{
		"is_required": reqAttributeRequired,
		"name":        reqAttributeName,
		"name_format": reqAttributeNameFormat,
		"values":      reqAttributeValue,
	}

	checks := checkIdentityProviderResourceAttr("ovh_me_identity_provider.sso", groupAttributeName, disableUsers, samlIDPMetadata, requestedAttributes)
	dataSourceChecks := checkIdentityProviderDataSourceAttr("data.ovh_me_identity_provider.sso", groupAttributeName, userAttributeName, ssoServiceUrl, disableUsers, certificateExpiration, certificateSubject, requestedAttributes)
	dataSourceChecks = append(dataSourceChecks, resource.TestCheckOutput("keys_present", "true"))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCredentials(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: preSetup,
				Check:  resource.ComposeTestCheckFunc(checks...),
			}, {
				Config: config,
				Check:  resource.ComposeTestCheckFunc(dataSourceChecks...),
			},
		},
	})
}

func checkIdentityProviderDataSourceAttr(name, group_attribute, user_attribute, sso_service_url, disable_users, certificateExpiration, certificateSubject string, requestedAttributes map[string]string) []resource.TestCheckFunc {
	checks := []resource.TestCheckFunc{}
	checks = append(checks, resource.TestCheckResourceAttr(name, "group_attribute_name", group_attribute))
	checks = append(checks, resource.TestCheckResourceAttr(name, "user_attribute_name", user_attribute))
	checks = append(checks, resource.TestCheckResourceAttr(name, "sso_service_url", sso_service_url))
	checks = append(checks, resource.TestCheckResourceAttr(name, "disable_users", disable_users))
	checks = append(checks, resource.TestCheckResourceAttr(name, "idp_signing_certificates.0.expiration", certificateExpiration))
	checks = append(checks, resource.TestCheckResourceAttr(name, "idp_signing_certificates.0.subject", certificateSubject))
	if requestedAttributes != nil {
		checks = append(checks, resource.TestCheckResourceAttr(name, "requested_attributes.0.is_required", requestedAttributes["is_required"]))
		checks = append(checks, resource.TestCheckResourceAttr(name, "requested_attributes.0.name", requestedAttributes["name"]))
		checks = append(checks, resource.TestCheckResourceAttr(name, "requested_attributes.0.name_format", requestedAttributes["name_format"]))
		checks = append(checks, resource.TestCheckResourceAttr(name, "requested_attributes.0.values.0", requestedAttributes["values"]))
	}
	return checks
}

const testAccMeIdentityProviderDataSourceConfig_preSetup = `
resource "ovh_me_identity_provider" "sso" {
    group_attribute_name = "%s"
	metadata = <<EOT
%s
EOT
    disable_users = %s
    requested_attributes {
      is_required = %s
      name = "%s"
      name_format = "%s"
      values = ["%s"]
    }
}`

const testAccMeIdentityProviderDataSourceConfig_keys = `
resource "ovh_me_identity_provider" "sso" {
    group_attribute_name = "%s"
	metadata = <<EOT
%s
EOT
    disable_users = %s
    requested_attributes {
      is_required = %s
      name = "%s"
      name_format = "%s"
      values = ["%s"]
    }
}

data "ovh_me_identity_provider" "sso" {}

output "keys_present" {
	value = tostring(
		data.ovh_me_identity_provider.sso.group_attribute_name == ovh_me_identity_provider.sso.group_attribute_name &&
		data.ovh_me_identity_provider.sso.requested_attributes.0.name == ovh_me_identity_provider.sso.requested_attributes.0.name
	)
}
`
