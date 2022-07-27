package ovh

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type CloudProjectKubeOIDCCreateOpts struct {
	ClientID  string `json:"clientId"`
	IssuerUrl string `json:"issuerUrl"`
}

type CloudProjectKubeOIDCUpdateOpts struct {
	ClientID  string `json:"clientId"`
	IssuerUrl string `json:"issuerUrl"`
}

type CloudProjectKubeOIDCResponse struct {
	ClientID  string `json:"clientId"`
	IssuerUrl string `json:"issuerUrl"`
}

func (opts *CloudProjectKubeOIDCCreateOpts) FromResource(d *schema.ResourceData) *CloudProjectKubeOIDCCreateOpts {
	return &CloudProjectKubeOIDCCreateOpts{
		ClientID:  d.Get("client_id").(string),
		IssuerUrl: d.Get("issuer_url").(string),
	}
}

func (opts *CloudProjectKubeOIDCUpdateOpts) FromResource(d *schema.ResourceData) *CloudProjectKubeOIDCUpdateOpts {
	return &CloudProjectKubeOIDCUpdateOpts{
		ClientID:  d.Get("client_id").(string),
		IssuerUrl: d.Get("issuer_url").(string),
	}
}

func (v CloudProjectKubeOIDCResponse) ToMap() map[string]interface{} {
	obj := make(map[string]interface{})
	obj["client_id"] = v.ClientID
	obj["issuer_url"] = v.IssuerUrl

	return obj
}
