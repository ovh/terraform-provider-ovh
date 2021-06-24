package ovh

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ovh/terraform-provider-ovh/ovh/helpers"
	"github.com/ovh/terraform-provider-ovh/ovh/helpers/hashcode"
)

func dataSourceIpLoadbalancing() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIpLoadbalancingRead,
		Schema: map[string]*schema.Schema{
			"ipv6": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					err := helpers.ValidateIpV6(v.(string))
					if err != nil {
						errors = append(errors, err)
					}
					return
				},
			},
			"ipv4": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					err := helpers.ValidateIpV4(v.(string))
					if err != nil {
						errors = append(errors, err)
					}
					return
				},
			},

			"zone": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},
			"offer": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"service_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"ip_loadbalancing": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					err := helpers.ValidateIp(v.(string))
					if err != nil {
						errors = append(errors, err)
					}
					return
				},
			},
			"state": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					err := helpers.ValidateStringEnum(v.(string), []string{"blacklisted", "deleted", "free", "ok", "quarantined", "suspended"})
					if err != nil {
						errors = append(errors, err)
					}
					return
				},
			},

			"vrack_eligibility": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"vrack_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"display_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"ssl_configuration": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					err := helpers.ValidateStringEnum(v.(string), []string{"intermediate", "modern"})
					if err != nil {
						errors = append(errors, err)
					}
					return
				},
			},

			// additional exported attributes
			"metrics_token": {
				Type:      schema.TypeString,
				Sensitive: true,
				Computed:  true,
			},
			"orderable_zone": {
				Type:     schema.TypeSet,
				Computed: true,
				Set: func(v interface{}) int {
					r := v.(map[string]interface{})
					return hashcode.String(r["name"].(string))
				},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"plan_code": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceIpLoadbalancingRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	response := []string{}
	v, exists := d.GetOk("service_name")
	if exists {
		log.Printf("[DEBUG] Will use provided iploadbalancing service")
		response = append(response, v.(string))
	} else {
		log.Printf("[DEBUG] Will list available iploadbalancing services")
		err := config.OVHClient.Get("/ipLoadbalancing", &response)
		if err != nil {
			return fmt.Errorf("Error calling /ipLoadbalancing:\n\t %q", err)
		}
	}

	filtered_iplbs := []*IpLoadbalancing{}

	for _, serviceName := range response {
		iplb := &IpLoadbalancing{}
		err := config.OVHClient.Get(fmt.Sprintf("/ipLoadbalancing/%s", serviceName), &iplb)

		if err != nil {
			return fmt.Errorf("Error calling /ipLoadbalancing/%s:\n\t %q", serviceName, err)
		}

		if v, ok := d.GetOk("ipv6"); ok && (iplb.IPv6 == nil || v.(string) != *iplb.IPv6) {
			continue
		}

		if v, ok := d.GetOk("ipv4"); ok && (iplb.IPv4 == nil || v.(string) != *iplb.IPv4) {
			continue
		}

		if v, ok := d.GetOk("zone"); ok && !zonesEquals(v.([]string), iplb.Zone) {
			continue
		}
		if v, ok := d.GetOk("offer"); ok && v.(string) != iplb.Offer {
			continue
		}
		if v, ok := d.GetOk("service_name"); ok && v.(string) != iplb.ServiceName {
			continue
		}
		if v, ok := d.GetOk("ip_loadbalancing"); ok && v.(string) != iplb.IpLoadbalancing {
			continue
		}
		if v, ok := d.GetOk("state"); ok && v.(string) != iplb.State {
			continue
		}
		if v, ok := d.GetOk("vrack_eligibility"); ok && v.(bool) != iplb.VrackEligibility {
			continue
		}
		if v, ok := d.GetOk("vrack_name"); ok && (iplb.VrackName == nil || v.(string) != *iplb.VrackName) {
			continue
		}
		if v, ok := d.GetOk("display_name"); ok && (iplb.DisplayName == nil || v.(string) != *iplb.DisplayName) {
			continue
		}
		if v, ok := d.GetOk("ssl_configuration"); ok && (iplb.SslConfiguration == nil || v.(string) != *iplb.SslConfiguration) {
			continue
		}
		filtered_iplbs = append(filtered_iplbs, iplb)
	}

	if len(filtered_iplbs) < 1 {
		return fmt.Errorf("Your query returned no results. " +
			"Please change your search criteria and try again.")
	}

	if len(filtered_iplbs) > 1 {
		return fmt.Errorf("Your query returned more than one result." +
			" Please try a more specific search criteria")
	}

	dataSourceIpLoadbalancingAttributes(d, filtered_iplbs[0])

	return nil
}

// dataSourceIpLoadbalancingAttributes populates the fields of an ipLoadbalancing datasource.
func dataSourceIpLoadbalancingAttributes(d *schema.ResourceData, iplb *IpLoadbalancing) error {
	log.Printf("[DEBUG] ovh_iploadbalancing details: %#v", iplb)

	if iplb.ServiceName == "" {
		return fmt.Errorf("serviceName cannot be empty")
	}
	if iplb.Zone == nil {
		return fmt.Errorf("zone cannot be nil")
	}
	if iplb.Offer == "" {
		return fmt.Errorf("offer cannot be empty")
	}
	if iplb.IpLoadbalancing == "" {
		return fmt.Errorf("ipLoadbalancing cannot be empty")
	}
	if iplb.State == "" {
		return fmt.Errorf("state cannot be empty")
	}

	d.SetId(iplb.ServiceName)
	d.Set("ipv6", iplb.IPv6)
	d.Set("ipv4", iplb.IPv4)
	d.Set("zone", iplb.Zone)
	d.Set("offer", iplb.Offer)
	d.Set("service_name", iplb.ServiceName)
	d.Set("ip_loadbalancing", iplb.IpLoadbalancing)
	d.Set("state", iplb.State)
	d.Set("vrack_eligibility", iplb.VrackEligibility)
	d.Set("vrack_name", iplb.VrackName)
	d.Set("display_name", iplb.DisplayName)
	d.Set("ssl_configuration", iplb.SslConfiguration)
	d.Set("metrics_token", iplb.MetricsToken)

	// Set the orderable_zone
	var orderableZone []map[string]interface{}
	for _, v := range iplb.OrderableZones {
		zone := make(map[string]interface{})
		zone["name"] = v.Name
		zone["plan_code"] = v.PlanCode

		orderableZone = append(orderableZone, zone)
	}
	err := d.Set("orderable_zone", orderableZone)
	if err != nil {
		log.Printf("[DEBUG] Unable to set orderable_zone: %s", err)
	}

	return nil
}

func zonesEquals(a, b []string) bool {
	if a == nil && b == nil {
		return true
	}

	if a == nil || b == nil {
		return false
	}

	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}
