package ovh

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ovh/terraform-provider-ovh/v2/ovh/helpers"
)

type CloudProjectKubeOIDCCreateOpts struct {
	ClientID       string   `json:"clientId"`
	IssuerUrl      string   `json:"issuerUrl"`
	UsernameClaim  string   `json:"usernameClaim"`
	UsernamePrefix string   `json:"usernamePrefix"`
	GroupsClaim    []string `json:"groupsClaim"`
	GroupsPrefix   string   `json:"groupsPrefix"`
	RequiredClaim  []string `json:"requiredClaim"`
	SigningAlgs    []string `json:"signingAlgorithms"`
	CaContent      string   `json:"caContent"`
}

type CloudProjectKubeOIDCUpdateOpts struct {
	ClientID       string   `json:"clientId"`
	IssuerUrl      string   `json:"issuerUrl"`
	UsernameClaim  string   `json:"usernameClaim"`
	UsernamePrefix string   `json:"usernamePrefix"`
	GroupsClaim    []string `json:"groupsClaim"`
	GroupsPrefix   string   `json:"groupsPrefix"`
	RequiredClaim  []string `json:"requiredClaim"`
	SigningAlgs    []string `json:"signingAlgorithms"`
	CaContent      string   `json:"caContent"`
}

type CloudProjectKubeOIDCResponse struct {
	ClientID       string   `json:"clientId"`
	IssuerUrl      string   `json:"issuerUrl"`
	UsernameClaim  string   `json:"usernameClaim"`
	UsernamePrefix string   `json:"usernamePrefix"`
	GroupsClaim    []string `json:"groupsClaim"`
	GroupsPrefix   string   `json:"groupsPrefix"`
	RequiredClaim  []string `json:"requiredClaim"`
	SigningAlgs    []string `json:"signingAlgorithms"`
	CaContent      string   `json:"caContent"`
}

func (opts *CloudProjectKubeOIDCCreateOpts) FromResource(d *schema.ResourceData) *CloudProjectKubeOIDCCreateOpts {
	opts.ClientID = d.Get("client_id").(string)
	opts.IssuerUrl = d.Get("issuer_url").(string)
	opts.UsernameClaim = d.Get("oidc_username_claim").(string)
	opts.UsernamePrefix = d.Get("oidc_username_prefix").(string)
	opts.GroupsClaim, _ = helpers.StringsFromSchema(d, "oidc_groups_claim")
	opts.GroupsPrefix = d.Get("oidc_groups_prefix").(string)
	opts.RequiredClaim, _ = helpers.StringsFromSchema(d, "oidc_required_claim")
	opts.SigningAlgs, _ = helpers.StringsFromSchema(d, "oidc_signing_algs")
	opts.CaContent = d.Get("oidc_ca_content").(string)

	return opts
}

func (opts *CloudProjectKubeOIDCUpdateOpts) FromResource(d *schema.ResourceData) *CloudProjectKubeOIDCUpdateOpts {
	opts.ClientID = d.Get("client_id").(string)
	opts.IssuerUrl = d.Get("issuer_url").(string)
	opts.UsernameClaim = d.Get("oidc_username_claim").(string)
	opts.UsernamePrefix = d.Get("oidc_username_prefix").(string)
	opts.GroupsClaim, _ = helpers.StringsFromSchema(d, "oidc_groups_claim")
	opts.GroupsPrefix = d.Get("oidc_groups_prefix").(string)
	opts.RequiredClaim, _ = helpers.StringsFromSchema(d, "oidc_required_claim")
	opts.SigningAlgs, _ = helpers.StringsFromSchema(d, "oidc_signing_algs")
	opts.CaContent = d.Get("oidc_ca_content").(string)

	return opts
}

func (v CloudProjectKubeOIDCResponse) ToMap() map[string]interface{} {
	obj := make(map[string]interface{})
	obj["client_id"] = v.ClientID
	obj["issuer_url"] = v.IssuerUrl
	obj["oidc_username_claim"] = v.UsernameClaim
	obj["oidc_username_prefix"] = v.UsernamePrefix
	obj["oidc_groups_claim"] = v.GroupsClaim
	obj["oidc_groups_prefix"] = v.GroupsPrefix
	obj["oidc_required_claim"] = v.RequiredClaim
	obj["oidc_signing_algs"] = v.SigningAlgs
	obj["oidc_ca_content"] = v.CaContent

	return obj
}
