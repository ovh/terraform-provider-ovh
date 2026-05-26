package ovh

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ovh/go-ovh/ovh"
)

type vpsSoftware struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Type   string `json:"type"`
	Status string `json:"status"`
}

func dataSourceVPSTemplateSoftware() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceVPSTemplateSoftwareRead,
		Schema: map[string]*schema.Schema{
			"service_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"template_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"software_id": {
				Type:     schema.TypeInt,
				Required: true,
			},

			// Computed
			"name":   {Type: schema.TypeString, Computed: true},
			"type":   {Type: schema.TypeString, Computed: true},
			"status": {Type: schema.TypeString, Computed: true},
		},
	}
}

func dataSourceVPSTemplateSoftwareRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)
	templateID := d.Get("template_id").(int)
	softwareID := d.Get("software_id").(int)

	sw := &vpsSoftware{}
	endpoint := fmt.Sprintf(
		"/vps/%s/templates/%d/software/%d",
		url.PathEscape(serviceName),
		templateID,
		softwareID,
	)
	if err := config.OVHClient.Get(endpoint, sw); err != nil {
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

	d.SetId(fmt.Sprintf("%s_%d_%d", serviceName, templateID, softwareID))
	d.Set("name", sw.Name)
	d.Set("type", sw.Type)
	d.Set("status", sw.Status)
	return nil
}
