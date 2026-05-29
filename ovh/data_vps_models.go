package ovh

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ovh/go-ovh/ovh"
	"github.com/ovh/terraform-provider-ovh/v2/ovh/helpers/hashcode"
)

func vpsModelElemSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name":    {Type: schema.TypeString, Computed: true},
			"offer":   {Type: schema.TypeString, Computed: true},
			"version": {Type: schema.TypeString, Computed: true},
			"vcore":   {Type: schema.TypeInt, Computed: true},
			"memory":  {Type: schema.TypeInt, Computed: true},
			"disk":    {Type: schema.TypeInt, Computed: true},
			"datacenter": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"available_options": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"maximum_additionnal_ip": {Type: schema.TypeInt, Computed: true},
		},
	}
}

func flattenVPSModels(models []VPSModel) []map[string]interface{} {
	out := make([]map[string]interface{}, 0, len(models))
	for _, m := range models {
		out = append(out, map[string]interface{}{
			"name":                   m.Name,
			"offer":                  m.Offer,
			"version":                m.Version,
			"vcore":                  m.Vcore,
			"memory":                 m.Memory,
			"disk":                   m.Disk,
			"datacenter":             m.Datacenter,
			"available_options":      m.AvailableOptions,
			"maximum_additionnal_ip": m.MaximumAdditionnalIp,
		})
	}
	return out
}

func dataSourceVPSModels() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceVPSModelsRead,
		Schema: map[string]*schema.Schema{
			"service_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"models": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     vpsModelElemSchema(),
			},
		},
	}
}

func dataSourceVPSModelsRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)

	endpoint := fmt.Sprintf("/vps/%s/models", url.PathEscape(serviceName))
	models := []VPSModel{}
	if err := config.OVHClient.Get(endpoint, &models); err != nil {
		if apiErr, ok := err.(*ovh.APIError); ok && apiErr.Code == 404 {
			msg := apiErr.Message
			switch {
			case strings.Contains(msg, "Got an invalid (or empty) URL"):
				return fmt.Errorf(
					"the OVHcloud API endpoint %s is not available on this VPS lineup. "+
						"This data source may only work on legacy VPS plans, or the endpoint "+
						"may have been deprecated. See the data source's documentation for "+
						"supported VPS generations.",
					endpoint)
			case strings.Contains(msg, "does not exist"):
				return fmt.Errorf(
					"the requested resource at %s does not exist (the VPS may not have "+
						"the required option subscribed, or the resource ID is wrong)",
					endpoint)
			}
		}
		return fmt.Errorf("calling GET %s: %w", endpoint, err)
	}

	names := make([]string, 0, len(models))
	for _, m := range models {
		names = append(names, m.Name)
	}
	d.SetId(fmt.Sprintf("%s/%s", serviceName, hashcode.Strings(names)))
	d.Set("models", flattenVPSModels(models))
	return nil
}
