package ovh

import (
	"fmt"
	"sort"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type DedicatedServerNetworking struct {
	Status      string                               `json:"status"`
	Description string                               `json:"description"`
	Interfaces  []DedicatedServerNetworkingInterface `json:"interfaces"`
}

type DedicatedServerNetworkingInterface struct {
	Macs []string `json:"macs"`
	Type string   `json:"type"`
}

type DedicatedServerNetworkingCreateOpts struct {
	Interfaces []DedicatedServerNetworkingInterface `json:"interfaces"`
}

type DedicatedServerNetworkingResponse struct {
	Status string
}

func (opts *DedicatedServerNetworkingCreateOpts) FromResource(d *schema.ResourceData) *DedicatedServerNetworkingCreateOpts {
	rawNetworkInterfaces := d.Get("interfaces").(*schema.Set).List()

	opts.Interfaces = make([]DedicatedServerNetworkingInterface, len(rawNetworkInterfaces))

	for i, rawNetworkInterface := range rawNetworkInterfaces {
		var networkingInterface DedicatedServerNetworkingInterface
		data := rawNetworkInterface.(map[string]interface{})

		if raw, ok := data["macs"]; ok {
			rawMacs := raw.(*schema.Set).List()
			networkingInterface.Macs = make([]string, len(rawMacs))
			for i, mac := range rawMacs {
				networkingInterface.Macs[i] = fmt.Sprint(mac)
			}

			// we want the MACs associated to an interface to be in a determist order to avoid false positive diff
			sort.Strings(networkingInterface.Macs)
		}

		if raw, ok := data["type"]; ok {
			networkingInterface.Type = raw.(string)
		}

		opts.Interfaces[i] = networkingInterface
	}

	// we want interfaces to be in a determist order to avoid false positive diff
	sort.SliceStable(opts.Interfaces, func(i, j int) bool {
		return opts.Interfaces[i].Type < opts.Interfaces[j].Type
	})

	return opts
}
