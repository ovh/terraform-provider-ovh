package ovh

import (
	"fmt"
	"sort"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceDedicatedServerNetworking() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceDedicatedServerNetworkingRead,
		Schema: map[string]*schema.Schema{
			"service_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The internal name of your dedicated server.",
			},

			// Computed
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Operation description",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Operation status",
			},
			"interfaces": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Interface or interfaces aggregation.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"macs": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Interface Mac address",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Interface type",
						},
						"aggregation_fallback": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Mac address of the LACP fallback interface",
						},
					},
				},
			},
		},
	}
}

func dataSourceDedicatedServerNetworkingRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)

	networkingDetails, err := getDedicatedServerNetworkingDetails(serviceName, config.OVHClient)
	if err != nil {
		return fmt.Errorf("Error retrieving networking details for %s:\n\t %q", serviceName, err)
	}

	d.SetId(serviceName)
	d.Set("description", networkingDetails.Description)
	d.Set("status", networkingDetails.Status)

	interfaces := make([]map[string]interface{}, len(networkingDetails.Interfaces))
	for i, details := range networkingDetails.Interfaces {
		iface := map[string]interface{}{
			"type": details.Type,
		}

		macs := append([]string(nil), details.Macs...)
		sort.Strings(macs)
		iface["macs"] = macs

		if details.AggregationFallback != nil {
			iface["aggregation_fallback"] = *details.AggregationFallback
		}

		interfaces[i] = iface
	}

	sort.SliceStable(interfaces, func(i, j int) bool {
		return interfaces[i]["type"].(string) < interfaces[j]["type"].(string)
	})

	if err := d.Set("interfaces", interfaces); err != nil {
		return fmt.Errorf("Error persisting interfaces in state for %s:\n\t %q", serviceName, err)
	}

	return nil
}
