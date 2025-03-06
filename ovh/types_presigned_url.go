package ovh

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type PresignedURL struct {
	Method        string            `json:"method"`
	URL           string            `json:"url"`
	SignedHeaders map[string]string `json:"signedHeaders"`
}

type PresignedURLInput struct {
	Expire    int    `json:"expire"`
	Method    string `json:"method"`
	Object    string `json:"object"`
	VersionId string `json:"versionId,omitempty"`
}

func (opts *PresignedURLInput) FromResource(d *schema.ResourceData) *PresignedURLInput {
	opts.Expire = d.Get("expire").(int)
	opts.Method = d.Get("method").(string)
	opts.Object = d.Get("object").(string)
	opts.VersionId = d.Get("version_id").(string)
	return opts
}
