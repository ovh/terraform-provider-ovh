package ovh

import (
	"fmt"
	"net/url"
	"sort"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ovh/go-ovh/ovh"
)

func dataSourceVPSTemplate() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceVPSTemplateRead,
		Schema: map[string]*schema.Schema{
			"service_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"template_id": {
				Type:     schema.TypeInt,
				Required: true,
			},

			// Computed
			"name":         {Type: schema.TypeString, Computed: true},
			"distribution": {Type: schema.TypeString, Computed: true},
			"bit_format":   {Type: schema.TypeInt, Computed: true},
			"available_language": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"locale": {Type: schema.TypeString, Computed: true},
			"software_ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeInt},
			},
		},
	}
}

func dataSourceVPSTemplateRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)
	templateID := d.Get("template_id").(int)

	t := &vpsTemplate{}
	tplEndpoint := fmt.Sprintf(
		"/vps/%s/templates/%d",
		url.PathEscape(serviceName),
		templateID,
	)
	if err := config.OVHClient.Get(tplEndpoint, t); err != nil {
		if apiErr, ok := err.(*ovh.APIError); ok && apiErr.Code == 404 {
			msg := apiErr.Message
			switch {
			case strings.Contains(msg, "Got an invalid (or empty) URL"):
				return fmt.Errorf(
					"the OVHcloud API endpoint %s is not available on this VPS lineup. "+
						"This data source may only work on legacy VPS plans, or the endpoint "+
						"may have been deprecated. See the data source's documentation for "+
						"supported VPS generations.",
					tplEndpoint)
			case strings.Contains(msg, "does not exist"):
				return fmt.Errorf(
					"the requested resource at %s does not exist (the VPS may not have "+
						"the required option subscribed, or the resource ID is wrong)",
					tplEndpoint)
			}
		}
		return fmt.Errorf("calling GET %s: %w", tplEndpoint, err)
	}

	bf, _ := strconv.Atoi(t.BitFormat)

	softwareIDs := []int{}
	swEndpoint := fmt.Sprintf(
		"/vps/%s/templates/%d/software",
		url.PathEscape(serviceName),
		templateID,
	)
	if err := config.OVHClient.Get(swEndpoint, &softwareIDs); err != nil {
		if apiErr, ok := err.(*ovh.APIError); ok && apiErr.Code == 404 {
			msg := apiErr.Message
			switch {
			case strings.Contains(msg, "Got an invalid (or empty) URL"):
				return fmt.Errorf(
					"the OVHcloud API endpoint %s is not available on this VPS lineup. "+
						"This data source may only work on legacy VPS plans, or the endpoint "+
						"may have been deprecated. See the data source's documentation for "+
						"supported VPS generations.",
					swEndpoint)
			case strings.Contains(msg, "does not exist"):
				return fmt.Errorf(
					"the requested resource at %s does not exist (the VPS may not have "+
						"the required option subscribed, or the resource ID is wrong)",
					swEndpoint)
			}
		}
		return fmt.Errorf("calling GET %s: %w", swEndpoint, err)
	}
	sort.Ints(softwareIDs)

	d.SetId(fmt.Sprintf("%s_%d", serviceName, templateID))
	d.Set("name", t.Name)
	d.Set("distribution", t.Distribution)
	d.Set("bit_format", bf)
	d.Set("available_language", t.AvailableLanguage)
	d.Set("locale", t.Locale)
	d.Set("software_ids", softwareIDs)
	return nil
}
