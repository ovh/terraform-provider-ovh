package ovh

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type DedicatedServerOlaAggregationCreateOpts struct {
	Name                    string   `json:"name"`
	VirtualNetworkInterface []string `json:"virtualNetworkInterfaces"`
}

func (opts *DedicatedServerOlaAggregationCreateOpts) FromResource(d *schema.ResourceData) *DedicatedServerOlaAggregationCreateOpts {
	opts.Name = d.Get("name").(string)
	virtualNetworkInterface := d.Get("virtual_network_interfaces").([]interface{})
	// Convert virtualNetworkInterface from []interface{} to []string
	opts.VirtualNetworkInterface = make([]string, len(virtualNetworkInterface))
	for i, v := range virtualNetworkInterface {
		opts.VirtualNetworkInterface[i] = fmt.Sprint(v)
	}

	return opts
}

type DedicatedServerOlaAggregationSingleDeleteOpts struct {
	VirtualNetworkInterface string `json:"virtualNetworkInterface"`
}

type DedicatedServerOlaAggregationDeleteOpts struct {
	VirtualNetworkInterface []string `json:"virtualNetworkInterfaces"`
}

func (opts *DedicatedServerOlaAggregationDeleteOpts) FromResource(d *schema.ResourceData) *DedicatedServerOlaAggregationDeleteOpts {
	virtualNetworkInterface := d.Get("virtual_network_interfaces").([]interface{})
	// Convert virtualNetworkInterface from []interface{} to []string
	opts.VirtualNetworkInterface = make([]string, len(virtualNetworkInterface))
	for i, v := range virtualNetworkInterface {
		opts.VirtualNetworkInterface[i] = fmt.Sprint(v)
	}

	return opts
}
