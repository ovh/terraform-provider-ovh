package ovh

import (
	"fmt"
	"log"
	"net/url"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceIpService() *schema.Resource {
	return &schema.Resource{
		Read: func(d *schema.ResourceData, meta interface{}) error {
			return dataSourceIpServiceRead(d, meta)
		},
		Schema: map[string]*schema.Schema{
			"service_name": {
				Type:        schema.TypeString,
				Description: "The service name",
				Required:    true,
			},

			//computed
			"can_be_terminated": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"country": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "Custom description on your ip",
				Computed:    true,
			},

			"ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"organisation_id": {
				Type:        schema.TypeString,
				Description: "IP block organisation Id",
				Computed:    true,
			},
			"routed_to": {
				Type:        schema.TypeList,
				Description: "Routage information",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"service_name": {
							Type:        schema.TypeString,
							Description: "Service where ip is routed to",
							Computed:    true,
						},
					},
				},
			},
			"type": {
				Type:        schema.TypeString,
				Description: "Possible values for ip type",
				Computed:    true,
			},
		},
	}
}

func dataSourceIpServiceRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	serviceName := d.Get("service_name").(string)

	log.Printf("[DEBUG] Will read ip service %s", serviceName)
	endpoint := fmt.Sprintf("/ip/service/%s",
		url.PathEscape(serviceName),
	)

	ip := &IpService{}
	if err := config.OVHClient.Get(endpoint, &ip); err != nil {
		return fmt.Errorf("Error calling GET %s:\n\t %q", endpoint, err)
	}

	for k, v := range ip.ToMap() {
		if k != "id" {
			d.Set(k, v)
		}
	}

	d.SetId(serviceName)

	return nil
}
