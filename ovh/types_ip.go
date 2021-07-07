package ovh

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type IpReverse struct {
	IpReverse string `json:"ipReverse"`
	Reverse   string `json:"reverse"`
}

func (v IpReverse) ToMap() map[string]interface{} {
	obj := make(map[string]interface{})
	obj["ip_reverse"] = v.IpReverse
	obj["reverse"] = v.Reverse
	return obj
}

type IpReverseCreateOpts struct {
	IpReverse string `json:"ipReverse"`
	Reverse   string `json:"reverse"`
}

func (opts *IpReverseCreateOpts) FromResource(d *schema.ResourceData) *IpReverseCreateOpts {
	opts.IpReverse = d.Get("ip_reverse").(string)
	opts.Reverse = d.Get("reverse").(string)

	return opts
}
