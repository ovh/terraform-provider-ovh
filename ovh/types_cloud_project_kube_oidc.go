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
	opts.ClientID = d.Get(kubeOidcClientIdKey).(string)
	opts.IssuerUrl = d.Get(kubeOidcIssuerUrlKey).(string)
	opts.UsernameClaim = d.Get(kubeOidcUsernameClaimKey).(string)
	opts.UsernamePrefix = d.Get(kubeOidcUsernamePrefixKey).(string)
	opts.GroupsClaim, _ = helpers.StringsFromSchema(d, kubeOidcGroupsClaimKey)
	opts.GroupsPrefix = d.Get(kubeOidcGroupsPrefixKey).(string)
	opts.RequiredClaim, _ = helpers.StringsFromSchema(d, kubeOidcRequiredClaimKey)
	opts.SigningAlgs, _ = helpers.StringsFromSchema(d, kubeOidcSigningAlgsKey)
	opts.CaContent = d.Get(kubeOidcCaContentKey).(string)

	return opts
}

func (opts *CloudProjectKubeOIDCUpdateOpts) FromResource(d *schema.ResourceData) *CloudProjectKubeOIDCUpdateOpts {
	opts.ClientID = d.Get(kubeOidcClientIdKey).(string)
	opts.IssuerUrl = d.Get(kubeOidcIssuerUrlKey).(string)
	opts.UsernameClaim = d.Get(kubeOidcUsernameClaimKey).(string)
	opts.UsernamePrefix = d.Get(kubeOidcUsernamePrefixKey).(string)
	opts.GroupsClaim, _ = helpers.StringsFromSchema(d, kubeOidcGroupsClaimKey)
	opts.GroupsPrefix = d.Get(kubeOidcGroupsPrefixKey).(string)
	opts.RequiredClaim, _ = helpers.StringsFromSchema(d, kubeOidcRequiredClaimKey)
	opts.SigningAlgs, _ = helpers.StringsFromSchema(d, kubeOidcSigningAlgsKey)
	opts.CaContent = d.Get(kubeOidcCaContentKey).(string)

	return opts
}

func (v CloudProjectKubeOIDCResponse) ToMap() map[string]interface{} {
	obj := make(map[string]interface{})
	obj[kubeOidcClientIdKey] = v.ClientID
	obj[kubeOidcIssuerUrlKey] = v.IssuerUrl
	obj[kubeOidcUsernameClaimKey] = v.UsernameClaim
	obj[kubeOidcUsernamePrefixKey] = v.UsernamePrefix
	obj[kubeOidcGroupsClaimKey] = v.GroupsClaim
	obj[kubeOidcGroupsPrefixKey] = v.GroupsPrefix
	obj[kubeOidcRequiredClaimKey] = v.RequiredClaim
	obj[kubeOidcSigningAlgsKey] = v.SigningAlgs
	obj[kubeOidcCaContentKey] = v.CaContent

	return obj
}
