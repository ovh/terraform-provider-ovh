package ovh

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type PresignedURL struct {
	Method string `json:"method"`
	URL    string `json:"url"`
}

type PresignedURLInput struct {
	Expire int    `json:"expire"`
	Method string `json:"method"`
	Object string `json:"object"`
}

func (opts *PresignedURLInput) FromResource(d *schema.ResourceData) *PresignedURLInput {
	opts.Expire = d.Get("expire").(int)
	opts.Method = d.Get("method").(string)
	opts.Object = d.Get("object").(string)
	return opts
}
