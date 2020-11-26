package ovh

import (
	"fmt"
	"net/url"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceIpLoadbalancingVrackNetwork() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIpLoadbalancingVrackNetworkRead,
		Schema: map[string]*schema.Schema{
			"service_name": {
				Type:        schema.TypeString,
				Description: "The internal name of your IPloadbalancer",
				Required:    true,
			},

			"vrack_network_id": {
				Type:        schema.TypeInt,
				Description: "Internal Load Balancer identifier of the vRack private network",
				Required:    true,
			},

			//Computed
			"display_name": {
				Type:        schema.TypeString,
				Description: "Human readable name for your vrack network",
				Computed:    true,
			},
			"nat_ip": {
				Type:        schema.TypeString,
				Description: "An IP block used as a pool of IPs by this Load Balancer to connect to the servers in this private network. The blck must be in the private network and reserved for the Load Balancer",
				Computed:    true,
			},
			"subnet": {
				Type:        schema.TypeString,
				Description: "IP block of the private network in the vRack",
				Computed:    true,
			},
			"vlan": {
				Type:        schema.TypeInt,
				Description: "VLAN of the private network in the vRack. 0 if the private network is not in a VLAN",
				Computed:    true,
			},
		},
	}
}

func dataSourceIpLoadbalancingVrackNetworkRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	endpoint := fmt.Sprintf(
		"/ipLoadbalancing/%s/vrack/network/%d",
		url.PathEscape(d.Get("service_name").(string)),
		d.Get("vrack_network_id").(int),
	)

	vn := &IpLoadbalancingVrackNetwork{}
	if err := config.OVHClient.Get(endpoint, &vn); err != nil {
		return fmt.Errorf("Error calling GET %s:\n\t %q", endpoint, err)
	}

	// set resource attributes
	for k, v := range vn.ToMap() {
		d.Set(k, v)
	}

	d.SetId(fmt.Sprintf("%s/%d", d.Get("service_name").(string), d.Get("vrack_network_id").(int)))
	return nil
}
