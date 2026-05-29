package ovh

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ovh/go-ovh/ovh"
)

func dataSourceVPSStatus() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceVPSStatusRead,
		Schema: map[string]*schema.Schema{
			"service_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"probes": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"service": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"port": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"state": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceVPSStatusRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)

	probes := []VPSStatusProbe{}
	endpoint := fmt.Sprintf(
		"/vps/%s/status",
		url.PathEscape(serviceName),
	)
	if err := config.OVHClient.Get(endpoint, &probes); err != nil {
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

	out := make([]map[string]interface{}, 0, len(probes))
	for _, p := range probes {
		out = append(out, map[string]interface{}{
			"service": p.Service,
			"port":    p.Port,
			"state":   p.State,
		})
	}

	d.SetId(serviceName)
	d.Set("probes", out)
	return nil
}
