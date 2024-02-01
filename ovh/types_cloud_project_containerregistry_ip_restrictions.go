package ovh

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type CloudProjectContainerRegistryIPRestriction struct {
	Description string `json:"description,omitempty"`
	IPBlock     string `json:"ipBlock"`
}

type CloudProjectContainerRegistryIPRestrictionCreateOpts struct {
	IPRestrictions []CloudProjectContainerRegistryIPRestriction `json:"ipRestrictions"`
}

func (opts *CloudProjectContainerRegistryIPRestrictionCreateOpts) FromResource(d *schema.ResourceData) *CloudProjectContainerRegistryIPRestrictionCreateOpts {
	opts.IPRestrictions = loadIPRestrictionsFromResource(d.Get("ip_restrictions"))

	return opts
}

func loadIPRestrictionsFromResource(i interface{}) []CloudProjectContainerRegistryIPRestriction {
	ips := make([]CloudProjectContainerRegistryIPRestriction, 0)

	iprestrictionsSet := i.([]interface{})

	for _, ipSet := range iprestrictionsSet {
		ips = append(ips, CloudProjectContainerRegistryIPRestriction{
			Description: ipSet.(map[string]interface{})["description"].(string),
			IPBlock:     ipSet.(map[string]interface{})["ip_block"].(string),
		})
	}

	return ips
}

func (r CloudProjectContainerRegistryIPRestriction) ToMap() map[string]interface{} {
	obj := make(map[string]interface{})
	obj["description"] = r.Description
	obj["ip_block"] = r.IPBlock

	return obj
}
