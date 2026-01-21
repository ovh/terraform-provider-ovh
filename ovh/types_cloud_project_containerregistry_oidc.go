package ovh

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type CloudProjectContainerRegistryOIDCProvider struct {
	// Fields required, omitempty for patch, requirements checked with ValidateParameters
	Name         string `json:"name"`
	Endpoint     string `json:"endpoint"`
	ClientID     string `json:"clientId"`
	ClientSecret string `json:"clientSecret"`
	Scope        string `json:"scope"`

	// Optional fields
	GroupFilter string `json:"groupFilter,omitempty"`
	GroupsClaim string `json:"groupsClaim,omitempty"`
	AdminGroup  string `json:"adminGroup,omitempty"`
	VerifyCert  bool   `json:"verifyCert,omitempty"`
	AutoOnboard bool   `json:"autoOnboard,omitempty"`
	UserClaim   string `json:"userClaim,omitempty"`
}

type CloudProjectContainerRegistryOIDCCreateOpts struct {
	DeleteUsers bool                                      `json:"deleteUsers"`
	Provider    CloudProjectContainerRegistryOIDCProvider `json:"provider"`
}

type CloudProjectContainerRegistryOIDCUpdateOpts struct {
	Provider CloudProjectContainerRegistryOIDCProvider `json:"provider"`
}

type CloudProjectContainerRegistryOIDCResponse struct {
	// Fields required, omitempty for patch, requirements checked with ValidateParameters
	Name     string `json:"name"`
	Endpoint string `json:"endpoint"`
	ClientID string `json:"clientId"`
	Scope    string `json:"scope"`

	// Non required
	GroupFilter string `json:"groupFilter"`
	GroupsClaim string `json:"groupsClaim"`
	AdminGroup  string `json:"adminGroup"`
	VerifyCert  bool   `json:"verifyCert"`
	AutoOnboard bool   `json:"autoOnboard"`
	UserClaim   string `json:"userClaim"`
}

func (opts *CloudProjectContainerRegistryOIDCCreateOpts) FromResource(d *schema.ResourceData) *CloudProjectContainerRegistryOIDCCreateOpts {
	opts.DeleteUsers = d.Get("delete_users").(bool)
	opts.Provider = *newCloudProjectContainerRegistryOIDCProvider(d)

	return opts
}

func (opts *CloudProjectContainerRegistryOIDCUpdateOpts) FromResource(d *schema.ResourceData) *CloudProjectContainerRegistryOIDCUpdateOpts {
	opts.Provider = *newCloudProjectContainerRegistryOIDCProvider(d)

	return opts
}

func newCloudProjectContainerRegistryOIDCProvider(d *schema.ResourceData) *CloudProjectContainerRegistryOIDCProvider {
	return &CloudProjectContainerRegistryOIDCProvider{
		Name:         d.Get("oidc_name").(string),
		Endpoint:     d.Get("oidc_endpoint").(string),
		ClientID:     d.Get("oidc_client_id").(string),
		ClientSecret: d.Get("oidc_client_secret").(string),
		Scope:        d.Get("oidc_scope").(string),
		GroupFilter:  d.Get("oidc_group_filter").(string),
		GroupsClaim:  d.Get("oidc_groups_claim").(string),
		AdminGroup:   d.Get("oidc_admin_group").(string),
		VerifyCert:   d.Get("oidc_verify_cert").(bool),
		AutoOnboard:  d.Get("oidc_auto_onboard").(bool),
		UserClaim:    d.Get("oidc_user_claim").(string),
	}
}

func (v CloudProjectContainerRegistryOIDCResponse) ToMap() map[string]interface{} {
	obj := make(map[string]interface{})
	obj["oidc_name"] = v.Name
	obj["oidc_endpoint"] = v.Endpoint
	obj["oidc_client_id"] = v.ClientID
	obj["oidc_scope"] = v.Scope
	obj["oidc_group_filter"] = v.GroupFilter
	obj["oidc_groups_claim"] = v.GroupsClaim
	obj["oidc_admin_group"] = v.AdminGroup
	obj["oidc_verify_cert"] = v.VerifyCert
	obj["oidc_auto_onboard"] = v.AutoOnboard
	obj["oidc_user_claim"] = v.UserClaim

	return obj
}
