package ovh

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type CloudProjectContainerRegistryIAMCreateOpts struct {
	DeleteUsers bool `json:"deleteUsers"`
}

func (opts *CloudProjectContainerRegistryIAMCreateOpts) FromResource(d *schema.ResourceData) *CloudProjectContainerRegistryIAMCreateOpts {
	opts.DeleteUsers = d.Get("delete_users").(bool)

	return opts
}
